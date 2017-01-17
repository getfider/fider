build:
	go build -o wchy-api main.go

run:
	./wchy-api

.DEFAULT_GOAL := build