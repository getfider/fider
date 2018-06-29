package mailgun_test

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/mock"

	"github.com/getfider/fider/app/pkg/email/mailgun"
	"github.com/getfider/fider/app/pkg/log/noop"

	. "github.com/getfider/fider/app/pkg/assert"
)

var client = mock.NewHTTPClient()
var sender = mailgun.NewSender(noop.NewLogger(), client, "mydomain.com", "mys3cr3tk3y")
var tenant = &models.Tenant{
	Subdomain: "got",
}

func TestSend_Success(t *testing.T) {
	RegisterT(t)
	os.Setenv("HOST_MODE", "MULTI")
	client.Reset()

	to := email.Recipient{
		Name:    "Jon Sow",
		Address: "jon.snow@got.com",
	}
	sender.Send(tenant, "echo_test", email.Params{
		"name": "Hello",
	}, "Fider Test", to)

	Expect(client.Requests).HasLen(1)
	Expect(client.Requests[0].URL.String()).Equals("https://api.mailgun.net/v3/mydomain.com/messages")
	Expect(client.Requests[0].Header.Get("Authorization")).Equals("Basic YXBpOm15czNjcjN0azN5")
	Expect(client.Requests[0].Header.Get("Content-Type")).Equals("application/x-www-form-urlencoded")

	bytes, err := ioutil.ReadAll(client.Requests[0].Body)
	Expect(err).IsNil()
	values, err := url.ParseQuery(string(bytes))
	Expect(err).IsNil()
	Expect(values).HasLen(5)
	Expect(values.Get("to")).Equals("Jon Sow <jon.snow@got.com>")
	Expect(values.Get("from")).Equals("Fider Test <noreply@random.org>")
	Expect(values.Get("subject")).Equals("Message to: Hello")
	Expect(values.Get("html")).Equals("Hello World Hello!")
	Expect(values["o:tag"][0]).Equals("template:echo_test")
	Expect(values["o:tag"][1]).Equals("tenant:got")
}

func TestSend_SkipEmptyAddress(t *testing.T) {
	RegisterT(t)
	client.Reset()

	to := email.Recipient{
		Name:    "Jon Sow",
		Address: "",
	}
	sender.Send(tenant, "echo_test", email.Params{
		"name": "Hello",
	}, "Fider Test", to)

	Expect(client.Requests).HasLen(0)
}

func TestSend_SkipUnlistedAddress(t *testing.T) {
	RegisterT(t)
	client.Reset()
	email.SetWhitelist("^.*@gmail.com$")

	to := email.Recipient{
		Name:    "Jon Sow",
		Address: "jon.snow@got.com",
	}
	sender.Send(tenant, "echo_test", email.Params{
		"name": "Hello",
	}, "Fider Test", to)

	Expect(client.Requests).HasLen(0)
}

func TestBatch_Success(t *testing.T) {
	RegisterT(t)
	client.Reset()
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
	sender.BatchSend(tenant, "echo_test", email.Params{}, "Fider Test", to)

	Expect(client.Requests).HasLen(1)
	Expect(client.Requests[0].URL.String()).Equals("https://api.mailgun.net/v3/mydomain.com/messages")
	Expect(client.Requests[0].Header.Get("Authorization")).Equals("Basic YXBpOm15czNjcjN0azN5")
	Expect(client.Requests[0].Header.Get("Content-Type")).Equals("application/x-www-form-urlencoded")

	bytes, err := ioutil.ReadAll(client.Requests[0].Body)
	Expect(err).IsNil()
	values, err := url.ParseQuery(string(bytes))
	Expect(err).IsNil()
	Expect(values).HasLen(6)
	Expect(values["to"]).HasLen(2)
	Expect(values["to"][0]).Equals("Jon Sow <jon.snow@got.com>")
	Expect(values["to"][1]).Equals("Arya Stark <arya.start@got.com>")
	Expect(values.Get("from")).Equals("Fider Test <noreply@random.org>")
	Expect(values.Get("subject")).Equals("Message to: %recipient.name%")
	Expect(values.Get("html")).Equals("Hello World %recipient.name%!")
	Expect(values["o:tag"]).HasLen(2)
	Expect(values["o:tag"][0]).Equals("template:echo_test")
	Expect(values["o:tag"][1]).Equals("tenant:got")
	Expect(values.Get("recipient-variables")).Equals("{\"arya.start@got.com\":{\"name\":\"Arya\"},\"jon.snow@got.com\":{\"name\":\"Jon\"}}")
}
