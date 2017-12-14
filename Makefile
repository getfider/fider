BUILD_TIME=$(shell date +"%Y.%m.%d.%H%M%S")

ENV_FILE=.test.env

ifeq ($(TRAVIS), true)
ENV_FILE=.ci.env
endif
	

# Building
build:
	rm -rf dist
	go build -ldflags='-s -w -X main.buildtime=${BUILD_TIME}' -o fider .
	webpack -p

lint: 
	./node_modules/.bin/tslint -c tslint.json 'public/**/*.{ts,tsx}'

# Testing
test:
	godotenv -f ${ENV_FILE} go test ./... -cover -p=1

coverage:
	godotenv -f ${ENV_FILE} courtney -o cover.out $$(go list ./...)

e2e:
	./scripts/e2e.sh

# Development
watch:
	rm -rf dist
	gin --buildArgs "-ldflags='-s -w -X main.buildtime=${BUILD_TIME}'" & 
	webpack --watch

run:
	godotenv -f .env ./fider

.DEFAULT_GOAL := build
