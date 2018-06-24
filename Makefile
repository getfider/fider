BUILD_TIME=$(shell date +"%Y.%m.%d.%H%M%S")
BUILD_NUMBER = $(shell echo $$CIRCLE_BUILD_NUM)

# Building
build-ui:
	rm -rf dist
	NODE_ENV=production npx webpack -p

build-server:
	go build -ldflags='-s -w -X main.buildtime=${BUILD_TIME} -X main.buildnumber=${BUILD_NUMBER}' -o fider .

build : build-server build-ui

lint: 
	npx tslint -c tslint.json 'public/**/*.{ts,tsx}' 'tests/**/*.{ts,tsx}'

lint-fix: 
	npx tslint -c tslint.json 'public/**/*.{ts,tsx}' 'tests/**/*.{ts,tsx}' --fix

# Testing
test-ui:
	TZ='GMT' npx jest ./public

test-server: build-server
	godotenv -f .test.env ./fider migrate
	godotenv -f .test.env go test ./... -race

test-e2e:
	./scripts/e2e.sh

test : test-server test-ui

coverage: build-server
	godotenv -f .test.env ./fider migrate
	godotenv -f .test.env go test ./... -coverprofile=cover.out -coverpkg=all -race

# Development
watch:
	rm -rf ./dist
	mkdir ./dist
	air -c air.conf & 
	npx webpack -w

run:
	godotenv -f .env ./fider

.DEFAULT_GOAL := build
