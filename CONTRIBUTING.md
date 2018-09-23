# Contributing

There are many ways you can contribute to Fider.

- **Send us a Pull Request** on GitHub. Make sure you read our [Getting Started](#getting-started-with-fider-codebase) guide to learn how to setup the development environment;
- **Report issues** and bug reports on https://github.com/getfider/fider/issues;
- **Give feedback** and vote on features you'd like to see at https://feedback.fider.io;
- **Spread the word** by starring us on GitHub. Tweet about the project and show it to your friends. The more people know about Fider, the bigger the community will be and more contributions will be made;
- **Support us financially** by donating any amount to our [OpenCollective](https://opencollective.com/fider) and help us continue our activities;

## Getting started with Fider codebase

<<<<<<< HEAD
Before start working on something that you intend to send a Pull Request, make sure there's an [GitHub Issue](https://github.com/getfider/fider/issues) open for that. If you're working on something not tracked yet, please open a new Issue before the Pull Request. If you have any question or need any help, leave a comment on the issue and we'll try our best to help you.
=======
Before you start working on something that you intend to send a Pull Request to Fider, make sure there's an [GitHub Issue](https://github.com/getfider/fider/issues) open for that. If you're working on something that is not yet tracked, please open a new Issue before the Pull Request.
>>>>>>> b6a43a63cc3410793abdd19a088a2bd09e9e6397

Fider is written in Go (backend) and TypeScript (frontend), but we also use things like Node.js, React and PostgreSQL.
If you know these technologies or would like to learn them, lucky you! This is the right place!

Install the following tools:

- Go 1.10+ (https://golang.org/)
- Node.js 10+ (https://nodejs.org/)
- Docker (https://www.docker.com/)
- cosmtrek/air (https://github.com/cosmtrek/air/)

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
