BUILD_TIME=$(shell date +"%Y.%m.%d.%H%M%S")

ifeq ($(WERCKER), true)
ENV_FILE=.ci.env
else
ENV_FILE=.test.env
endif

test:
	godotenv -f ${ENV_FILE} go test ./... -cover -p=1

coverage:
	echo 'mode: atomic' > cover.out
	go list ./... | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverpkg=./... -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> cover.out'
	rm coverage.tmp

build:
	go build -a -ldflags='-X main.buildtime=${BUILD_TIME}' -o fider .

watch:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'"

watch-ssl:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'" --certFile etc/server.crt --keyFile etc/server.key

run:
	godotenv -f .env ./fider

.DEFAULT_GOAL := build
