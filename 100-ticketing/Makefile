export PATH := $(PATH):/usr/local/go/bin/:$(GOPATH)/bin/

help:
	@echo "to start this apps, run 'example-docker-start'"
	@echo "to stop this apps, run 'example-docker-stop'"

build:
	@go build -o ./bin/app
run:
	@go run main.go
swagger:
	@echo "creating swagger api documentation"
	@swag init --outputTypes yaml --output docs
example-docker-start:
	@echo "creating example docker using file docker-compose-example.yaml"
	@docker compose -f docker-compose-example.yaml up --build -t1 -d
example-docker-stop:
	@echo "creating example docker using file docker-compose-example.yaml"
	@docker compose -f docker-compose-example.yaml down -t1