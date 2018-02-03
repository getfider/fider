package middlewares_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/worker"
	. "github.com/onsi/gomega"
)

func TestWebSetup(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.WebSetup(log.NewNoopLogger()))
	status, _ := server.Execute(func(c web.Context) error {
		Expect(c.ActiveTransaction()).NotTo(BeNil())
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestWebSetup_Failure(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.WebSetup(log.NewNoopLogger()))
	status, _ := server.Execute(func(c web.Context) error {
		return c.Failure(errors.New("Something went wrong..."))
	})

	Expect(status).To(Equal(http.StatusInternalServerError))
}

func TestWebSetup_Panic(t *testing.T) {
	RegisterTestingT(t)

	server, _ := mock.NewServer()
	server.Use(middlewares.WebSetup(log.NewNoopLogger()))
	status, _ := server.Execute(func(c web.Context) error {
		panic("Boom!")
	})

	Expect(status).To(Equal(http.StatusInternalServerError))
}

func TestWorkerSetup(t *testing.T) {
	RegisterTestingT(t)

	c := worker.NewContext("0", "Any Task", log.NewNoopLogger())
	mw := middlewares.WorkerSetup(log.NewNoopLogger())
	err := mw(func(c *worker.Context) error {
		Expect(c.Services()).NotTo(BeNil())
		return errors.New("Not Found")
	})(c)
	Expect(err).NotTo(BeNil())
}
