# Simple Ticketing App

## Description

This is demo ticketing application build with Golang, JWT Authentication, MongoDB, packed as docker container.

`docker-compose` is used to start all needed services during dev.

On prod env, you should not use this docker-compose yml:

- mongodb is usually a separate cluster of multiple servers


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
  - SERVICE_ENDPOINT_USER - ```http://localhost:8000/api/v1/user```
  - SERVICE_ENDPOINT_ROUTE - ```http://localhost:8000/api/v1/route```
  - SERVICE_ENDPOINT_TICKET - ```http://localhost:8000/api/v1/ticket```
  - ORIGIN - `*`

## Project Acceptance Criteria

- This apps should have ticket service, that only will get request from booking service
- This apps should have booking service, to create booking, update booking status completed or cancelled.
  - when create new booking, booking service automatically request to ticket service to create new ticket
- This apps should have CORS to restrict API access from specific domain
- Communication between services (booking, ticket, user) must be done through RestAPI, this aims to maintain independently of each service and database.

## API Documentation

   Please check `/docs/swagger.yaml` or open `http://localhost:8000/swagger/index.html`

## Sample Route CSV file

Please check `/sample_data/ticket_data.csv`

## Directory structure
    - /path/to/base/directory
        - docker-compose-example.yaml  => docker compose example configuration
        - main.go => main golang file
        - app/router  => routing configuration
        - app/controllers  => routing controllers
        - app/adapter  => for external connection, use for mongodb
        - app/models => for get and insert data to mongodb
        - sample_data => directory csv file
        - Makefile => using for generated api docs

## Quick Start
- Enter `cd /path/to/base/directory`
- Clone this repo:

      git clone https://gitlab.playcourt.id/dpe/curriculum/be-golang/100-ticketing.git
- Enter repo directory `cd 100-ticketing`
- Copy example docker compose config `cp docker-compose-example.yaml docker-compose.yaml`
- Edit config file with your parameter. this will be your local configuration, not synced to git.
- Build and start services:

       docker compose -f docker-compose.yaml up --build -t1 -d

- Stop services:

      docker compose -f docker-compose.yaml down -t1

- Asciinema Preview

[![asciicast](https://asciinema.org/a/572808.png)](https://asciinema.org/a/572808)

