BUILD_TIME=$(shell date +"%Y.%m.%d.%H%M%S")

ifeq ($(WERCKER), true)
ENV_FILE=.ci.env
else
ENV_FILE=.test.env
endif

test:
	godotenv -f ${ENV_FILE} go test $$(go list ./... | grep -v /vendor/) -cover -p=1

coverage:
	godotenv -f ${ENV_FILE} ./cover.sh $$(go list ./... | grep -v /vendor/)

build:
	go build -a -ldflags='-X main.buildtime=${BUILD_TIME}' -o fider .

watch:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'"

run:
	fider

.DEFAULT_GOAL := build
