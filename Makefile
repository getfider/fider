SHELL := /bin/bash

.PONEY: all run tslint

build-web:
	npx webpack
migrate: 
	go run main.go migrate
run:
	go run main.go
tslint:
	npx tslint  'public/**/*.{ts,tsx}' 'tests/**/*.{ts,tsx}'
testui:
	npx jest ./public
lint:
	golangci-lint run
test:
	docker-compose -f docker-compose-ci.yml up -d
	sleep 30
	godotenv -f .ci.env go run main.go migrate
	godotenv -f .ci.env go test ./... -race
	docker-compose -f docker-compose-ci.yml kill && docker-compose -f docker-compose-ci.yml rm -f	
