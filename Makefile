SHELL := /bin/bash

.PONEY: all run tslint

build-web:
	npx webpack
build:
	go build -ldflags -s -w -X main.buildtime=2020.10.28.153959 -X main.buildnumber= -o fider .
migrate: 
	go run main.go migrate
run-dev:
	env $(grep -v '^#' .env | xargs -0) go run main.go
run:
	go run main.go
tslint:
	npx tslint  'public/**/*.{ts,tsx}' 'tests/**/*.{ts,tsx}'
testui:
	npx jest ./public
lint:
	golangci-lint run
test: 
	docker-compose -f docker-compose-test.yml up -d
	sleep 15
	godotenv -f .test.env go run main.go migrate
	godotenv -f .test.env go test ./... -race
	docker-compose -f docker-compose-test.yml kill && docker-compose -f docker-compose-test.yml rm -f	
