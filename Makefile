BUILD_TIME="$(shell date +"%Y.%m.%d.%H%M%S")"

test:
	go test $$(go list ./... | grep -v /vendor/) -coverprofile=coverage.out

lint:
	golint -set_exit_status

setup-ci:
	go get -u github.com/kardianos/govendor
	go get -u github.com/golang/lint/golint
	go get -u github.com/mattn/goveralls
	govendor sync

run-ci: lint build
	goveralls -service=travis-ci

migrate:
ifeq ($(ENV), local)
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