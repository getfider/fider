BUILD_TIME="$(shell date +"%Y.%m.%d.%H%M%S")"

define tag_docker
	@if [ "$(TRAVIS_BRANCH)" = "master" ]; then \
		docker tag $(1) $(1):stable; \
	fi
	@if [ "$(TRAVIS_BRANCH)" = "dev" ]; then \
		docker tag $(1) $(1):staging; \
	fi
endef

test:
	godotenv -f .test.env go test $$(go list ./... | grep -v /vendor/) -cover

setup-ci:
	go get github.com/kardianos/govendor
	go get github.com/joho/godotenv/cmd/godotenv
	govendor sync
	govendor install +vendor

goveralls:
	godotenv -f .test.env goveralls -service=travis-ci

dockerize:
ifeq ($(TRAVIS_PULL_REQUEST), false)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build
	docker build -f Dockerfile -t WeCanHearYou/wechy .
	$(call tag_docker, WeCanHearYou/wechy)
	docker push WeCanHearYou/wechy
endif

migrate:
ifeq ($(GO_ENV), development)
	migrate -url postgres://wechy:wechy-pw@localhost:5555/wechy?sslmode=disable -path ./migrations up
else
	migrate -url ${DATABASE_URL} -path ./migrations up
endif

build:
	go build -ldflags='-X main.buildtime=${BUILD_TIME}'

watch:
	gin --buildArgs "-ldflags='-X main.buildtime=${BUILD_TIME}'"

run:
	wechy

.DEFAULT_GOAL := build
