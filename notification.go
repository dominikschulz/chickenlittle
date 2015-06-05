package main

import (
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	stopChan chan string
)

type NotificationsInProgress struct {
	Stoppers map[string]chan bool
	Mu       sync.Mutex
}

func StartNotificationEngine() {
	// Initialize our map of Stopper channels
	// UUID -> channel
	NIP.Stoppers = make(map[string]chan bool)

	log.Println("StartNotificationEngine()")

	for {

		select {
		// IMPROVEMENT: We could implement a close() of planChan to indicate that the service is shutting down
		//              and instruct all notifications to cease
		case np := <-planChan:

			// We've received a new notification plan
			log.Printf("Got plan: %+v\n", np)

			// Get the plan's UUID
			id := np.ID.String()

			NIP.Mu.Lock()

			// Create a new Stopper channel for this plan
			NIP.Stoppers[id] = make(chan bool)

			// Launch a goroutine to handle plan processing
			go planHandler(np, NIP.Stoppers[id])

			NIP.Mu.Unlock()
		case stopUUID := <-stopChan:
			// We've received a request to stop a notification plan
			NIP.Mu.Lock()

			// Check to see if the requested UUID is actually in progress
			_, prs := NIP.Stoppers[stopUUID]
			if prs {
				// It's in progress, so we'll send a message on its Stopper to
				// be received by the goroutine executing the plan
				NIP.Stoppers[stopUUID] <- true
			}
			NIP.Mu.Unlock()
		}
	}

}

func planHandler(np *NotificationPlan, sc <-chan bool) {

	var timerChan <-chan time.Time
	var tickerChan <-chan time.Time

	log.Println("planHandler()")

	uuid := np.ID.String()

	for n, s := range np.Steps {

		log.Println("[", uuid, "]", "STEP", n)
		log.Println("[", uuid, "]", "Method:", s.Method)

		if n == len(np.Steps)-1 {
			// Last step, so we use a Ticker and NotifyEveryPeriod
			tickerChan = time.NewTicker(s.NotifyEveryPeriod).C
			log.Println("[", uuid, "]", "[Waiting", strconv.FormatFloat(s.NotifyEveryPeriod.Minutes(), 'f', 1, 64), "minutes]")

		} else {
			// Not the last step, so we use a Timer and NotifyUntilPeriod
			timerChan = time.NewTimer(s.NotifyUntilPeriod).C
			log.Println("[", uuid, "]", "[Waiting", strconv.FormatFloat(s.NotifyUntilPeriod.Minutes(), 'f', 1, 64), "minutes]")
		}

	timerLoop:
		for {
			select {
			case <-timerChan:
				log.Println("[", uuid, "]", "Step timer expired.  Proceeding!")
				break timerLoop
			case <-tickerChan:
				log.Println("[", uuid, "]", "**Tick**  Retry contact method!")
				log.Println("[", uuid, "]", "Waiting", strconv.FormatFloat(s.NotifyEveryPeriod.Minutes(), 'f', 1, 64), "minutes]")
			case <-sc:
				log.Println("[", uuid, "]", "Manual stop received.  Exiting loops.")
				NIP.Mu.Lock()
				defer NIP.Mu.Unlock()
				delete(NIP.Stoppers, uuid)
				return
			}
		}

		log.Println("Loop broken")
	}
}
