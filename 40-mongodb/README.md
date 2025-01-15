# Using Mongodb in Golang

Using Mongodb in Golang. Based on `routing` course:

## Setup Go Server And Mongodb
- Use the following additional ENV for go server, mongo server
    - ADDRESS - ```0.0.0.0:8000```
    - MONGO_USERNAME - ```root```
    - MONGO_PASSWORD - ```devtestpassword```
    - MONGO_HOST - ```mongodb```
    - MONGO_PORT - ```27017```
    - MONGO_DBNAME - ```default_database```
    - MONGO_INITDB_ROOT_USERNAME - ```root```
    - MONGO_INITDB_ROOT_PASSWORD  - ```devtestpassword```

## Prequisite

`routing` course

## Project Acceptance Criteria:
- Candidate should be able to add new API route to create new item and get item list
- Request GET to root url should response json string:

      { "message": "Hello. Try GET/POST to /api/v1/item" }

- Request GET to http://{{host}}/api/v1/item should response:
  - Response content-type: application/json
  - Response status: 200 (OK)

          {
              "count": 3,
              "data": [
                        {
                            "_id": "633448951994c4a3b515ce19",
                            "name": "testName",
                            "qty": 1,
                            "createdAt": "2022-09-28T20:13:57.090Z",
                            "updatedAt": "2022-09-28T20:24:35.945Z"
                        },
                        {
                            "_id": "63344fec1994c4a3b515d14b",
                            "name": "item1",
                            "qty": 1,
                            "createdAt": "2022-09-28T20:45:16.415Z",
                            "updatedAt": "2022-09-28T20:45:38.928Z"
                        },
                        {
                             "_id": "633452f91994c4a3b515d298",
                             "name": "test1",
                             "qty": 1,
                             "createdAt": "2022-09-28T20:58:17.719Z",
                             "updatedAt": "2022-09-28T22:31:45.348Z"
                        }
              ]
          }
  
- Request POST to http://{{host}}/api/v1/item
    - Request content-type: application/json
    - Have request validation to check request body should have field:
        - name
        - qty
    - Response body if request body don't have name or qty object
      - Response content-type: application/json
      - Response status: 400 (Bad Request)
      
            { "message": "name and qty cannot be empty" }

   - Response body if request body have name and qty object
        - Response content-type: application/json
        - Response status: 200 (OK)

              {
                 "_id": "633452f91994c4a3b515d298",
                 "name": "test1",
                 "qty": 1,
                 "createdAt": "2022-09-28T20:58:17.719Z",
                 "updatedAt": "2022-09-28T21:06:31.427Z"
              }


## Directory structure
    - /path/to/base/directory
        - docker-compose-example.yaml  => docker compose example configuration
        - main.go => main golang file
        - app/controllers  => routing controllers
        - app/adapter  => for external connection, use for mongodb
        - app/models => for get and insert data to mongodb

## Quick Start
- Enter `cd /path/to/base/directory`
- Clone this repo:

      git clone https://gitlab.playcourt.id/dpe/curriculum/be-golang/40-mongodb.git
- Enter repo directory `cd 40-mongodb`
- Copy example docker compose config `cp docker-compose-example.yaml docker-compose.yaml`
- Edit config file with your parameter. this will be your local configuration, not synced to git.
- Build and start services:

       docker compose -f docker-compose.yaml up --build -t1 -d
- Test:

      curl -v http://localhost:8000/
      curl -v http://localhost:8000/api/v1/item
      curl -v -d '{"name":"testName", "qty":1}' -H "Content-Type: application/json" -X POST http://localhost:8000/api/v1/item
- Stop services:

      docker compose -f docker-compose.yaml down -t1
- Asciinema Preview

[![asciicast](https://asciinema.org/a/567953.png)](https://asciinema.org/a/567953)