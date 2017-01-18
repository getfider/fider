test:
	go test

lint:
	golint -set_exit_status

setup:
  go get -u github.com/kardianos/govendor
  go get github.com/golang/lint/golint
  govendor sync

ci: lint test

build:
	go build -o wchy-api main.go

run: build
	./wchy-api

.DEFAULT_GOAL := build