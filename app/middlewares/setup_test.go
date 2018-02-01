package middlewares_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/pkg/email"

	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/web"
	. "github.com/onsi/gomega"
)

func TestSetup(t *testing.T) {
	RegisterTestingT(t)

	db := dbx.New()
	defer db.Close()

	server, _ := mock.NewServer()
	server.Use(middlewares.Setup(db, email.NewNoopSender()))
	status, _ := server.Execute(func(c web.Context) error {
		Expect(c.ActiveTransaction()).NotTo(BeNil())
		return c.NoContent(http.StatusOK)
	})

	Expect(status).To(Equal(http.StatusOK))
}

func TestSetup_Failure(t *testing.T) {
	RegisterTestingT(t)

	db := dbx.New()
	defer db.Close()

	server, _ := mock.NewServer()
	server.Use(middlewares.Setup(db, email.NewNoopSender()))
	status, _ := server.Execute(func(c web.Context) error {
		return c.Failure(errors.New("Something went wrong..."))
	})

	Expect(status).To(Equal(http.StatusInternalServerError))
}

func TestSetup_Panic(t *testing.T) {
	RegisterTestingT(t)

	db := dbx.New()
	defer db.Close()

	server, _ := mock.NewServer()
	server.Use(middlewares.Setup(db, email.NewNoopSender()))
	status, _ := server.Execute(func(c web.Context) error {
		panic("Boom!")
	})

	Expect(status).To(Equal(http.StatusInternalServerError))
}
