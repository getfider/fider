package handler_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/WeCanHearYou/wchy/mock"
	. "github.com/onsi/gomega"
)

type falsyHealthCheckService struct{}

func (svc falsyHealthCheckService) IsDatabaseOnline() bool {
	return false
}

func TestStatusHandler(t *testing.T) {
	RegisterTestingT(t)

	ctx := &context.WchyContext{
		Health: &falsyHealthCheckService{},
		Settings: context.WchySettings{
			BuildTime: "today",
		},
	}

	server := mock.NewServer()
	server.Register(handler.Status(ctx))
	status, query := server.Request()

	Expect(query.String("build")).To(Equal("today"))
	Expect(query.Bool("healthy", "database")).To(Equal(false))
	Expect(status).To(Equal(200))
}
