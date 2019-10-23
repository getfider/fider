package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/services/billing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/dto"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestCreateEditBillingPaymentInfo_InvalidInput(t *testing.T) {
	RegisterT(t)
	bus.Init(billing.Service{})

	testCases := []struct {
		expected []string
		input    *dto.CreateEditBillingPaymentInfo
	}{
		{
			expected: []string{"card", "name", "email", "addressLine1", "addressLine2", "addressCity", "addressPostalCode", "addressCountry"},
			input:    &dto.CreateEditBillingPaymentInfo{},
		},
		{
			expected: []string{"card", "email", "addressCity", "addressPostalCode", "addressCountry"},
			input: &dto.CreateEditBillingPaymentInfo{
				Name:           "John",
				AddressLine1:   "Street 1",
				AddressLine2:   "Street 2",
				Email:          "jo@a",
				AddressCountry: "PP",
			},
		},
		{
			expected: []string{"card", "email", "addressCity", "addressPostalCode", "addressCountry"},
			input: &dto.CreateEditBillingPaymentInfo{
				Name:           "John",
				AddressLine1:   "Street 1",
				AddressLine2:   "Street 2",
				Email:          "jo@a",
				AddressCountry: "US",
				Card: &dto.CreateEditBillingPaymentInfoCard{
					Country: "IE",
				},
			},
		},
	}

	for _, testCase := range testCases {
		action := &actions.CreateEditBillingPaymentInfo{
			Model: testCase.input,
		}
		ctx := context.WithValue(context.Background(), app.TenantCtxKey, &models.Tenant{ID: 2})
		result := action.Validate(ctx, nil)
		ExpectFailed(result, testCase.expected...)
	}
}

func TestCreateEditBillingPaymentInfo_ValidInput(t *testing.T) {
	RegisterT(t)
	bus.Init(billing.Service{})

	action := &actions.CreateEditBillingPaymentInfo{
		Model: &dto.CreateEditBillingPaymentInfo{
			Name:              "Jon Snow",
			AddressLine1:      "Street 1",
			AddressLine2:      "Street 2",
			AddressCity:       "New York",
			AddressPostalCode: "12345",
			AddressState:      "NY",
			Email:             "jon.show@got.com",
			AddressCountry:    "US",
			Card: &dto.CreateEditBillingPaymentInfoCard{
				Token:   "tok_visa",
				Country: "US",
			},
		},
	}
	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &models.Tenant{ID: 2})
	result := action.Validate(ctx, nil)
	ExpectSuccess(result)
}

func TestCreateEditBillingPaymentInfo_VATNumber(t *testing.T) {
	RegisterT(t)
	bus.Init(billing.Service{})

	ctx := context.WithValue(context.Background(), app.TenantCtxKey, &models.Tenant{ID: 2})

	action := &actions.CreateEditBillingPaymentInfo{
		Model: &dto.CreateEditBillingPaymentInfo{
			Name:              "Jon Snow",
			AddressLine1:      "Street 1",
			AddressLine2:      "Street 2",
			AddressCity:       "New York",
			AddressPostalCode: "12345",
			AddressState:      "NY",
			Email:             "jon.show@got.com",
			AddressCountry:    "IE",
			VATNumber:         "IE0",
			Card: &dto.CreateEditBillingPaymentInfoCard{
				Token:   "tok_visa",
				Country: "IE",
			},
		},
	}
	result := action.Validate(ctx, nil)
	ExpectFailed(result, "vatNumber")

	action.Model.VATNumber = "GB270600730"
	result = action.Validate(ctx, nil)
	ExpectFailed(result, "vatNumber")

	action.Model.VATNumber = "IE6388047A"
	result = action.Validate(ctx, nil)
	ExpectFailed(result, "vatNumber")

	action.Model.VATNumber = "IE6388047V"
	result = action.Validate(ctx, nil)
	ExpectSuccess(result)

	action.Model.VATNumber = ""
	result = action.Validate(ctx, nil)
	ExpectSuccess(result)
}
