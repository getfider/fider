package smtp_test

import (
	"context"
	gosmtp "net/smtp"
	"regexp"
	"testing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/services/email"
	"github.com/getfider/fider/app/services/email/smtp"
)

type request struct {
	servername string
	auth       gosmtp.Auth
	from       string
	to         []string
	body       []byte
}

var ctx context.Context

var requests = make([]request, 0)

func mockSend(localname, servername string, enableStartTLS bool, auth gosmtp.Auth, from string, to []string, body []byte) error {
	requests = append(requests, request{servername, auth, from, to, body})
	return nil
}

func reset() {
	ctx = context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{
		Subdomain: "got",
	})
	smtp.Send = mockSend
	requests = make([]request, 0)
	bus.Init(smtp.Service{})
}

func TestSend_Success(t *testing.T) {
	RegisterT(t)
	reset()

	bus.Publish(ctx, &cmd.SendMail{
		From: dto.Recipient{Name: "Fider Test"},
		To: []dto.Recipient{
			{
				Name:    "Jon Sow",
				Address: "jon.snow@got.com",
			},
		},
		TemplateName: "echo_test",
		Props: dto.Props{
			"name": "Hello",
		},
	})

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

	bus.Publish(ctx, &cmd.SendMail{
		From: dto.Recipient{Name: "Fider Test"},
		To: []dto.Recipient{
			{
				Name:    "Jon Sow",
				Address: "",
			},
		},
		TemplateName: "echo_test",
		Props: dto.Props{
			"name": "Hello",
		},
	})

	Expect(requests).HasLen(0)
}

func TestSend_SkipUnlistedAddress(t *testing.T) {
	RegisterT(t)
	reset()
	email.SetAllowlist("^.*@gmail.com$")

	bus.Publish(ctx, &cmd.SendMail{
		From: dto.Recipient{Name: "Fider Test"},
		To: []dto.Recipient{
			{
				Name:    "Jon Sow",
				Address: "jon.snow@got.com",
			},
		},
		TemplateName: "echo_test",
		Props: dto.Props{
			"name": "Hello",
		},
	})

	Expect(requests).HasLen(0)
}

func TestBatch_Success(t *testing.T) {
	RegisterT(t)
	reset()
	email.SetAllowlist("")

	bus.Publish(ctx, &cmd.SendMail{
		From: dto.Recipient{Name: "Fider Test"},
		To: []dto.Recipient{
			{
				Name:    "Jon Sow",
				Address: "jon.snow@got.com",
				Props: dto.Props{
					"name": "Jon",
				},
			},
			{
				Name:    "Arya Stark",
				Address: "arya.start@got.com",
				Props: dto.Props{
					"name": "Arya",
				},
			},
		},
		TemplateName: "echo_test",
	})

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
