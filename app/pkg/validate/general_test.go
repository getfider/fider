package validate_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
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
	} {
		messages := validate.WebhookURL(rawurl)
		Expect(messages).HasLen(0)
	}
}

func TestEmailNotDisposable_ToggleOff(t *testing.T) {
	RegisterT(t)

	// Toggle is off: always returns nil regardless of the email domain
	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{BlockDisposableEmails: false})
	Expect(validate.EmailNotDisposable(ctx, "foo@mailinator.com")).HasLen(0)
	Expect(validate.EmailNotDisposable(ctx, "foo@gmail.com")).HasLen(0)
}

func TestEmailNotDisposable_NoTenant(t *testing.T) {
	RegisterT(t)

	// No tenant in context: always returns nil
	Expect(validate.EmailNotDisposable(context.Background(), "foo@mailinator.com")).HasLen(0)
}

func TestEmailNotDisposable_ToggleOn_Blocks(t *testing.T) {
	RegisterT(t)

	// Register a stub handler that returns no tenant-level rules (bundled list only)
	bus.AddHandler(func(ctx context.Context, q *query.GetEmailDomainRules) error {
		return nil
	})

	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{BlockDisposableEmails: true})

	// mailinator.com is in the bundled list, should be blocked
	messages := validate.EmailNotDisposable(ctx, "foo@mailinator.com")
	Expect(len(messages) > 0).IsTrue()

	// gmail.com is not in the bundled list, should be allowed
	Expect(validate.EmailNotDisposable(ctx, "foo@gmail.com")).HasLen(0)
}

func TestEmailNotDisposable_ToggleOn_TenantAllowOverrides(t *testing.T) {
	RegisterT(t)

	// Tenant allow rule overrides the bundled deny list
	bus.AddHandler(func(ctx context.Context, q *query.GetEmailDomainRules) error {
		q.Result.Allow = []*entity.EmailDomainRule{{Domain: "mailinator.com"}}
		return nil
	})

	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{BlockDisposableEmails: true})

	// mailinator.com is in the bundled list but also in tenant allow, so it should pass
	Expect(validate.EmailNotDisposable(ctx, "foo@mailinator.com")).HasLen(0)
}

func TestEmailNotDisposable_ToggleOn_TenantDenyBlocks(t *testing.T) {
	RegisterT(t)

	// Tenant deny rule blocks a domain not in the bundled list
	bus.AddHandler(func(ctx context.Context, q *query.GetEmailDomainRules) error {
		q.Result.Deny = []*entity.EmailDomainRule{{Domain: "customdisposable.example"}}
		return nil
	})

	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{BlockDisposableEmails: true})

	messages := validate.EmailNotDisposable(ctx, "foo@customdisposable.example")
	Expect(len(messages) > 0).IsTrue()
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
