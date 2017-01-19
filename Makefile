BUILD_TIME="$(shell date +"%Y.%m.%d.%H%M%S")"

test:
	go test

lint:
	golint -set_exit_status

setup:
	go get -u github.com/kardianos/govendor
	go get -u github.com/golang/lint/golint
	govendor sync

ci: lint test build

migrate:
ifeq ($(ENV), local)
	migrate -url postgres://wchy:wchy-pw@localhost:5555/wchy?sslmode=disable -path ./db/migrations up
else
	go get -u github.com/mattes/migrate
	migrate -url ${DATABASE_URL} -path ./db/migrations up
endif

build:
	go build -ldflags='-X main.buildtime=${BUILD_TIME}'

watch:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'"

run:
	wchy-api

.DEFAULT_GOAL := build