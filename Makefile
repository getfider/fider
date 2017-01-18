test:
	go test

lint:
	golint -set_exit_status

setup:
	go get -u github.com/kardianos/govendor
	go get github.com/golang/lint/golint
	govendor sync

ci: lint test build

migrate:
ifeq ($(ENV), local)
	migrate -url postgres://wchy:wchy-pw@localhost:5555/wchy?sslmode=disable -path ./db/migrations up
else
	migrate -url ${DATABASE_URL} -path ./db/migrations up
endif

build:
	go build -o wchy-api main.go

run: build
	./wchy-api

.DEFAULT_GOAL := build