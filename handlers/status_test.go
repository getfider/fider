package handlers_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handlers"
	. "github.com/onsi/gomega"
)

type inMemoryHealthCheckService struct {
	status bool
}

func (svc inMemoryHealthCheckService) IsDatabaseOnline() bool {
	return svc.status
}

func TestStatusHandler(t *testing.T) {
	RegisterTestingT(t)
	ctx := context.WchyContext{
		Health: &inMemoryHealthCheckService{status: false},
		Settings: context.WchySettings{
			BuildTime: "today",
		},
	}

	server := NewTestServer()
	server.register(handlers.Status(ctx))
	status, query := server.request()

	Expect(query.String("build")).To(Equal("today"))
	Expect(query.Bool("healthy", "database")).To(Equal(false))
	Expect(status).To(Equal(200))
}
