package smtp

import (
	gosmtp "net/smtp"
)

type agnosticAuth struct {
	identity, username, password, host string
	auth                               gosmtp.Auth
}

func (a *agnosticAuth) FindAuth(server *gosmtp.ServerInfo) gosmtp.Auth {
	for _, auth := range server.Auth {
		switch auth {
		case "LOGIN":
			return LoginAuth(a.username, a.password)
		case "PLAIN":
			return gosmtp.PlainAuth(a.identity, a.username, a.password, a.host)
		case "CRAM-MD5":
			return gosmtp.CRAMMD5Auth(a.username, a.password)
		default:
			continue
		}
	}
	return nil
}

// AgnosticAuth returns an Auth that match the correct authentication
// thanks to gosmtp.ServerInfo
func AgnosticAuth(identity, username, password, host string) gosmtp.Auth {
	return &agnosticAuth{identity, username, password, host, nil}
}

func (a *agnosticAuth) Start(server *gosmtp.ServerInfo) (string, []byte, error) {
	a.auth = a.FindAuth(server)
	return a.auth.Start(server)
}

func (a *agnosticAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	return a.auth.Next(fromServer, more)
}
