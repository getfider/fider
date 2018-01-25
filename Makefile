BUILD_TIME=$(shell date +"%Y.%m.%d.%H%M%S")

# Building
build:
	rm -rf dist
	go build -ldflags='-s -w -X main.buildtime=${BUILD_TIME}' -o fider .
	./node_modules/.bin/webpack -p

lint: 
	./node_modules/.bin/tslint -c tslint.json 'public/**/*.{ts,tsx}'

# Testing
test:
	godotenv -f .test.env go test ./... -cover -p=1

coverage:
	godotenv -f .test.env courtney -o cover.out $$(go list ./...)

e2e-single:
	./scripts/e2e.sh single

e2e-multi:
	./scripts/e2e.sh multi

e2e-build:
	./scripts/e2e.sh build

e2e:
	./scripts/e2e.sh

# Development
watch:
	rm -rf dist
	gin --buildArgs "-ldflags='-s -w -X main.buildtime=${BUILD_TIME}'" & 
	./node_modules/.bin/webpack --watch

run:
	godotenv -f .env ./fider

.DEFAULT_GOAL := build
