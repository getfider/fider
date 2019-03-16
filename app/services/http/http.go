package http

import (
	"github.com/getfider/fider/app/pkg/bus"
)

func init() {
	bus.Register(&HTTPClientService{})
}

type HTTPClientService struct{}

func (s HTTPClientService) Enabled() bool {
	return true
}

func (s HTTPClientService) Init() {
	bus.AddHandler(s, HTTPGetRequest)
	bus.AddHandler(s, HTTPPostRequest)
}

type BasicAuth struct {
	User     string
	Password string
}
