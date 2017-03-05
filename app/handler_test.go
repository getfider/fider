package app_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/app"
	"github.com/WeCanHearYou/wchy/mock"
	. "github.com/onsi/gomega"
)

type falsyHealthCheckService struct{}

func (svc falsyHealthCheckService) IsDatabaseOnline() bool {
	return false
}

func TestStatusHandler(t *testing.T) {
	RegisterTestingT(t)

	service := &falsyHealthCheckService{}
	settings := &app.WchySettings{
		BuildTime: "today",
	}

	server := mock.NewServer()
	status, query := server.Execute(app.Status(service, settings))

	Expect(query.String("build")).To(Equal("today"))
	Expect(query.Bool("healthy", "database")).To(Equal(false))
	Expect(status).To(Equal(200))
}
