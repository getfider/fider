package actions

import (
	"context"
	"fmt"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/goenning/vat"
)

// CreateEditBillingPaymentInfo is used to create/edit billing payment info
type CreateEditBillingPaymentInfo struct {
	Model *dto.CreateEditBillingPaymentInfo
}

// Initialize the model
func (input *CreateEditBillingPaymentInfo) Initialize() interface{} {
	input.Model = new(dto.CreateEditBillingPaymentInfo)
	return input.Model
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *CreateEditBillingPaymentInfo) IsAuthorized(ctx context.Context, user *models.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (input *CreateEditBillingPaymentInfo) Validate(ctx context.Context, user *models.User) *validate.Result {
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

	if input.Model.AddressPostalCode == "" {
		result.AddFieldFailure("addressPostalCode", "Postal Code is required.")
	}

	getPaymentInfo := &query.GetPaymentInfo{}
	err := bus.Dispatch(ctx, getPaymentInfo)
	if err != nil {
		return validate.Error(err)
	}

	current := getPaymentInfo.Result
	isNew := current == nil
	isUpdate := current != nil && input.Model.Card == nil
	isReplacing := current != nil && input.Model.Card != nil

	if (isNew || isReplacing) && (input.Model.Card == nil || input.Model.Card.Token == "") {
		result.AddFieldFailure("card", "Card information is required.")
	}

	if input.Model.AddressCountry == "" {
		result.AddFieldFailure("addressCountry", "Country is required.")
	} else {
		err := bus.Dispatch(ctx, &query.GetCountryByCode{Code: input.Model.AddressCountry})
		if err != nil {
			if err == app.ErrNotFound {
				result.AddFieldFailure("addressCountry", fmt.Sprintf("'%s' is not a valid country code.", input.Model.AddressCountry))
			} else {
				return validate.Error(err)
			}
		}

		if (isNew || isReplacing) && input.Model.Card != nil && input.Model.AddressCountry != input.Model.Card.Country {
			result.AddFieldFailure("addressCountry", "Country doesn't match with card issue country.")
		} else if isUpdate && input.Model.AddressCountry != current.CardCountry {
			result.AddFieldFailure("addressCountry", "Country doesn't match with card issue country.")
		}

		if isReplacing || isUpdate {
			prevIsEU := vat.IsEU(current.AddressCountry)
			nextIsEU := vat.IsEU(input.Model.AddressCountry)
			if prevIsEU != nextIsEU {
				result.AddFieldFailure("currency", "Billing currency cannot be changed.")
			}
		}
	}

	if input.Model.VATNumber != "" && vat.IsEU(input.Model.AddressCountry) && (isNew || input.Model.VATNumber != current.VATNumber) {
		valid, euCC := vat.ValidateNumberFormat(input.Model.VATNumber)
		if !valid {
			result.AddFieldFailure("vatNumber", "VAT Number is an invalid format.")
		} else if euCC != input.Model.AddressCountry {
			result.AddFieldFailure("vatNumber", "VAT Number doesn't match with selected country.")
		} else {
			resp, err := vat.Query(input.Model.VATNumber)
			if err != nil {
				if err == vat.ErrInvalidVATNumberFormat {
					result.AddFieldFailure("vatNumber", "VAT Number is an invalid format.")
				} else {
					log.Error(ctx, errors.Wrap(err, "failed to validate VAT Number '%s'", input.Model.VATNumber))
					result.AddFieldFailure("vatNumber", "We couldn't validate your VAT Number right now, please try again soon.")
				}
			}
			if !resp.IsValid {
				result.AddFieldFailure("vatNumber", "VAT Number is invalid.")
			}
		}
	}

	return result
}
