# Routing in Golang

Example of routing with golang, adding GET and POST route

## Setup Go Server
- Use the following additional ENV for go server
    - ADDRESS - ```0.0.0.0:8000```

## Project Acceptance Criteria:
- Candidate should be able to add new API route to create new item and get item list
- Request GET to root url should response json string:

      { "message": "Hello. Try GET to /api/v1/item" }

- Request GET to http://{{host}}/api/v1/item should response:
  - Response content-type: application/json
  - Response status: 200 (OK)

          { "message": "This route is a placeholder that will get all items later" }
  
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

              { "message": "This route is a placeholder that will save name (name), qty (qty) later" }


## Directory structure
    - /path/to/base/directory
        - docker-compose-example.yaml  => docker compose example configuration
        - main.go => main golang file
        - app/controllers  => routing controllers

## Quick Start
- Enter `cd /path/to/base/directory`
- Clone this repo:

      git clone https://gitlab.playcourt.id/dpe/curriculum/be-golang/30-routing.git
- Enter repo directory `cd 30-routing`
- `cp docker-compose-example.yaml docker-compose.yaml`
  This will be your local configuration, not synced to git.
- Build and start services:

       docker compose -f docker-compose.yaml up --build -t1 -d
- Test:

      curl -v http://localhost:8000/
      curl -v http://localhost:8000/api/v1/item
      curl -v -d '{"name":"testName", "qty":1}' -H "Content-Type: application/json" -X POST http://localhost:8000/api/v1/item
- Stop services:

      docker compose -f docker-compose.yaml down -t1
- Asciinema Preview

[![asciicast](https://asciinema.org/a/567952.png)](https://asciinema.org/a/567952)