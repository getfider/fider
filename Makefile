## This is a self-documented Makefile. For usage information, run `make help`:
##
## For more information, refer to https://suva.sh/posts/well-documented-makefiles/

LDFLAGS += -X github.com/getfider/fider/app/pkg/env.buildnumber=${BUILDNUMBER}



##@ Running

run: ## Run Fider
	godotenv -f .env ./fider

migrate: ## Run all database migrations
	godotenv -f .env ./fider migrate



##@ Building

build: build-server build-ssr build-ui ## Build server and ui

build-server: ## Build server
	go build -ldflags '-s -w $(LDFLAGS)' -o fider .

build-ui: ## Build all UI assets
	NODE_ENV=production npx webpack-cli

build-ssr: ## Build SSR script and locales
	npx lingui extract public/
	npx lingui compile
	NODE_ENV=production node esbuild.config.js



##@ Testing

test: test-server test-ui ## Test server and ui code

test-server: build-server build-ssr ## Run all server tests
	godotenv -f .test.env ./fider migrate
	godotenv -f .test.env go test ./... -race

test-ui: ## Run all UI tests
	TZ=GMT npx jest ./public

coverage-server: build-server build-ssr ## Run all server tests (with code coverage)
	godotenv -f .test.env ./fider migrate
	godotenv -f .test.env go test ./... -coverprofile=cover.out -coverpkg=all -p=8 -race



##@ Running (Watch Mode)

watch:
	make -j4 watch-server watch-ui

watch-server: migrate ## Build and run server in watch mode
	air -c air.conf

watch-ui: ## Build and run server in watch mode
	npx webpack-cli -w



##@ Linting

lint: lint-server lint-ui ## Lint server and ui

lint-server: ## Lint server code
	golangci-lint run

lint-ui: ## Lint ui code
	npx eslint .



##@ Miscellaneous

clean: ## Remove all build-generated content
	rm -rf ./dist
	rm -f ssr.js

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
