# Contributing to Fider

## Ways to contribute

- **Send us a Pull Request** on GitHub. Make sure you read our [Getting Started](#getting-started-with-fider-codebase) guide.
- **Reporting issues** and bug reports on (https://github.com/getfider/fider).
- **Giving feedback** and voting on features you'd like to see at https://feedback.fider.io.
- **Spread the word!** Star us on GitHub, tweet about us, show it to your friends. The more people are aware of Fider, the bigger the community and the contributions will be.

## Getting started with Fider codebase

Before youstart working on something that you intend to send a Pull Request back to Fider, make sure there's an [GitHub Issue](https://github.com/getfider/fider/issues) open for that. If you're working on something not tracked yet, please open a new Issue before the Pull Request.

Fider is mainly written in Go and TypeScript, but we're also using things like Node.js, React and PostgreSQL. 
If you know these technologies or would like to learn them, lucky you! This is the right place!

Install the following tools:

- Go 1.10 (https://golang.org/)
- Node.js 8.x (https://nodejs.org/)
- Docker (https://www.docker.com/)
- github.com/codegangsta/gin (https://github.com/codegangsta/gin/)

To setup your development workspace:

1) clone the repository into `$GOPATH/src/github.com/getfider/fider`
2) run `npm install` to install front end packages 
3) run `docker-compose up -d pgdev` to start a local PostgreSQL database on Docker.
4) run `cp .example.env .env` to create a local environment configuration file.
5) run `make watch` to start the application.
6) Navigate to `http://localhost:3000/` and 🎉! You should see the sign up page of Fider!

To run the unit tests:

1) run `docker-compose up -d pgtest` to start a test-only PostgreSQL database on Docker.
2) run `make test`.
