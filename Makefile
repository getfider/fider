BUILD_TIME=$(shell date +"%Y.%m.%d.%H%M%S")

# Building
build:
	rm -rf dist
	go build -ldflags='-s -w -X main.buildtime=${BUILD_TIME}' -o fider .
	NODE_ENV=production npx webpack -p

lint: 
	npx tslint -c tslint.json 'public/**/*.{ts,tsx}' 'tests/**/*.{ts,tsx}'

lint-fix: 
	npx tslint -c tslint.json 'public/**/*.{ts,tsx}' 'tests/**/*.{ts,tsx}' --fix

# Testing
test-ui:
	rm -rf ./output/public
	npx tsc -p ./tsconfig.json
	npx jest ./output/public

test-server:
	godotenv -f .test.env go test ./... -p=1

test : test-server test-ui

coverage:
	godotenv -f .test.env go test ./... -p=1 -coverprofile=cover.out -coverpkg=all

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
	npx webpack -w

run:
	godotenv -f .env ./fider

.DEFAULT_GOAL := build
