BUILD_TIME=$(shell date +"%Y.%m.%d.%H%M%S")

ENV_FILE=.test.env

ifeq ($(WERCKER), true)
ENV_FILE=.ci.env
endif
ifeq ($(TRAVIS), true)
ENV_FILE=.ci.env
endif

test:
	godotenv -f ${ENV_FILE} go test ./... -cover -p=1

coverage:
	godotenv -f ${ENV_FILE} courtney -o cover.out $$(go list ./...)

build:
	go build -a -ldflags='-X main.buildtime=${BUILD_TIME}' -o fider .

watch:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'"

watch-ssl:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'" --certFile etc/server.crt --keyFile etc/server.key

run:
	godotenv -f .env ./fider

.DEFAULT_GOAL := build
