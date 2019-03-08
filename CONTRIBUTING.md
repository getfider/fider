# Contributing

There are many ways you can contribute to Fider.

- **Send us a Pull Request** on GitHub. Make sure you read our [Getting Started](#getting-started-with-fider-codebase) guide to learn how to setup the development environment;
- **Report issues** and bug reports on https://github.com/getfider/fider/issues;
- **Give feedback** and vote on features you'd like to see at https://feedback.fider.io;
- **Spread the word** by starring us on GitHub. Tweet about the project and show it to your friends. The more people know about Fider, the bigger the community will be and more contributions will be made;
- **Support us financially** by donating any amount to our [OpenCollective](https://opencollective.com/fider) and help us continue our activities;

## Getting started with Fider codebase

Before start working on something that you intend to send a Pull Request, make sure there's an [GitHub Issue](https://github.com/getfider/fider/issues) open for that. If you're working on something not tracked yet, please open a new Issue before the Pull Request. If you have any question or need any help, leave a comment on the issue and we'll try our best to help you.

Fider is written in Go (backend) and TypeScript (frontend), but we also use things like Node.js, React and PostgreSQL.
If you know these technologies or would like to learn them, lucky you! This is the right place!

#### 1. Install the following tools:

| Software  | How to install | What is it used for |
|---|---|---|
| Go 1.11+ | https://golang.org/ | To compile server side code |
| Node.js 10+ | https://nodejs.org/ or run `nvm use` if you have nvm installed | To compile TypeScript and bundle all the client side code |
| Docker | https://www.docker.com/ | To start local PostgreSQL instances |
| cosmtrek/air | `go get github.com/cosmtrek/air/cmd/air/` | Live reload for Go apps. When you change the code, it automatically recompiles the application |
| joho/godotenv | `go get github.com/joho/godotenv/cmd/godotenv/` | To load environment variables from a `.env` so that you don't have to change your machine environment variables |
| magefile/mage | `go get github.com/magefile/mage/` | A cross-platform Make alternative |

#### 2. To setup your development workspace:

1. if it's your first time using Go, ensure that have your `GOPATH` and `PATH` variables correctly setup. This guide can help you on that https://golang.org/doc/code.html#GOPATH
2. clone the repository into `$GOPATH/src/github.com/getfider/fider`.
3. run `npm install` to install client side packages .
4. run `docker-compose up -d` to start a local PostgreSQL database on Docker.
5. run `cp .example.env .env` to create a local environment configuration file.

- **Important:** Fider has a strong dependency on an email delivery service. You'll need to edit `.env` file and configure the `EMAIL_*` environment variables with your own SMTP server details. If you don't have an SMTP server, you can either sign up for a [Mailgun account](https://www.mailgun.com/) (it's Free) or sign up for a [Mailtrap account](https://mailtrap.io), which is a free SMTP mocking server. If you prefer not to setup an email service, keep an eye on the server logs. Sometimes it's necessary to navigate to some URLs that are only sent by email, but are also written to the logs.

#### 3. To start the application

1. run `mage watch` to start the application on watch mode. The application will be reloaded everytime a file is changed. Alternatively, it's also possible to start Fider by running `mage build` and `mage run`.
2. Navigate to `http://localhost:3000/` and 🎉! You should see the sign up page of Fider!

#### 4. To run the unit tests:

1. run `mage test` to run both UI and Server unit tests.

## Common Issues

#### 1. It doesn't work on Windows

This is a known [Issue #434](https://github.com/getfider/fider/issues/434). If you're a Windows user and want to contribute to Fider, please help us resolve this.

#### 2. godotenv: not found

This happens when godotenv was not (or incorrectly) installed. Install it by running `go get github.com/joho/godotenv/cmd/godotenv/`.

#### 3. mage watch throws 'too many open files' error

macOS has a small limit on how many files can be open at the same time. This limit is usually OK for most users, but developers tools usuage require a larger limit. [Learn how to resolve this](https://www.macobserver.com/tips/deep-dive/evade-macos-many-open-files-error-pushing-limits/).

#### 4. mage isn't found even after installing

Include %GOPATH%\bin in PATH environment variable and restart the terminal
