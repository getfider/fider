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
}

func TestInviteUsers_Valid(t *testing.T) {
	RegisterTestingT(t)

	recipients := make([]string, 31)
	for i := 0; i < len(recipients); i++ {
		recipients[i] = fmt.Sprintf("jon.snow%d@got.com", i)
	}

	action := &actions.InviteUsers{Model: &models.InviteUsers{
		Subject: "Share your feedback.",
		Message: "Use this link to join our community: %invite%",
		Recipients: []string{
			"jon.snow@got.com",
			"arya.stark@got.com",
		},
	}}
	ExpectSuccess(action.Validate(nil, services))
}
