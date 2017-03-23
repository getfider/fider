# WeCanHearYou

[![Build Status](https://travis-ci.org/WeCanHearYou/wechy.svg?branch=dev)](https://travis-ci.org/WeCanHearYou/wechy) | [![Coverage Status](https://coveralls.io/repos/github/WeCanHearYou/wechy/badge.svg?branch=dev)](https://coveralls.io/github/WeCanHearYou/wechy?branch=dev)

# What is WeCanHearYou ?

Visit http://we.canhearyou.com for information on what it is and how to use it.

# How to contribute?

WeCanHearYou is coded in Go, React and TypeScript. We also use Node.js to run some development tools. So if you know these technologies or would like to, please feel free to send us some PR!

Tools you'll need: Go 1.8, Glide, Node.js 6+, Yarn, Docker and Git, of course :)

1) clone this repository into $GOPATH/src/GitHub.com/WeCanHearYou/wechy
2) run `yarn` to install froentend packages 
3) run `npm run build:watch` to pack the froentend source into a bundle. It'll 
4) run `docker-conpose up -d` to start a  local Postgres database on Docker.
5) run `make watch` to start the application.
6) Navigate to ...