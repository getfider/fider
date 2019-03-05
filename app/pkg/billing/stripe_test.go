package billing_test

import (
	"testing"
	"time"

	"github.com/getfider/fider/app/pkg/env"

	"github.com/getfider/fider/app/pkg/billing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestCreateCustomer_WithSubscription(t *testing.T) {
	RegisterT(t)
	if !env.IsBillingEnabled() {
		return
	}

	tenant := &models.Tenant{
		ID:        2,
		Name:      "Game Inc.",
		Subdomain: "gameinc",
		Billing:   &models.TenantBilling{},
	}
	client := billing.NewClient()
	client.SetCurrentTenant(tenant)
	customerID, err := client.CreateCustomer("")
	Expect(err).IsNil()
	Expect(customerID).IsNotEmpty()
	Expect(customerID).Equals(tenant.Billing.StripeCustomerID)

	err = client.UpdatePaymentInfo(&models.CreateEditBillingPaymentInfo{
		Email: "jon.snow@got.com",
		Card: &models.CreateEditBillingPaymentInfoCard{
			Token: "tok_visa",
		},
	})
	Expect(err).IsNil()

	err = client.Subscribe("plan_EKTT1YWe1Zmrtp")
	Expect(err).IsNil()
	Expect(tenant.Billing.StripeSubscriptionID).IsNotEmpty()
	Expect(tenant.Billing.StripePlanID).Equals("plan_EKTT1YWe1Zmrtp")

	inv, err := client.GetUpcomingInvoice()
	Expect(err).IsNil()
	Expect(int(inv.AmountDue)).Equals(900)
	Expect(inv.Currency).Equals("USD")

	err = client.CancelSubscription()
	Expect(err).IsNil()
	Expect(tenant.Billing.SubscriptionEndsAt).IsNotNil()

	err = client.DeleteCustomer()
	Expect(err).IsNil()
}

var forUnitTestingTenant = &models.Tenant{
	ID:        5,
	Name:      "For Unit Testing (DO NOT DELETE)",
	Subdomain: "unittesting",
	Billing: &models.TenantBilling{
		StripeCustomerID: "cus_EICBuXBIkhI2EV",
	},
}

