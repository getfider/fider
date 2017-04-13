package infra_test

import (
	"testing"

	"github.com/WeCanHearYou/wechy/app/infra"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/models"
	. "github.com/onsi/gomega"
)

func TestStatusHandler(t *testing.T) {
	RegisterTestingT(t)

	settings := &models.WechySettings{
		BuildTime: "today",
	}

	server := mock.NewServer()
	status, query := server.Execute(infra.Status(settings))

	Expect(query.String("build")).To(Equal("today"))
	Expect(status).To(Equal(200))
}
