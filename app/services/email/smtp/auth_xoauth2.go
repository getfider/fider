package smtp

import (
	"fmt"
	gosmtp "net/smtp"

	"github.com/getfider/fider/app/pkg/errors"
)

type xoauth2Auth struct {
	user  string
	token string
	host  string
}

func XOAuth2Auth(user, token, host string) gosmtp.Auth {
	return &xoauth2Auth{
		user:  user,
		token: token,
		host:  host,
	}
}

func (a *xoauth2Auth) Start(server *gosmtp.ServerInfo) (proto string, toServer []byte, err error) {
	if server.Name != a.host {
		return "", nil, errors.New("smtp: wrong host name")
	}

	if !server.TLS {
	    return "", nil, errors.New("smtp: XOAUTH2 requires TLS")
	}

	resp := fmt.Sprintf("user=%s\x01auth=Bearer %s\x01\x01", a.user, a.token)
	return "XOAUTH2", []byte(resp), nil
}

func (a *xoauth2Auth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return nil, errors.New("smtp: unexpected server challenge for XOAUTH2")
	}
	return nil, nil
}
