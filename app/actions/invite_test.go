package actions_test

import (
	"fmt"
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestInviteUsers_Empty(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{Model: &models.InviteUsers{}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "subject", "message", "recipients")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_Oversized(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject:    "Join us and share your ideas. Because we have a cool website and this subject needs to be very long",
		Message:    "Use this link to join %invite%",
		Recipients: []string{"jon.snow@got.com"},
	}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "subject")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_MissingInvite(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject:    "Share your feedback.",
		Message:    "Please!",
		Recipients: []string{"jon.snow@got.com"},
	}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "message")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_TooManyRecipients(t *testing.T) {
	RegisterT(t)

	recipients := make([]string, 31)
	for i := 0; i < len(recipients); i++ {
		recipients[i] = fmt.Sprintf("jon.snow%d@got.com", i)
	}

	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject:    "Share your feedback.",
		Message:    "Use this link to join %invite%",
		Recipients: recipients,
	}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "recipients")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_InvalidRecipient(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"jon.snow",
			"@got.com",
		},
	}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "recipients")
	Expect(action.Invitations).IsNil()
}

func TestInviteUsers_Valid(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"",
			"jon.snow@got.com",
			"arya.stark@got.com",
		},
	}}

	ExpectSuccess(action.Validate(nil, services))

	Expect(action.Invitations).HasLen(2)

	Expect(action.Invitations[0].Email).Equals("jon.snow@got.com")
	Expect(action.Invitations[0].VerificationKey).IsNotEmpty()

	Expect(action.Invitations[1].Email).Equals("arya.stark@got.com")
	Expect(action.Invitations[1].VerificationKey).IsNotEmpty()
}

func TestInviteUsers_IgnoreAlreadyRegistered(t *testing.T) {
	RegisterT(t)

	theTenant := &models.Tenant{ID: 1, Name: "The Tenant"}
	services.Users.SetCurrentTenant(theTenant)
	services.Users.Register(&models.User{
		Name:   "Tony",
		Email:  "tony.stark@avengers.com",
		Tenant: theTenant,
	})
	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"tony.stark@avengers.com",
			"jon.snow@got.com",
			"arya.stark@got.com",
		},
	}}

	ExpectSuccess(action.Validate(nil, services))

	Expect(action.Invitations).HasLen(2)

	Expect(action.Invitations[0].Email).Equals("jon.snow@got.com")
	Expect(action.Invitations[0].VerificationKey).IsNotEmpty()

	Expect(action.Invitations[1].Email).Equals("arya.stark@got.com")
	Expect(action.Invitations[1].VerificationKey).IsNotEmpty()
}

func TestInviteUsers_ShouldFail_WhenAllRecipientsIgnored(t *testing.T) {
	RegisterT(t)

	theTenant := &models.Tenant{ID: 1, Name: "The Tenant"}
	services.Users.SetCurrentTenant(theTenant)
	services.Users.Register(&models.User{
		Name:   "Tony",
		Email:  "tony.stark@avengers.com",
		Tenant: theTenant,
	})
	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"tony.stark@avengers.com",
		},
	}}

	ExpectFailed(action.Validate(nil, services), "recipients")
}

func TestInviteUsers_SampleInvite_IgnoreRecipients(t *testing.T) {
	RegisterT(t)

	action := &actions.InviteUsers{
		IsSampleInvite: true,
		Model: &models.InviteUsers{
			Subject: "Share your feedback.",
			Message: "Use this link to join our community: %invite%",
		},
	}

	ExpectSuccess(action.Validate(nil, services))
}
