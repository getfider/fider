package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
)

// InviteUsers is used to invite new users into Fider
type InviteUsers struct {
	Subject        string   `json:"subject"`
	Message        string   `json:"message"`
	Recipients     []string `json:"recipients" format:"lower"`
	IsSampleInvite bool              `json:"-"`
	Invitations    []*UserInvitation `json:"-"`
}

// DefaultInviteSubject returns the default invite subject for the given tenant.
// Mirrors the placeholder shown in the admin UI so non-pro tenants get a sensible default.
func DefaultInviteSubject(tenant *entity.Tenant) string {
	return fmt.Sprintf("[%s] We would like to hear from you!", tenant.Name)
}

// DefaultInviteMessage returns the default invite message for the given tenant and inviter.
// Mirrors the placeholder shown in the admin UI so non-pro tenants get a sensible default.
func DefaultInviteMessage(tenant *entity.Tenant, inviter *entity.User) string {
	return fmt.Sprintf("Hi,\n\nWe are inviting you to join the %s feedback site, a place where you can vote, discuss and share your ideas and thoughts on how to improve our services!\n\nClick the link below to join!\n\n%s\n\nRegards,\n%s (%s)",
		tenant.Name, app.InvitePlaceholder, inviter.Name, tenant.Name)
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *InviteUsers) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (action *InviteUsers) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	// On hosted Fider, only pro tenants may customize the invite copy. Free tenants
	// are forced onto safe defaults to prevent abuse via phishing-style copy.
	// Self-hosted (billing disabled) is unaffected.
	if env.IsBillingEnabled() {
		if tenant, ok := ctx.Value(app.TenantCtxKey).(*entity.Tenant); ok && tenant != nil && !tenant.IsPro {
			action.Subject = DefaultInviteSubject(tenant)
			action.Message = DefaultInviteMessage(tenant, user)
		}
	}

	if action.Subject == "" {
		result.AddFieldFailure("subject", "Subject is required.")
	} else if len(action.Subject) > 70 {
		result.AddFieldFailure("subject", "Subject must have less than 70 characters.")
	}

	if action.Message == "" {
		result.AddFieldFailure("message", "Message is required.")
	} else if !strings.Contains(action.Message, app.InvitePlaceholder) {
		msg := fmt.Sprintf("Your message is missing the invitation link placeholder. Please add '%s' to your message.", app.InvitePlaceholder)
		result.AddFieldFailure("message", msg)
	}

	//When it's a sample invite, we skip recipients validation
	if !action.IsSampleInvite {

		if len(action.Recipients) == 0 {
			result.AddFieldFailure("recipients", "At least one recipient is required.")
		} else if len(action.Recipients) > 30 {
			result.AddFieldFailure("recipients", "Too many recipients. We limit at 30 recipients per invite.")
		}

		for _, email := range action.Recipients {
			if email != "" {
				messages := validate.Email(ctx, email)
				result.AddFieldFailure("recipients", messages...)
			}
		}

		if result.Ok {
			action.Invitations = make([]*UserInvitation, 0)
			for _, email := range action.Recipients {
				if email != "" {
					err := bus.Dispatch(ctx, &query.GetUserByEmail{Email: email})
					if errors.Cause(err) == app.ErrNotFound {
						action.Invitations = append(action.Invitations, &UserInvitation{
							Email:           email,
							VerificationKey: entity.GenerateEmailVerificationKey(),
						})
					}
				}
			}

			if len(action.Invitations) == 0 {
				result.AddFieldFailure("recipients", "All these addresses have already been registered on this site.")
			}
		}

	}

	return result
}

//UserInvitation is the model used to register an invite sent to an user
type UserInvitation struct {
	Email           string
	VerificationKey string
}

//GetEmail returns the invited user's email
func (e *UserInvitation) GetEmail() string {
	return e.Email
}

//GetName returns empty for this kind of process
func (e *UserInvitation) GetName() string {
	return ""
}

//GetUser returns the current user performing this action
func (e *UserInvitation) GetUser() *entity.User {
	return nil
}

//GetKind returns EmailVerificationKindUserInvitation
func (e *UserInvitation) GetKind() enum.EmailVerificationKind {
	return enum.EmailVerificationKindUserInvitation
}