func TestUpdatePaymentInfo(t *testing.T) {
	RegisterT(t)
	if !env.IsBillingEnabled() {
		return
	}

	var firstCardID string

	client := billing.NewClient()
	client.SetCurrentTenant(forUnitTestingTenant)

	err := client.ClearPaymentInfo()
	Expect(err).IsNil()

	//Creating a new card
	err = client.UpdatePaymentInfo(&models.CreateEditBillingPaymentInfo{
		Email:     "jon.snow@got.com",
		VATNumber: "IE1234",
		Card: &models.CreateEditBillingPaymentInfoCard{
			Token: "tok_visa",
		},
		AddressCountry: "IE",
	})
	Expect(err).IsNil()

	info, err := client.GetPaymentInfo()
	Expect(err).IsNil()
	firstCardID = info.StripeCardID

	Expect(info.StripeCardID).IsNotEmpty()
	Expect(info.Email).Equals("jon.snow@got.com")
	Expect(info.VATNumber).Equals("IE1234")
	Expect(info.CardBrand).Equals("Visa")
	Expect(info.CardCountry).Equals("US")
	Expect(info.CardLast4).Equals("4242")
	Expect(int(info.CardExpMonth)).Equals(int(time.Now().Month()))
	Expect(int(info.CardExpYear)).Equals(time.Now().Year() + 1)
	Expect(info.Name).Equals("")
	Expect(info.AddressLine1).Equals("")
	Expect(info.AddressLine2).Equals("")
	Expect(info.AddressCity).Equals("")
	Expect(info.AddressState).Equals("")
	Expect(info.AddressPostalCode).Equals("")
	Expect(info.AddressCountry).Equals("")

	//Update existing card
	err = client.UpdatePaymentInfo(&models.CreateEditBillingPaymentInfo{
		Email:             "jon.snow@got.com",
		Name:              "Jon Snow",
		AddressLine1:      "Street 1",
		AddressLine2:      "Av. ABC",
		AddressCity:       "New York",
		AddressState:      "NYC",
		AddressPostalCode: "12098",
		AddressCountry:    "US",
	})
	Expect(err).IsNil()

	info, err = client.GetPaymentInfo()
	Expect(err).IsNil()
	Expect(info.Name).Equals("Jon Snow")
	Expect(info.VATNumber).Equals("")
	Expect(info.CardLast4).Equals("4242")
	Expect(info.AddressLine1).Equals("Street 1")
	Expect(info.AddressLine2).Equals("Av. ABC")
	Expect(info.AddressCity).Equals("New York")
	Expect(info.AddressState).Equals("NYC")
	Expect(info.AddressPostalCode).Equals("12098")
	Expect(info.AddressCountry).Equals("US")

	//Replace card
	err = client.UpdatePaymentInfo(&models.CreateEditBillingPaymentInfo{
		Email: "jon.snow@got.com",
		Card: &models.CreateEditBillingPaymentInfoCard{
			Token: "tok_br",
		},
	})
	Expect(err).IsNil()

	info, err = client.GetPaymentInfo()
	Expect(err).IsNil()
	Expect(info.StripeCardID).IsNotEmpty()
	Expect(info.StripeCardID).NotEquals(firstCardID)
	Expect(info.Email).Equals("jon.snow@got.com")
	Expect(info.VATNumber).Equals("")
	Expect(info.CardBrand).Equals("Visa")
	Expect(info.CardCountry).Equals("BR")
	Expect(info.CardLast4).Equals("0002")
	Expect(int(info.CardExpMonth)).Equals(int(time.Now().Month()))
	Expect(int(info.CardExpYear)).Equals(time.Now().Year() + 1)
	Expect(info.Name).Equals("")
	Expect(info.AddressLine1).Equals("")
	Expect(info.AddressLine2).Equals("")
	Expect(info.AddressCity).Equals("")
	Expect(info.AddressState).Equals("")
	Expect(info.AddressPostalCode).Equals("")
	Expect(info.AddressCountry).Equals("")
}

func TestListPlans(t *testing.T) {
	RegisterT(t)
	if !env.IsBillingEnabled() {
		return
	}

	client := billing.NewClient()
	plans, err := client.ListPlans("US")
	Expect(err).IsNil()
	Expect(plans).HasLen(3)

	Expect(plans[0].ID).Equals("plan_EKTT1YWe1Zmrtp")
	Expect(plans[0].Name).Equals("Starter")
	Expect(plans[0].Currency).Equals("USD")
	Expect(plans[0].MaxUsers).Equals(200)

	Expect(plans[1].ID).Equals("plan_DoK187GZcnFpKY")
	Expect(plans[1].Name).Equals("Business (monthly)")
	Expect(plans[1].Currency).Equals("USD")
	Expect(plans[1].MaxUsers).Equals(0)

	Expect(plans[2].ID).Equals("plan_DpN9SkJMjNTvLd")
	Expect(plans[2].Name).Equals("Business (yearly)")
	Expect(plans[2].Currency).Equals("USD")
	Expect(plans[2].MaxUsers).Equals(0)

	plans, err = client.ListPlans("DE")
	Expect(err).IsNil()
	Expect(plans).HasLen(3)

	Expect(plans[0].ID).Equals("plan_EKTSnrGmj5BuKI")
	Expect(plans[0].Name).Equals("Starter")
	Expect(plans[0].Currency).Equals("EUR")

	Expect(plans[0].MaxUsers).Equals(200)
	Expect(plans[1].ID).Equals("plan_EKPnahGhiTEnCc")
	Expect(plans[1].Name).Equals("Business (monthly)")
	Expect(plans[1].Currency).Equals("EUR")
	Expect(plans[1].MaxUsers).Equals(0)

	Expect(plans[2].ID).Equals("plan_EKTU4xD7LNI9dO")
	Expect(plans[2].Name).Equals("Business (yearly)")
	Expect(plans[2].Currency).Equals("EUR")
	Expect(plans[2].MaxUsers).Equals(0)
}
