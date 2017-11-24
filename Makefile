BUILD_TIME=$(shell date +"%Y.%m.%d.%H%M%S")

ENV_FILE=.test.env

ifeq ($(TRAVIS), true)
ENV_FILE=.ci.env
endif

ifeq ($(TRAVIS_BRANCH), master)
DOCKER_TAG=master
endif	

ifdef TRAVIS_PULL_REQUEST
ifneq ($(TRAVIS_PULL_REQUEST), false)
DOCKER_TAG=$(TRAVIS_PULL_REQUEST)
endif	
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

dockerize:
	docker build -t getfider/fider .
	docker login -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}"
	docker tag getfider/fider getfider/fider:${DOCKER_TAG}
	docker push getfider/fider:${DOCKER_TAG}

run:
	godotenv -f .env ./fider

.DEFAULT_GOAL := build
