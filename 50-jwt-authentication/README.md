# Using JWT Authentication in Golang

Using JWT authentication in Golang. Based on `mongodb` course:

- Add `user` service to record user data, create admin, create user and login
- Add `auth` middleware to generate token and verify token
- Add Bearer token authentication to access API user/create, item/create
- User entity have data name, email and password (password minimum 6 character)

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

`mongodb` course

## Project Acceptance Criteria:

- In this project we need to create user service with JWT authentication
- Create new role model with schema:

      name: String (unique)

- Create new user model with schema:

      email: String (unique)
      name: String
      password: String (minimal length 6 characters)
      roles: Array of role
      createdAt: Datetime
      updatedAt: Datetime

- Create these 2 roles automatically when the server starts for the first time only.

  Role to be created:

    - admin
    - user

- This User service will have features:

    - Create admin user (need secret-key to access this feature)
    - Login user
    - Get Me (get current login user profile)

  Manage user features, below features only available for role admin:

    - Create user
    - Find all user
    - Find all user filter by name
    - Find user by id
    - Update user by id
    - Delete user by id

- Update existing item routes:

    - Find all items (This feature accessible by public user without authentication)

      Request GET to http://{{host}}/api/v1/item

    - Create new item (Only login user can access this feature )

      Request POST to http://{{host}}/api/v1/item

- Please refer to `/docs/swagger.yaml` for the detail API Contract
- This App service and MongoDB service should will running on docker and all service running from docker-compose
- docker-compose file should have mongodb service configuration, node app configurations and all config ENV should be declared on the docker-compose configurations.

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

      git clone https://gitlab.playcourt.id/dpe/curriculum/be-golang/50-jwt-authentication.git
- Enter repo directory `cd 50-jwt-authentication`
- Copy example docker compose config `cp docker-compose-example.yaml docker-compose.yaml`
- Edit config file with your parameter. this will be your local configuration, not synced to git.
- Build and start services:

       docker compose -f docker-compose.yaml up --build -t1 -d

- Stop services:

      docker compose -f docker-compose.yaml down -t1

- Asciinema Preview

[![asciicast](https://asciinema.org/a/568788.png)](https://asciinema.org/a/568788)
