# Add Route Service to Import, Create, Find, Edit Ticketing Route

Create route services to manage route record, and also route data can be imported from csv file:

- Add `route` service to import route data from csv file, create, find and edit route.
- Please refer to `/docs/swagger.yaml` for detail API Contract

- Sample import route file: `/sample_data/ticket_data.csv`

- This App service and MongoDB service should be running on docker and all service running from docker-compose
- docker-compose file should have mongodb service configuration, node app configurations and all config ENV should be declared on the docker-compose configurations.

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
  - SECRET - ```aOFNMxyVIZfAANsT```
  - SECRET_KEY - ```BctbLulGvxijNQKi```

## Prequisite

`jwt-authenticatio` course

## Directory structure
    - /path/to/base/directory
        - docker-compose-example.yaml  => docker compose example configuration
        - main.go => main golang file
        - app/router  => routing configuration
        - app/controllers  => routing controllers
        - app/adapter  => for external connection, use for mongodb
        - app/models => for get and insert data to mongodb
        - sample_data => directory csv file

## Quick Start
- Enter `cd /path/to/base/directory`
- Clone this repo:

      git clone https://gitlab.playcourt.id/dpe/curriculum/be-golang/60-upload-file.git
- Enter repo directory `cd 60-upload-file`
- Copy example docker compose config `cp docker-compose-example.yaml docker-compose.yaml`
- Edit config file with your parameter. this will be your local configuration, not synced to git.
- Build and start services:

       docker compose -f docker-compose.yaml up --build -t1 -d

- Stop services:

      docker compose -f docker-compose.yaml down -t1
- Test

  - Import route from file (This feature only accessible for role admin)

        curl -X 'POST' 'http://localhost:10081/api/v1/route/import' \
        -H 'Authorization: Bearer <token>' \
        -H 'Content-Type: application/json' \
        -F 'routeFile=@<file_location>'

  - Create new route (This feature only accessible for role admin)

        curl -X 'POST' 'http://localhost:10081/api/v1/route' \
        -H 'Authorization: Bearer <token>' \
        -H 'Content-Type: application/json' \
        -d '{"from":<from_city>, "to": <destination_city>,"price": <price>,"departureTime": <departure_time>}'

  - Get route with filter

        curl -X 'GET' 'http://localhost:10081/api/v1/route?from=<city>' \
        -H 'Authorization: Bearer <token>' \
        -H 'Content-Type: application/json'

  - Get Route by id

        curl -X 'GET' 'http://localhost:10081/api/v1/route/<route_id>' \
        -H 'Authorization: Bearer <token>' \
        -H 'Content-Type: application/json'

  - Edit route by id (This feature only accessible for role admin)

        curl -X 'PUT' 'http://localhost:10081/api/v1/route/<route_id>' \
        -H 'Authorization: Bearer <token>' \
        -H 'Content-Type: application/json' \
        -d '{"price": <price>,"departureTime": <departure_time>}'

- Asciinema Preview

[![asciicast](https://asciinema.org/a/569673.png)](https://asciinema.org/a/569673)
