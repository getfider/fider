SHELL := /bin/bash

.PONEY: all run

build-web:
	npx webpack
migrate: 
	go run main.go migrate
run:
	go run main.go

