package entity_test

import (
	"testing"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestTenant_CanCreateAccountByEmail(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		tenant     entity.Tenant
		signIn     bool
		invitation bool
	}{
		{entity.Tenant{IsEmailAuthAllowed: true}, true, true},
		{entity.Tenant{IsEmailAuthAllowed: true, IsSignUpDisabled: true}, false, true},
		{entity.Tenant{IsEmailAuthAllowed: true, IsPrivate: true}, false, true},
		{entity.Tenant{IsEmailAuthAllowed: false}, false, true},
	}

	for _, tc := range testCases {
		Expect(tc.tenant.CanCreateAccountByEmail(enum.EmailVerificationKindSignIn)).Equals(tc.signIn)
		Expect(tc.tenant.CanCreateAccountByEmail(enum.EmailVerificationKindUserInvitation)).Equals(tc.invitation)
	}
}

func TestTenant_CanCreateAccountByOAuth(t *testing.T) {
	RegisterT(t)

	trusted := &entity.OAuthConfig{IsTrusted: true}
	untrusted := &entity.OAuthConfig{IsTrusted: false}

	testCases := []struct {
		tenant    entity.Tenant
		builtIn   bool
		untrusted bool
		trusted   bool
	}{
		{entity.Tenant{}, true, true, true},
		{entity.Tenant{IsSignUpDisabled: true}, false, false, true},
		{entity.Tenant{IsPrivate: true}, false, false, true},
		// Email login doesn't affect OAuth sign-ups
		{entity.Tenant{IsEmailAuthAllowed: false}, true, true, true},
	}

	for _, tc := range testCases {
		Expect(tc.tenant.CanCreateAccountByOAuth(nil)).Equals(tc.builtIn)
		Expect(tc.tenant.CanCreateAccountByOAuth(untrusted)).Equals(tc.untrusted)
		Expect(tc.tenant.CanCreateAccountByOAuth(trusted)).Equals(tc.trusted)
	}
}
