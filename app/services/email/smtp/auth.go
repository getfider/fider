package smtp

import (
	gosmtp "net/smtp"

	"github.com/getfider/fider/app/pkg/errors"
)

type agnosticAuth struct {
	identity, username, password, host string
	auth                               gosmtp.Auth
}

func (a *agnosticAuth) createAuth(mode string) gosmtp.Auth {
	switch mode {
	case "LOGIN":
		return LoginAuth(a.username, a.password)
	case "PLAIN":
		return gosmtp.PlainAuth(a.identity, a.username, a.password, a.host)
	case "CRAM-MD5":
		return gosmtp.CRAMMD5Auth(a.username, a.password)
	default:
		return nil
	}
}

// AgnosticAuth returns an Auth that match the correct authentication
// thanks to gosmtp.ServerInfo
func AgnosticAuth(identity, username, password, host string) gosmtp.Auth {
	return &agnosticAuth{identity, username, password, host, nil}
}

func (a *agnosticAuth) Start(server *gosmtp.ServerInfo) (string, []byte, error) {
	for _, auth := range server.Auth {
		a.auth = a.createAuth(auth)
		if a.auth == nil {
			continue
		}

		proto, toServer, err := a.auth.Start(server)
		if err == nil {
			return proto, toServer, err
		}
	}

	return "", nil, errors.New("could not find any suitable auth mechanism")
}

func (a *agnosticAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	return a.auth.Next(fromServer, more)
}
