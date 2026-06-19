package actions_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
)

func withBillingEnabled(t *testing.T) {
	t.Helper()
	originalKey := env.Config.Stripe.SecretKey
	originalMode := env.Config.HostMode
	env.Config.Stripe.SecretKey = "sk_test"
	env.Config.HostMode = "multi"
	t.Cleanup(func() {
		env.Config.Stripe.SecretKey = originalKey
		env.Config.HostMode = originalMode
	})
}

func ctxWithTenant(tenant *entity.Tenant) context.Context {
	return context.WithValue(context.Background(), app.TenantCtxKey, tenant)
}

func TestInviteUsers_Empty(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "subject", "message", "recipients")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_Oversized(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{
		Subject:    "Join us and share your ideas. Because we have a cool website and this subject needs to be very long",
		Message:    "Use this link to join %invite%",
		Recipients: []string{"jon.snow@got.com"},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "subject")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_MissingInvite(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{
		Subject:    "Share your feedback.",
		Message:    "Please!",
		Recipients: []string{"jon.snow@got.com"},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "message")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_TooManyRecipients(t *testing.T) {
	RegisterT(t)

	recipients := make([]string, 31)
	for i := 0; i < len(recipients); i++ {
		recipients[i] = fmt.Sprintf("jon.snow%d@got.com", i)
	}

	action := &actions.InviteUsers{
		Subject:    "Share your feedback.",
		Message:    "Use this link to join %invite%",
		Recipients: recipients,
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "recipients")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_InvalidRecipient(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"jon.snow",
			"@got.com",
		},
	}
	result := action.Validate(context.Background(), nil)
	ExpectFailed(result, "recipients")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_Valid(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	action := &actions.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"",
			"jon.snow@got.com",
			"arya.stark@got.com",
		},
	}

	ExpectSuccess(action.Validate(context.Background(), nil))

	Expect(action.Invitations).HasLen(2)

	Expect(action.Invitations[0].Email).Equals("jon.snow@got.com")
	Expect(action.Invitations[0].VerificationKey).IsNotEmpty()

	Expect(action.Invitations[1].Email).Equals("arya.stark@got.com")
	Expect(action.Invitations[1].VerificationKey).IsNotEmpty()
}

func TestInviteUsers_IgnoreAlreadyRegistered(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == "tony.stark@avengers.com" {
			q.Result = &entity.User{Email: q.Email}
			return nil
		}
		return app.ErrNotFound
	})

	action := &actions.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"tony.stark@avengers.com",
			"jon.snow@got.com",
			"arya.stark@got.com",
		},
	}

	ExpectSuccess(action.Validate(context.Background(), nil))

	Expect(action.Invitations).HasLen(2)

	Expect(action.Invitations[0].Email).Equals("jon.snow@got.com")
	Expect(action.Invitations[0].VerificationKey).IsNotEmpty()

	Expect(action.Invitations[1].Email).Equals("arya.stark@got.com")
	Expect(action.Invitations[1].VerificationKey).IsNotEmpty()
}

func TestInviteUsers_ShouldFail_WhenAllRecipientsIgnored(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		q.Result = &entity.User{Email: q.Email}
		return nil
	})

	action := &actions.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"tony.stark@avengers.com",
		},
	}

	ExpectFailed(action.Validate(context.Background(), nil), "recipients")
}

func TestInviteUsers_SampleInvite_IgnoreRecipients(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{
		IsSampleInvite: true,
		Subject:        "Share your feedback.",
		Message:        "Use this link to join our community: %invite%",
	}

	ExpectSuccess(action.Validate(context.Background(), nil))
}

func TestInviteUsers_NonProTenant_OverridesCopy(t *testing.T) {
	RegisterT(t)
	withBillingEnabled(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	tenant := &entity.Tenant{Name: "Acme", IsPro: false}
	inviter := &entity.User{Name: "Jon Snow"}
	action := &actions.InviteUsers{
		Subject:    "Free Vinted vouchers!! Click here",
		Message:    "Tap %invite% to claim your prize",
		Recipients: []string{"victim@example.com"},
	}

	ExpectSuccess(action.Validate(ctxWithTenant(tenant), inviter))

	Expect(action.Subject).Equals(actions.DefaultInviteSubject(tenant))
	Expect(action.Message).Equals(actions.DefaultInviteMessage(tenant, inviter))
}

func TestInviteUsers_ProTenant_KeepsCustomCopy(t *testing.T) {
	RegisterT(t)
	withBillingEnabled(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	tenant := &entity.Tenant{Name: "Acme", IsPro: true}
	inviter := &entity.User{Name: "Jon Snow"}
	customSubject := "Share your feedback."
	customMessage := "Use this link to join our community: %invite%"
	action := &actions.InviteUsers{
		Subject:    customSubject,
		Message:    customMessage,
		Recipients: []string{"customer@example.com"},
	}

	ExpectSuccess(action.Validate(ctxWithTenant(tenant), inviter))

	Expect(action.Subject).Equals(customSubject)
	Expect(action.Message).Equals(customMessage)
}

func TestInviteUsers_BillingDisabled_KeepsCustomCopy(t *testing.T) {
	RegisterT(t)
	// Billing disabled: self-hosted Fider, no Stripe key configured.
	originalKey := env.Config.Stripe.SecretKey
	env.Config.Stripe.SecretKey = ""
	t.Cleanup(func() { env.Config.Stripe.SecretKey = originalKey })

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	tenant := &entity.Tenant{Name: "Acme", IsPro: false}
	inviter := &entity.User{Name: "Jon Snow"}
	customSubject := "Share your feedback."
	customMessage := "Use this link to join our community: %invite%"
	action := &actions.InviteUsers{
		Subject:    customSubject,
		Message:    customMessage,
		Recipients: []string{"customer@example.com"},
	}

	ExpectSuccess(action.Validate(ctxWithTenant(tenant), inviter))

	Expect(action.Subject).Equals(customSubject)
	Expect(action.Message).Equals(customMessage)
}

func TestInviteUsers_DefaultCopyPassesValidation(t *testing.T) {
	RegisterT(t)

	tenant := &entity.Tenant{Name: "Acme"}
	inviter := &entity.User{Name: "Jon Snow"}

	subject := actions.DefaultInviteSubject(tenant)
	message := actions.DefaultInviteMessage(tenant, inviter)

	Expect(len(subject) <= 70).IsTrue()
	Expect(message).ContainsSubstring(app.InvitePlaceholder)
}
