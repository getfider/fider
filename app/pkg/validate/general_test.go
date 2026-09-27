package validate_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/validate"
)

func TestInvalidEmail(t *testing.T) {
	RegisterT(t)

	for _, email := range []string{
		"hello",
		"",
		"my@company",
		"my @company.com",
		"my@.company.com",
		"my+company.com",
		".my@company.com",
		"my@company@other.com",
		"my@company@other.com",
		rand.String(200) + "@gmail.com",
	} {
		messages := validate.Email(context.Background(), email)
		Expect(len(messages) > 0).IsTrue()
	}
}

func TestValidEmail(t *testing.T) {
	RegisterT(t)

	for _, email := range []string{
		"hello@company.com",
		"hello+alias@company.com",
		"abc@gmail.com",
	} {
		messages := validate.Email(context.Background(), email)
		Expect(messages).HasLen(0)
	}
}

func TestInvalidURL(t *testing.T) {
	RegisterT(t)

	for _, rawurl := range []string{
		"http//google.com",
		"google.com",
		"google",
		rand.String(301),
		"my@company",
	} {
		messages := validate.URL(context.Background(), rawurl)
		Expect(len(messages) > 0).IsTrue()
	}
}

func TestValidURL(t *testing.T) {
	RegisterT(t)

	for _, rawurl := range []string{
		"http://example.org",
		"https://example.org/oauth",
		"https://example.org/oauth?test=abc",
	} {
		messages := validate.URL(context.Background(), rawurl)
		Expect(messages).HasLen(0)
	}
}

func TestWebhookURL_BlockedAddresses(t *testing.T) {
	RegisterT(t)

	for _, rawurl := range []string{
		"http://localhost/hook",
		"http://localhost:8080/hook",
		"http://LOCALHOST/hook",
		"http://127.0.0.1/hook",
		"http://127.0.0.2/hook",
		"http://[::1]/hook",
		"http://10.0.0.1/hook",
		"http://10.255.255.255/hook",
		"http://172.16.0.1/hook",
		"http://172.31.255.255/hook",
		"http://192.168.1.1/hook",
		"http://192.168.0.100:8080/hook",
		"http://169.254.169.254/latest/meta-data/",
		"http://169.254.1.1/hook",
		"http://0.0.0.0/hook",
		"http://100.64.0.1/hook",
		"http://198.18.0.1/hook",
		"http://224.0.0.1/hook",
		"http://255.255.255.255/hook",
		"http://[::ffff:127.0.0.1]/hook",
		"http://[::ffff:169.254.169.254]/hook",
		"http://[fd00::1]/hook",
		"http://[fe80::1]/hook",
		"http://[ff02::1]/hook",
		"http://[64:ff9b::a9fe:a9fe]/latest/meta-data/",
		"http://[64:ff9b::7f00:1]/hook",
		"http://[64:ff9b:1::1]/hook",
		"http://[2002:a9fe:a9fe::1]/hook",
		"http://[2002:7f00:1::1]:8080/hook",
		"http://[2001:0:4136:e378:8000:63bf:80ff:fffe]/hook",
		"http://192.0.0.192/opc/v2/instance/",
		"http://[fec0::1]/hook",
		"http://[::ffff:0:a9fe:a9fe]/latest/meta-data/",
		"ftp://example.com/hook",
		"file:///etc/passwd",
		"gopher://example.com",
	} {
		messages := validate.WebhookURL(rawurl)
		Expect(len(messages) > 0).IsTrue()
	}
}

func TestWebhookURL_AllowedAddresses(t *testing.T) {
	RegisterT(t)

	for _, rawurl := range []string{
		"https://hooks.slack.com/services/T00/B00/xxx",
		"https://discord.com/api/webhooks/123/abc",
		"http://example.com/webhook",
		"https://203.0.113.1/hook",
		"https://[2001:4860:4860::8888]/hook",
		"https://[2002:0808:0808::1]/hook",
		"https://[64:ff9b::808:808]/hook",
	} {
		messages := validate.WebhookURL(rawurl)
		Expect(messages).HasLen(0)
	}
}

func TestWebhookURL_PrivateIPsAllowedWhenOptedIn(t *testing.T) {
	RegisterT(t)

	original := env.Config.AllowPrivateNetworkTargets
	env.Config.AllowPrivateNetworkTargets = true
	t.Cleanup(func() { env.Config.AllowPrivateNetworkTargets = original })

	for _, rawurl := range []string{
		"http://localhost/hook",
		"http://127.0.0.1/hook",
		"http://10.0.0.1/hook",
		"http://172.16.0.1/hook",
		"http://192.168.1.1:8080/hook",
		"http://[::1]/hook",
		"http://internal.lan/hook",
	} {
		messages := validate.WebhookURL(rawurl)
		Expect(messages).HasLen(0)
	}
}

func TestWebhookURL_OptInDoesNotBypassFormatValidation(t *testing.T) {
	RegisterT(t)

	original := env.Config.AllowPrivateNetworkTargets
	env.Config.AllowPrivateNetworkTargets = true
	t.Cleanup(func() { env.Config.AllowPrivateNetworkTargets = original })

	for _, rawurl := range []string{
		"ftp://example.com/hook",
		"file:///etc/passwd",
		"not a url at all",
		"",
	} {
		messages := validate.WebhookURL(rawurl)
		Expect(len(messages) > 0).IsTrue()
	}
}

func TestInvalidCNAME(t *testing.T) {
	RegisterT(t)

	for _, cname := range []string{
		"hello",
		"hellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello.com",
		"",
		"my",
		"name.com/abc",
		"feedback.test.fider.io",
		"test.fider.io",
		"@google.com",
	} {
		messages := validate.CNAME(context.Background(), cname)
		Expect(len(messages) > 0).IsTrue()
	}
}

func TestValidHostname(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.IsCNAMEAvailable) error {
		q.Result = true
		return nil
	})

	for _, cname := range []string{
		"google.com",
		"feedback.fider.io",
		"my.super.domain.com",
		"jon-snow.got.com",
		"got.com",
		"hi.m",
	} {
		messages := validate.CNAME(context.Background(), cname)
		Expect(messages).HasLen(0)
	}
}

func TestValidCNAME_Availability(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.IsCNAMEAvailable) error {
		q.Result = q.CNAME != "footbook.com" && q.CNAME != "fider.yourcompany.com" && q.CNAME != "feedback.newyork.com"
		return nil
	})

	for _, cname := range []string{
		"footbook.com",
		"fider.yourcompany.com",
		"feedback.newyork.com",
	} {
		messages := validate.CNAME(context.Background(), cname)
		Expect(len(messages) > 0).IsTrue()
	}

	for _, cname := range []string{
		"fider.footbook.com",
		"yourcompany.com",
		"anything.com",
	} {
		messages := validate.CNAME(context.Background(), cname)
		Expect(messages).HasLen(0)
	}
}
