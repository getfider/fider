[![Build Status](https://travis-ci.org/WeCanHearYou/wechy.svg?branch=master)](https://travis-ci.org/WeCanHearYou/wechy)
[![Coverage Status](https://coveralls.io/repos/github/WeCanHearYou/wechy/badge.svg?branch=master)](https://coveralls.io/github/WeCanHearYou/wechy?branch=master)

# What is WeCanHearYou?

Visit http://we.canhearyou.com for information on what it is and how to use it.

# How to run it locally?

WeCanHearYou is coded in Go and TypeScript, but also using things like Node.js, React and Postgres. So if you know these technologies or would like to learn, you came to the right place.

Tools you'll need: Go 1.8, Glide, Node.js 6+, Yarn, Docker and Git :)

1) clone this repository into `$GOPATH/src/GitHub.com/WeCanHearYou/wechy`
2) run `yarn` to install froentend packages 
3) run `npm run build:watch` to pack the froentend source into a bundle. It'll 
4) run `docker-conpose up -d` to start a  local Postgres database on Docker.
5) run `make watch` to start the application.
6) Navigate to `http://orange.dev.canhearyou.com:3000` and boom! Welcome to your new local WeCanHearYou development copy!

Change some code, fix some bugs, implement some features and send us your Pull Request!