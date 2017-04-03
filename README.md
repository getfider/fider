[![Build](https://travis-ci.org/WeCanHearYou/wechy.svg?branch=master)](https://travis-ci.org/WeCanHearYou/wechy)
[![codecov](https://codecov.io/gh/WeCanHearYou/wechy/branch/master/graph/badge.svg)](https://codecov.io/gh/WeCanHearYou/wechy)


# What is WeCanHearYou?

Visit http://we.canhearyou.com for information on what it is and how to use it.

# How to run it locally?

WeCanHearYou is mainly written in Go and TypeScript, but we're also using things like Node.js, React and PostgreSQL. So if you know these technologies or would like to learn, you came to the right place.

Tools you'll need:

- Go 1.8+ (https://golang.org)
- Node.js 6+ (https://nodejs.org/)
- Yarn (https://yarnpkg.com/)
- Docker (https://www.docker.com/)

Step by step:

1) clone this repository into `$GOPATH/src/GitHub.com/WeCanHearYou/wechy`
2) run `yarn` to install froentend packages 
3) run `npm run build:watch` to pack the froentend source into a bundle. It'll 
4) run `docker-compose up -d` to start a local PostgreSQL database on Docker.
5) run `make watch` to start the application.
6) Navigate to `http://orange.dev.canhearyou.com:3000` and boom! Welcome to your new local WeCanHearYou development copy!

Change some code, fix some bugs, implement some features and send us your Pull Request!
