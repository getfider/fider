package smtp

import (
	gosmtp "net/smtp"

	"github.com/getfider/fider/app/pkg/errors"
)

type agnosticAuth struct {
	identity, username, password, host string
	auth                               string
}

// AgnosticAuth returns an Auth that match the correct authentication
// thanks to gosmtp.ServerInfo
func AgnosticAuth(identity, username, password, host string) gosmtp.Auth {
	return &agnosticAuth{identity, username, password, host, ""}
}

func (a *agnosticAuth) Start(server *gosmtp.ServerInfo) (string, []byte, error) {
	for _, auth := range server.Auth {
		if a.auth = auth; a.auth == "LOGIN" {
			return LoginAuth(a.username, a.password).Start(server)
		} else if a.auth == "PLAIN" {
			return gosmtp.PlainAuth(a.identity, a.username, a.password, a.host).Start(server)
		} else if a.auth == "CRAM-MD5" {
			return gosmtp.CRAMMD5Auth(a.username, a.password).Start(server)
		}
	}
	return "", nil, errors.New("unknown auth method")
}

func (a *agnosticAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if a.auth == "LOGIN" {
		return LoginAuth(a.username, a.password).Next(fromServer, more)
	} else if a.auth == "PLAIN" {
		return gosmtp.PlainAuth(a.identity, a.username, a.password, a.host).Next(fromServer, more)
	} else if a.auth == "CRAM-MD5" {
		return gosmtp.CRAMMD5Auth(a.username, a.password).Next(fromServer, more)
	}
	return nil, errors.New("unknown auth method")
}
