BUILD_TIME="$(shell date +"%Y.%m.%d.%H%M%S")"

# Load .test.env file and export all env variables
include .test.env
export $(shell sed 's/=.*//' .test.env)

define tag_docker
	@if [ "$(TRAVIS_BRANCH)" = "master" ]; then \
		docker tag $(1) $(1):stable; \
	fi
	@if [ "$(TRAVIS_BRANCH)" = "dev" ]; then \
		docker tag $(1) $(1):staging; \
	fi
endef

test:
	go test $$(go list ./... | grep -v /vendor/) -cover

setup-ci:
	go get github.com/kardianos/govendor
	govendor sync
	govendor install +vendor

goveralls:
	goveralls -service=travis-ci

dockerize:
ifeq ($(TRAVIS_PULL_REQUEST), false)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
	docker build -f Dockerfile -t wecanhearyou/wchy .
	$(call tag_docker, wecanhearyou/wchy)
	docker push wecanhearyou/wchy
endif

migrate:
ifeq ($(GO_ENV), development)
	migrate -url postgres://wchy:wchy-pw@localhost:5555/wchy?sslmode=disable -path ./migrations up
else
	migrate -url ${DATABASE_URL} -path ./migrations up
endif

build:
	go build -ldflags='-X main.buildtime=${BUILD_TIME}'

watch:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'"

run:
	wchy

.DEFAULT_GOAL := build
