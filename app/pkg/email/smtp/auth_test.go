package smtp_test

import (
	gosmtp "net/smtp"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/email/smtp"
)

func TestAgnosticAuth_Login(t *testing.T) {
	RegisterT(t)

	auth := smtp.AgnosticAuth("", "jon", "s3cr3t", "test.fider.io")
	proto, bytes, err := auth.Start(&gosmtp.ServerInfo{
		Auth: []string{"LOGIN"},
	})

	Expect(err).IsNil()
	Expect(proto).Equals("LOGIN")
	Expect(bytes).Equals([]byte("jon"))

	bytes, err = auth.Next([]byte("Username:"), true)
	Expect(err).IsNil()
	Expect(bytes).Equals([]byte("jon"))

	bytes, err = auth.Next([]byte("Password:"), true)
	Expect(err).IsNil()
	Expect(bytes).Equals([]byte("s3cr3t"))
}

func TestAgnosticAuth_NoMatchingAuth(t *testing.T) {
	RegisterT(t)

	auth := smtp.AgnosticAuth("", "jon", "s3cr3t", "test.fider.io")
	proto, bytes, err := auth.Start(&gosmtp.ServerInfo{
		Auth: []string{"FAKE-MD5"},
	})

	Expect(err).IsNotNil()
	Expect(proto).IsEmpty()
	Expect(bytes).IsNil()
}
