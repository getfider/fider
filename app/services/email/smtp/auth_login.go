package smtp

import (
	gosmtp "net/smtp"
)

type loginAuth struct {
	username, password string
}

// LoginAuth returns an Auth that implements the LOGIN authentication
// mechanism as defined in Internet-Draft draft-murchison-sasl-login-00.
// The LOGIN mechanism is still used by some SMTP server.
func LoginAuth(username, password string) gosmtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *gosmtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
