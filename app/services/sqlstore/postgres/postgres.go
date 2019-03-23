package postgres

import (
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/stripe/stripe-go/client"
)

var stripeClient *client.API

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "PostgreSQL"
}

func (s Service) Category() string {
	return "sqlstore"
}

func (s Service) Enabled() bool {
	return true
}

func (s Service) Init() {
	bus.AddHandler(storeEvent)
}
