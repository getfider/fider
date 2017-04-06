BUILD_TIME=$(shell date +"%Y.%m.%d.%H%M%S")

define tag_docker
	@if [ "$(TRAVIS_BRANCH)" = "master" ]; then \
		docker tag $(1) $(1):staging; \
	fi
endef

# remove this. check if test is running on WERCKER, if YES, then change .env file
ci-test:
	godotenv -f .ci.env go test $$(go list ./... | grep -v /vendor/) -cover -p=1

test:
	godotenv -f .test.env go test $$(go list ./... | grep -v /vendor/) -cover -p=1

dockerize:
ifeq ($(TRAVIS_PULL_REQUEST), false)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
	docker build -f Dockerfile -t wecanhearyou/wechy .
	$(call tag_docker, wecanhearyou/wechy)
	docker push wecanhearyou/wechy
endif

coverage:
	rm -rf coverage.txt
	for d in $$(go list ./... | grep -v vendor); do \
			godotenv -f .test.env go test -p=1 -race -coverprofile=profile.out -covermode=atomic $$d ; \
			if [ -f profile.out ]; then \
					cat profile.out >> coverage.txt ; \
					rm profile.out ; \
			fi \
	done

build:
	go build -ldflags='-X main.buildtime=${BUILD_TIME}'

watch:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'"

run:
	wechy

.DEFAULT_GOAL := build
