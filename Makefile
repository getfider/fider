test:
	go test

build:
	go build -o wchy-api main.go

run: build
	./wchy-api

.DEFAULT_GOAL := build