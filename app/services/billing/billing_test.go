package billing_test

import (
	"context"
	"testing"
	"time"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/services/billing"
)

func TestCreateCustomer_WithSubscription(t *testing.T) {
	RegisterT(t)
	if !env.IsBillingEnabled() {
		return
	}

	bus.Init(billing.Service{})

	tenant := &models.Tenant{
		ID:        2,
		Name:      "Game Inc.",
		Subdomain: "gameinc",
		Billing:   &models.TenantBilling{},
	}
	ctx := context.WithValue(context.Background(), app.TenantCtxKey, tenant)

	err := bus.Dispatch(ctx, &cmd.CreateBillingCustomer{})
	Expect(err).IsNil()
	Expect(tenant.Billing.StripeCustomerID).IsNotEmpty()

	err = bus.Dispatch(ctx, &cmd.UpdatePaymentInfo{
		Input: &dto.CreateEditBillingPaymentInfo{
			Email: "jon.snow@got.com",
			Card: &dto.CreateEditBillingPaymentInfoCard{
				Token: "tok_visa",
			},
		},
	})
	Expect(err).IsNil()

	createSubscription := &cmd.CreateBillingSubscription{PlanID: "plan_EKTT1YWe1Zmrtp"}
	err = bus.Dispatch(ctx, createSubscription)
	Expect(err).IsNil()
	Expect(tenant.Billing.StripeSubscriptionID).IsNotEmpty()
	Expect(tenant.Billing.StripePlanID).Equals("plan_EKTT1YWe1Zmrtp")

	invoiceQuery := &query.GetUpcomingInvoice{}
	err = bus.Dispatch(ctx, invoiceQuery)
	Expect(err).IsNil()
	Expect(int(invoiceQuery.Result.AmountDue)).Equals(900)
	Expect(invoiceQuery.Result.Currency).Equals("USD")

	err = bus.Dispatch(ctx, &cmd.CancelBillingSubscription{})
	Expect(err).IsNil()
	Expect(tenant.Billing.SubscriptionEndsAt).IsNotNil()

	err = bus.Dispatch(ctx, &cmd.DeleteBillingCustomer{})
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

	bus.Init(billing.Service{})

	var firstCardID string

	ctx := context.WithValue(context.Background(), app.TenantCtxKey, forUnitTestingTenant)

	err := bus.Dispatch(ctx, &cmd.ClearPaymentInfo{})
	Expect(err).IsNil()

	//Creating a new card
	err = bus.Dispatch(ctx, &cmd.UpdatePaymentInfo{
		Input: &dto.CreateEditBillingPaymentInfo{
			Email:     "jon.snow@got.com",
			VATNumber: "IE1234",
			Card: &dto.CreateEditBillingPaymentInfoCard{
				Token: "tok_visa",
			},
			AddressCountry: "IE",
		},
	})
	Expect(err).IsNil()

	paymentInfoQuery := &query.GetPaymentInfo{}
	err = bus.Dispatch(ctx, paymentInfoQuery)
	Expect(err).IsNil()
	info := paymentInfoQuery.Result

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
	err = bus.Dispatch(ctx, &cmd.UpdatePaymentInfo{
		Input: &dto.CreateEditBillingPaymentInfo{
			Email:             "jon.snow@got.com",
			Name:              "Jon Snow",
			AddressLine1:      "Street 1",
			AddressLine2:      "Av. ABC",
			AddressCity:       "New York",
			AddressState:      "NYC",
			AddressPostalCode: "12098",
			AddressCountry:    "US",
		},
	})
	Expect(err).IsNil()

	paymentInfoQuery = &query.GetPaymentInfo{}
	err = bus.Dispatch(ctx, paymentInfoQuery)
	Expect(err).IsNil()
	info = paymentInfoQuery.Result

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
	err = bus.Dispatch(ctx, &cmd.UpdatePaymentInfo{
		Input: &dto.CreateEditBillingPaymentInfo{
			Email: "jon.snow@got.com",
			Card: &dto.CreateEditBillingPaymentInfoCard{
				Token: "tok_br",
			},
		},
	})
	Expect(err).IsNil()

	paymentInfoQuery = &query.GetPaymentInfo{}
	err = bus.Dispatch(ctx, paymentInfoQuery)
	Expect(err).IsNil()
	info = paymentInfoQuery.Result

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

	bus.Init(billing.Service{})
	ctx := context.Background()

	q := &query.ListBillingPlans{CountryCode: "US"}
	err := bus.Dispatch(ctx, q)
	Expect(err).IsNil()
	Expect(q.Result).HasLen(3)

	Expect(q.Result[0].ID).Equals("plan_EKTT1YWe1Zmrtp")
	Expect(q.Result[0].Name).Equals("Starter")
	Expect(q.Result[0].Currency).Equals("USD")
	Expect(q.Result[0].MaxUsers).Equals(200)

	Expect(q.Result[1].ID).Equals("plan_DoK187GZcnFpKY")
	Expect(q.Result[1].Name).Equals("Business (monthly)")
	Expect(q.Result[1].Currency).Equals("USD")
	Expect(q.Result[1].MaxUsers).Equals(0)

	Expect(q.Result[2].ID).Equals("plan_DpN9SkJMjNTvLd")
	Expect(q.Result[2].Name).Equals("Business (yearly)")
	Expect(q.Result[2].Currency).Equals("USD")
	Expect(q.Result[2].MaxUsers).Equals(0)

	q = &query.ListBillingPlans{CountryCode: "DE"}
	err = bus.Dispatch(ctx, q)
	Expect(err).IsNil()
	Expect(q.Result).HasLen(3)

	Expect(q.Result[0].ID).Equals("plan_EKTSnrGmj5BuKI")
	Expect(q.Result[0].Name).Equals("Starter")
	Expect(q.Result[0].Currency).Equals("EUR")

	Expect(q.Result[0].MaxUsers).Equals(200)
	Expect(q.Result[1].ID).Equals("plan_EKPnahGhiTEnCc")
	Expect(q.Result[1].Name).Equals("Business (monthly)")
	Expect(q.Result[1].Currency).Equals("EUR")
	Expect(q.Result[1].MaxUsers).Equals(0)

	Expect(q.Result[2].ID).Equals("plan_EKTU4xD7LNI9dO")
	Expect(q.Result[2].Name).Equals("Business (yearly)")
	Expect(q.Result[2].Currency).Equals("EUR")
	Expect(q.Result[2].MaxUsers).Equals(0)
}
