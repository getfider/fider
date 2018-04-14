package actions_test

import (
	"fmt"
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"

	. "github.com/onsi/gomega"
)

func TestInviteUsers_Empty(t *testing.T) {
	RegisterTestingT(t)

	action := &actions.InviteUsers{Model: &models.InviteUsers{}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "subject", "message", "recipients")
	Expect(action.Invitations).To(BeNil())
}

func TestInviteUsers_Oversized(t *testing.T) {
	RegisterTestingT(t)

	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject:    "Join us and share your ideas. Because we have a cool website and this subject needs to be very long",
		Message:    "Use this link to join %invite%",
		Recipients: []string{"jon.snow@got.com"},
	}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "subject")
	Expect(action.Invitations).To(BeNil())
}

func TestInviteUsers_MissingInvite(t *testing.T) {
	RegisterTestingT(t)

	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject:    "Share your feedback.",
		Message:    "Please!",
		Recipients: []string{"jon.snow@got.com"},
	}}
	result := action.Validate(nil, services)
	ExpectFailed(result, "message")
	Expect(action.Invitations).To(BeNil())
}

func TestInviteUsers_TooManyRecipients(t *testing.T) {
	RegisterTestingT(t)

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
	Expect(action.Invitations).To(BeNil())
}

func TestInviteUsers_InvalidRecipient(t *testing.T) {
	RegisterTestingT(t)

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
	Expect(action.Invitations).To(BeNil())
}

func TestInviteUsers_Valid(t *testing.T) {
	RegisterTestingT(t)

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

	Expect(action.Invitations).To(HaveLen(2))

	Expect(action.Invitations[0].Email).To(Equal("jon.snow@got.com"))
	Expect(action.Invitations[0].VerificationKey).NotTo(BeEmpty())

	Expect(action.Invitations[1].Email).To(Equal("arya.stark@got.com"))
	Expect(action.Invitations[1].VerificationKey).NotTo(BeEmpty())
}

func TestInviteUsers_IgnoreAlreadyRegistered(t *testing.T) {
	RegisterTestingT(t)

	theTenant := &models.Tenant{ID: 1, Name: "The Tenant"}
	services.Users.SetCurrentTenant(theTenant)
	services.Users.Register(&models.User{
		Name:   "Tony",
		Email:  "tony.start@avengers.com",
		Tenant: theTenant,
	})
	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"tony.start@avengers.com",
			"jon.snow@got.com",
			"arya.stark@got.com",
		},
	}}

	ExpectSuccess(action.Validate(nil, services))

	Expect(action.Invitations).To(HaveLen(2))

	Expect(action.Invitations[0].Email).To(Equal("jon.snow@got.com"))
	Expect(action.Invitations[0].VerificationKey).NotTo(BeEmpty())

	Expect(action.Invitations[1].Email).To(Equal("arya.stark@got.com"))
	Expect(action.Invitations[1].VerificationKey).NotTo(BeEmpty())
}
