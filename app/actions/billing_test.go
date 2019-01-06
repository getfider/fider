package actions_test

import (
	"testing"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestCreateEditBillingPaymentInfo_InvalidInput(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		expected []string
		input    *models.CreateEditBillingPaymentInfo
	}{
		{
			expected: []string{"card", "name", "email", "addressLine1", "addressLine2", "addressCity", "addressState", "addressPostalCode", "addressCountry"},
			input:    &models.CreateEditBillingPaymentInfo{},
		},
		{
			expected: []string{"card", "email", "addressCity", "addressState", "addressPostalCode", "addressCountry"},
			input: &models.CreateEditBillingPaymentInfo{
				Name:           "John",
				AddressLine1:   "Street 1",
				AddressLine2:   "Street 2",
				Email:          "jo@a",
				AddressCountry: "PP",
			},
		},
		{
			expected: []string{"card", "email", "addressCity", "addressState", "addressPostalCode", "addressCountry"},
			input: &models.CreateEditBillingPaymentInfo{
				Name:           "John",
				AddressLine1:   "Street 1",
				AddressLine2:   "Street 2",
				Email:          "jo@a",
				AddressCountry: "US",
				Card: &models.CreateEditBillingPaymentInfoCard{
					Country: "IE",
				},
			},
		},
	}

	for _, testCase := range testCases {
		action := &actions.CreateEditBillingPaymentInfo{
			Model: testCase.input,
		}
		services.Billing.SetCurrentTenant(&models.Tenant{ID: 2})
		result := action.Validate(nil, services)
		ExpectFailed(result, testCase.expected...)
	}
}

func TestCreateEditBillingPaymentInfo_ValidInput(t *testing.T) {
	RegisterT(t)

	action := &actions.CreateEditBillingPaymentInfo{
		Model: &models.CreateEditBillingPaymentInfo{
			Name:              "Jon Snow",
			AddressLine1:      "Street 1",
			AddressLine2:      "Street 2",
			AddressCity:       "New York",
			AddressPostalCode: "12345",
			AddressState:      "NY",
			Email:             "jon.show@got.com",
			AddressCountry:    "US",
			Card: &models.CreateEditBillingPaymentInfoCard{
				Token:   "tok_visa",
				Country: "US",
			},
		},
	}
	services.Billing.SetCurrentTenant(&models.Tenant{ID: 2})
	result := action.Validate(nil, services)
	ExpectSuccess(result)
}
