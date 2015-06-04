# chickenlittle
A RESTful service to get ahold of people, quickly.

# Notice
This is not ready for public consumption.  Sorry.

# People API

## Get list of people
**Request**
```
GET /people
```

**Example Response**
```json
{
  "people": [
    {
      "username": "arthur",
      "fullname": "King Arthur"
    },
    {
      "username": "lancelot",
      "fullname": "Sir Lancelot"
    }
  ],
  "message": "",
  "error": ""
}
```

## Fetch details for a person
**Request**
```
GET /people/USERNAME
```

**Example Response**
```json
{
  "people": [
    {
      "username": "lancelot",
      "fullname": "Sir Lancelot"
    }
  ],
  "message": "",
  "error": ""
}
```

## Create a new person
**Request**
```
POST /people

{
  "username": "lancelot",
  "fullname": "Sir Lancelot"
}
```

**Example Response**
```json
{
  "message": "User lancelot created",
  "error": ""
}
```

## Update a person
**Request**
```
PUT /people/USERNAME

{
  "fullname": "Sir Lancelot the Brave"
}
```

**Example Response**
```json
{
  "people": [
    {
      "username": "lancelot",
      "fullname": "Sir Lancelot the Brave"
    }
  ],
  "message": "User lancelot updated",
  "error": ""
}
```

## Delete a person
**Request**
```
DELETE /people/USERNAME
```

**Example Response**
```json
{
  "message": "User lancelot deleted",
  "error": ""
}
```

# Notification Plan API

The notfication plan is stored as a JSON array of steps to take when notifying a person.  The order of the array is the order in which the steps are followed.  An example plan looks like this:

```json
[
    {
        "method": 1,
        "data": "2108675309",
        "notify_every_period": 0,
        "notify_until_period": 900000000000
    },
        {
        "method": 3,
        "data": "lancelot@roundtable.org.uk",
        "notify_every_period": 0,
        "notify_until_period": 900000000000
    },
    {
        "method": 2,
        "data": "2105551212",
        "notify_every_period": 300000000000,
        "notify_until_period": 0
    }
]
```

The fields of a step are as follows:

| Field | Description |
|:-------|:-------------|
|```method```| **Method of notification**  (1 = Voice, 2 = SMS, 3 = E-mail)|
|```data```| **The data relevant to the method of notification** Can be phone number, SMS number, or email address. *double-quotes are mandatory, even for numbers!*|
|```notify_every_period```|**Period of time in which to repeat a notification**  Time is stored in nanoseconds.  1 minute = 60000000000.  This is only relevant to the *last* notification step in the array, since the last step is the only one repeated *ad infinitum* until the person responds.  A ```0``` value indicates that this step will only be followed once and not repeated.  If this field is set for a step that's not the last in the array, it will be ignored. |
|```notify_unit_period```|**Period of time in which the service waits for a response before proceeding to the next notification step in the array**  Time is stored in nanoseconds.  1 minute = 60000000000.  A ```0``` value is not valid for this field and will result in the step being skipped.  If this field is set for the very last step in the array, it will be ignored. |
