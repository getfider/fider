package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
)

// InviteUsers is used to invite new users into Fider
type InviteUsers struct {
	IsSampleInvite bool
	Model          *models.InviteUsers
	Invitations    []*models.UserInvitation
}

// Initialize the model
func (input *InviteUsers) Initialize() interface{} {
	input.Model = new(models.InviteUsers)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *InviteUsers) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (input *InviteUsers) Validate(ctx context.Context, user *models.User) *validate.Result {
	result := validate.Success()

	if input.Model.Subject == "" {
		result.AddFieldFailure("subject", "Subject is required.")
	} else if len(input.Model.Subject) > 70 {
		result.AddFieldFailure("subject", "Subject must have less than 70 characters.")
	}

	if input.Model.Message == "" {
		result.AddFieldFailure("message", "Message is required.")
	} else if !strings.Contains(input.Model.Message, app.InvitePlaceholder) {
		msg := fmt.Sprintf("Your message is missing the invitation link placeholder. Please add '%s' to your message.", app.InvitePlaceholder)
		result.AddFieldFailure("message", msg)
	}

	//When it's a sample invite, we skip recipients validation
	if !input.IsSampleInvite {

		if len(input.Model.Recipients) == 0 {
			result.AddFieldFailure("recipients", "At least one recipient is required.")
		} else if len(input.Model.Recipients) > 30 {
			result.AddFieldFailure("recipients", "Too many recipients. We limit at 30 recipients per invite.")
		}

		for _, email := range input.Model.Recipients {
			if email != "" {
				messages := validate.Email(email)
				result.AddFieldFailure("recipients", messages...)
			}
		}

		if result.Ok {
			input.Invitations = make([]*models.UserInvitation, 0)
			for _, email := range input.Model.Recipients {
				if email != "" {
					err := bus.Dispatch(ctx, &query.GetUserByEmail{Email: email})
					if errors.Cause(err) == app.ErrNotFound {
						input.Invitations = append(input.Invitations, &models.UserInvitation{
							Email:           email,
							VerificationKey: models.GenerateSecretKey(),
						})
					}
				}
			}

			if len(input.Invitations) == 0 {
				result.AddFieldFailure("recipients", "All these addresses have already been registered on this site.")
			}
		}

	}

	return result
}
