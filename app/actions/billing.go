package actions

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
)

// CreateEditBillingPaymentInfo is used to create/edit billing payment info
type CreateEditBillingPaymentInfo struct {
	Model *models.CreateEditBillingPaymentInfo
}

// Initialize the model
func (input *CreateEditBillingPaymentInfo) Initialize() interface{} {
	input.Model = new(models.CreateEditBillingPaymentInfo)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateEditBillingPaymentInfo) IsAuthorized(user *models.User, services *app.Services) bool {
	return user != nil && user.IsAdministrator()
}

// Validate is current model is valid
func (input *CreateEditBillingPaymentInfo) Validate(user *models.User, services *app.Services) *validate.Result {
	result := validate.Success()

	if input.Model.Name == "" {
		result.AddFieldFailure("name", "Name is required.")
	}

	if input.Model.Email == "" {
		result.AddFieldFailure("email", "Email is required")
	} else {
		messages := validate.Email(input.Model.Email)
		if len(messages) > 0 {
			result.AddFieldFailure("email", messages...)
		}
	}

	if input.Model.AddressLine1 == "" {
		result.AddFieldFailure("addressLine1", "Address Line 1 is required.")
	}

	if input.Model.AddressLine2 == "" {
		result.AddFieldFailure("addressLine2", "Address Line 2 is required.")
	}

	if input.Model.AddressCity == "" {
		result.AddFieldFailure("addressCity", "City is required.")
	}

	if input.Model.AddressState == "" {
		result.AddFieldFailure("addressState", "State/Region is required.")
	}

	if input.Model.AddressPostalCode == "" {
		result.AddFieldFailure("addressPostalCode", "Postal Code is required.")
	}

	if input.Model.AddressCountry == "" {
		result.AddFieldFailure("addressCountry", "Country is required.")
	}

	if input.Model.Card != nil && input.Model.AddressCountry != input.Model.Card.Country {
		result.AddFieldFailure("addressCountry", "Country that doesn't match with card issue country.")
	}

	return result
}
