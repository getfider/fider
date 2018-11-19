package smtp_test

import (
	gosmtp "net/smtp"
	"regexp"
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/email/smtp"
	"github.com/getfider/fider/app/pkg/log/noop"
	"github.com/getfider/fider/app/pkg/worker"
)

type request struct {
	servername string
	auth       gosmtp.Auth
	from       string
	to         []string
	body       []byte
}

var logger = noop.NewLogger()
var sender = smtp.NewSender(logger, "localhost", "1234", "us3r", "p4ss")
var tenant = &models.Tenant{
	Subdomain: "got",
}

var ctx = worker.NewContext("ID-1", worker.Task{Name: "TaskName"}, nil, logger)

var requests = make([]request, 0)

func mockSend(localname, servername string, auth gosmtp.Auth, from string, to []string, body []byte) error {
	requests = append(requests, request{servername, auth, from, to, body})
	return nil
}

func reset() {
	ctx.SetTenant(tenant)
	sender.ReplaceSend(mockSend)
	requests = make([]request, 0)
}

func TestSend_Success(t *testing.T) {
	RegisterT(t)
	reset()

	to := email.Recipient{
		Name:    "Jon Sow",
		Address: "jon.snow@got.com",
	}
	err := sender.Send(ctx, "echo_test", email.Params{
		"name": "Hello",
	}, "Fider Test", to)

	Expect(err).IsNil()
	Expect(requests).HasLen(1)
	Expect(requests[0].servername).Equals("localhost:1234")
	Expect(requests[0].auth).Equals(smtp.AgnosticAuth("", "us3r", "p4ss", "localhost"))
	Expect(requests[0].from).Equals("noreply@random.org")
	Expect(requests[0].to).Equals([]string{"jon.snow@got.com"})
	Expect(string(requests[0].body)).ContainsSubstring("From: \"Fider Test\" <noreply@random.org>\r\nReply-To: noreply@random.org\r\nTo: \"Jon Sow\" <jon.snow@got.com>\r\nSubject: Message to: Hello\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nDate: ")
	Expect(string(requests[0].body)).ContainsSubstring("Message-ID: ")
	Expect(string(requests[0].body)).ContainsSubstring("Hello World Hello!")

	var validID = regexp.MustCompile(`.*Message-ID: <[a-z0-9\-].*\.[0-9].*@.*>.*`)
	Expect(validID.MatchString(string(requests[0].body))).IsTrue()
}
func TestSend_SkipEmptyAddress(t *testing.T) {
	RegisterT(t)
	reset()

	to := email.Recipient{
		Name:    "Jon Sow",
		Address: "",
	}
	err := sender.Send(ctx, "echo_test", email.Params{
		"name": "Hello",
	}, "Fider Test", to)

	Expect(err).IsNil()
	Expect(requests).HasLen(0)
}

func TestSend_SkipUnlistedAddress(t *testing.T) {
	RegisterT(t)
	reset()
	email.SetWhitelist("^.*@gmail.com$")

	to := email.Recipient{
		Name:    "Jon Sow",
		Address: "jon.snow@got.com",
	}
	err := sender.Send(ctx, "echo_test", email.Params{
		"name": "Hello",
	}, "Fider Test", to)

	Expect(err).IsNil()
	Expect(requests).HasLen(0)
}

func TestBatch_Success(t *testing.T) {
	RegisterT(t)
	reset()
	email.SetWhitelist("")

	to := []email.Recipient{
		email.Recipient{
			Name:    "Jon Sow",
			Address: "jon.snow@got.com",
			Params: email.Params{
				"name": "Jon",
			},
		},
		email.Recipient{
			Name:    "Arya Stark",
			Address: "arya.start@got.com",
			Params: email.Params{
				"name": "Arya",
			},
		},
	}

	err := sender.BatchSend(ctx, "echo_test", email.Params{}, "Fider Test", to)
	Expect(err).IsNil()

	Expect(requests).HasLen(2)

	Expect(requests[0].servername).Equals("localhost:1234")
	Expect(requests[0].auth).Equals(smtp.AgnosticAuth("", "us3r", "p4ss", "localhost"))
	Expect(requests[0].from).Equals("noreply@random.org")
	Expect(requests[0].to).Equals([]string{"jon.snow@got.com"})
	Expect(string(requests[0].body)).ContainsSubstring("From: \"Fider Test\" <noreply@random.org>\r\nReply-To: noreply@random.org\r\nTo: \"Jon Sow\" <jon.snow@got.com>\r\nSubject: Message to: Jon\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nDate: ")
	Expect(string(requests[0].body)).ContainsSubstring("Message-ID: ")
	Expect(string(requests[0].body)).ContainsSubstring("Hello World Jon!")

	Expect(requests[1].servername).Equals("localhost:1234")
	Expect(requests[1].auth).Equals(smtp.AgnosticAuth("", "us3r", "p4ss", "localhost"))
	Expect(requests[1].from).Equals("noreply@random.org")
	Expect(requests[1].to).Equals([]string{"arya.start@got.com"})
	Expect(string(requests[1].body)).ContainsSubstring("From: \"Fider Test\" <noreply@random.org>\r\nReply-To: noreply@random.org\r\nTo: \"Arya Stark\" <arya.start@got.com>\r\nSubject: Message to: Arya\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nDate: ")
	Expect(string(requests[1].body)).ContainsSubstring("Message-ID: ")
	Expect(string(requests[1].body)).ContainsSubstring("Hello World Arya!")
}
