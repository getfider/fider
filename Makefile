SHELL := /bin/bash

.PONEY: all run

build-web:
	npx webpack
run:
	source .dev.env
	go run main.go migrate
	go run main.go
