BUILD_TIME="$(shell date +"%Y.%m.%d.%H%M%S")"

define tag_docker
	@if [ "$(TRAVIS_BRANCH)" = "master" ]; then \
		docker tag $(1) $(1):stable; \
	fi
	@if [ "$(TRAVIS_BRANCH)" = "dev" ]; then \
		docker tag $(1) $(1):staging; \
	fi
endef

test:
	GO_ENV=test ginkgo -r -cover

lint:
	golint -set_exit_status

setup-ci:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/kardianos/govendor
	go get github.com/modocache/gover
	go get github.com/golang/lint/golint
	go get github.com/mattn/goveralls
	govendor sync

run-ci: test
	gover
	goveralls -coverprofile=gover.coverprofile -service=travis-ci
ifeq ($(TRAVIS_PULL_REQUEST), false)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
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