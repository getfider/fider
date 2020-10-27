SHELL := /bin/bash

.PONEY: all run tslint

build-web:
	npx webpack
migrate: 
	go run main.go migrate
run:
	go run main.go
tslint:
	npx tslint  'public/**/*.{ts,tsx}' 'tests/**/*.{ts,tsx}'
