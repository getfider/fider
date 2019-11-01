package billing

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

var stripeClient *client.API

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "Stripe"
}

func (s Service) Category() string {
	return "billing"
}

func (s Service) Enabled() bool {
	return env.IsBillingEnabled()
}

func (s Service) Init() {
	stripe.DefaultLeveledLogger = &stripe.LeveledLogger{Level: 0}
	stripeClient = &client.API{}
	stripeClient.Init(env.Config.Stripe.SecretKey, nil)

	bus.AddHandler(listPlans)
	bus.AddHandler(getPlanByID)
	bus.AddHandler(cancelSubscription)
	bus.AddHandler(subscribe)
	bus.AddHandler(getUpcomingInvoice)
	bus.AddHandler(createCustomer)
	bus.AddHandler(deleteCustomer)
	bus.AddHandler(getPaymentInfo)
	bus.AddHandler(clearPaymentInfo)
	bus.AddHandler(updatePaymentInfo)
	bus.AddHandler(getAllCountries)
	bus.AddHandler(getCountryByCode)
}

func using(ctx context.Context, handler func(tenant *models.Tenant) error) error {
	tenant, _ := ctx.Value(app.TenantCtxKey).(*models.Tenant)
	return handler(tenant)
}
