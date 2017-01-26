BUILD_TIME="$(shell date +"%Y.%m.%d.%H%M%S")"

define tag_docker
	@if [ "$(TRAVIS_BRANCH)" = "master" -a "$(TRAVIS_PULL_REQUEST)" = "false" ]; then \
		docker tag $(1) $(1):stable; \
	fi
	@if [ "$(TRAVIS_BRANCH)" = "dev" -a "$(TRAVIS_PULL_REQUEST)" = "false" ]; then \
		docker tag $(1) $(1):staging; \
	fi
	@if [ "$(TRAVIS_PULL_REQUEST)" != "false" ]; then \
		docker tag $(1) $(1):PR_$(TRAVIS_PULL_REQUEST); \
	fi
endef

test:
	go test $$(go list ./... | grep -v /vendor/) -cover

lint:
	golint -set_exit_status

setup-ci:
	go get -u github.com/kardianos/govendor
	go get -u github.com/golang/lint/golint
	go get -u github.com/mattn/goveralls
	govendor sync

run-ci: lint
	goveralls -service=travis-ci
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
	docker build -f Dockerfile -t wecanhearyou/wchy-api .
	$(call tag_docker, wecanhearyou/wchy-api)
	docker push wecanhearyou/wchy-api

migrate:
ifeq ($(ENV), development)
	migrate -url postgres://wchy:wchy-pw@localhost:5555/wchy?sslmode=disable -path ./migrations up
else
	migrate -url ${DATABASE_URL} -path ./migrations up
endif

build:
	go build -ldflags='-X main.buildtime=${BUILD_TIME}'

watch:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'"

run:
	wchy-api

.DEFAULT_GOAL := build