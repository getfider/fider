BUILD_TIME="$(shell date +"%Y.%m.%d.%H%M%S")"

test:
	go test $$(go list ./... | grep -v /vendor/) -cover

lint:
	golint -set_exit_status

setup-ci:
	go get -u github.com/kardianos/govendor
	go get -u github.com/golang/lint/golint
	go get -u github.com/mattn/goveralls
	govendor sync

run-ci: lint build
	goveralls -service=travis-ci
  if [ "${TRAVIS_PULL_REQUEST}" = "false" ]; then
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
    docker login -u $DOCKER_USER -p $DOCKER_PASS
    export REPO=wecanhearyou/wchy-api
    export TAG=`if [ "$TRAVIS_BRANCH" == "dev" ]; then echo "staging"; else echo "latest" ; fi`
    docker build -f Dockerfile -t $REPO:$TAG .
    docker push $REPO
  fi

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