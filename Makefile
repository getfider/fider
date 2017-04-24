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

old-coverage:
	rm -rf coverage.txt
	for d in $$(go list ./... | grep -v vendor); do \
			godotenv -f ${ENV_FILE} go test -p=1 -race -coverprofile=profile.out -coverpkg=github.com/WeCanHearYou/wechy/app/... -covermode=atomic $$d ; \
			if [ -f profile.out ]; then \
					cat profile.out >> coverage.txt ; \
					rm profile.out ; \
			fi \
	done

build:
	go build -a -ldflags='-X main.buildtime=${BUILD_TIME}' -o wechy .

watch:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'"

run:
	wechy

.DEFAULT_GOAL := build
