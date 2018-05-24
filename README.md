# SP-Friends

SP-friends is a friends management web service written in Go, and use `neo4j` as database.

## Prerequisites

1. Golang
2. Docker
3. Docker Compose

## Setup (Mac)

Following setup using `docker-machine` is used and tested on Mac.

1. Create `docker-machine create default` if no default machine.
2. Setup env with `eval "$(docker-machine env default)"`.
3. Clone repo into `$GOPATH/src/github.com/bobintornado/spfriends`.
4. `dep ensure`
4. Run `make all`.
5. Run `make up` to start the app

## Dev

Run `make all && make up` to rebuild the app and start again

Run `make test` to run all the tests.

## Todo

1. Add more tests

## Tech-Stack

Given the nature of complex relationships between people (friend, subscribe and block), I choose to use a graph database to model the data.

I choose neo4j since it's popular(community support), mature(tested) and easy to start.

I design the API as more of a RPC style rather than RESTful style given the query is easier thought as relationship oritened rather than resources oriented.

## Starter-Kit

Use [ardanlabs/service](https://github.com/ardanlabs/service) as the starter kit.

# API

## Overall style

RPC via http, using JSON as data encoding format.

## Overview

| endpoint             | purpose                                                         | method | 
|----------------------|-----------------------------------------------------------------|--------|
| /v1/createFriendship | create a friendship relationship between two people             | POST   |
| /v1/getFriendsOfUser | get all friendship relationship of a single person              | POST   |
| /v1/getCommonFriends | get commons friends of two people                               | POST   |
| /v1/subscribe        | create a subscribe relationship from a person to another person | POST   |
| /v1/block            | create a block relationship from a person to another person     | POST   |
| /v1/getUpdateList    | get the list of people who can get the update given a update    | POST   |

## Example

Assuming `DOCKER_HOST` is `192.168.99.100`.

### createFriendship
Request
```
curl -X POST \
  http://192.168.99.100:3000/v1/createFriendship \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
  "friends":
    [
      "andy@example.com",
      "john@example.com"
    ]
}'
```
Successful Response 
```
{
    "success": true
}
```
Error Response 
```
{
  "error": "friendship request: \u0026{Friends:[andy@example.com john@example.com]}: User block existed, can't create friendship"
}
```

### getFriendsOfUser
Request 
```
curl -X POST \
  http://192.168.99.100:3000/v1/getFriendsOfUser \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "andy@example.com"
}'
```

Successful Response 
```
{
    "success": true,
    "friends": [
        "john@example.com"
    ],
    "count": 1
}
```
Error Response 
```
{
    "error": "stack trace or reason of error"
}
```

### getCommonFriends
Request 
```
curl -X POST \
  http://192.168.99.100:3000/v1/getCommonFriends \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
  "friends":
    [
      "andy@example.com",
      "john@example.com"
    ]
}'
```
Successful Response when there is common friends
```
{
    "success": true,
    "friends": [
        "bob@example.com"
    ],
    "count": 1
}
```

Successful Response when there is no common friends
```
{
    "success": true,
}
```

### Subscribe
Request
```
curl -X POST \
  http://192.168.99.100:3000/v1/subscribe \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
  "requestor": "lisa@example.com",
  "target": "john@example.com"
}'
```

Successful Response 
```
{
    "success": true
}
```
Error Response 
```
{
    "error": "stack trace or reason of error"
}
```

### Block
Request
```
curl -X POST \
  http://192.168.99.100:3000/v1/block \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
  "requestor": "andy@example.com",
  "target": "john@example.com"
}'
```
Successful Response
```
{
    "success": true
}
```
Error Response 
```
{
    "error": "stack trace or reason of error"
}
```

### getUpdateList
Request
```
curl -X POST \
  http://192.168.99.100:3000/v1/getUpdateList \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 75d4fdc6-c8fe-4b20-8f0c-9045a9316c2f' \
  -d '{
  "sender":  "john@example.com",
  "text": "Hello World! kate@example.com and bob@gmail.com"
}'
```

Successful Response
```
{
    "success": true,
    "recipients": [
        "kate@example.com",
        "bob@gmail.com",
        "abc@example.com",
        "mcc@example.com",
        "bob@example.com",
        "lisa@example.com"
    ]
}
```
