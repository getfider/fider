package handler_test

import (
	"testing"

	"github.com/WeCanHearYou/wchy/context"
	"github.com/WeCanHearYou/wchy/handler"
	"github.com/WeCanHearYou/wchy/mock"
	"github.com/WeCanHearYou/wchy/model"
	. "github.com/onsi/gomega"
)

func TestIndexHandler(t *testing.T) {
	RegisterTestingT(t)

	ctx := &context.WchyContext{}

	server := mock.NewServer()
	server.Set("Tenant", &model.Tenant{Name: "Any Tenant"})
	server.Register(handler.Index(ctx))
	status, _ := server.Request()

	Expect(status).To(Equal(200))
}
