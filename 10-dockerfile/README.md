# Managing Golang Container with Dockerfile

A basic go webapp, packed as docker container
This project is a simple Hello World web server application develop with Go.

## Project Criteria:

- Web server developed with Go
- Web server should response with string "Hello World" when client request to the server URL
- Web server application should run in docker container

## Prequisite

Docker Desktop or other docker variants working. See `docker desktop` course

## Build & Run Example

- Build container, tag it as `dockerfile-go-image`

      docker build . -t dockerfile-go-image
- Run built image, name it `dockerfile-go-test`

      docker run --name dockerfile-go-test --rm -d -p 10081:8080 dockerfile-go-image
- Access web server

      curl http://localhost:10081
- Stop container. Since we start the container with `--rm`, it will be automatically deleted when stopped

      docker kill dockerfile-go-test
- Asciinema Preview

[![asciicast](https://asciinema.org/a/567943.png)](https://asciinema.org/a/567943)
