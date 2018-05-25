package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestHealthHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		Execute(handlers.Health())

	Expect(code).Equals(http.StatusOK)
}

func TestPageHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		Execute(handlers.Page("The Title", "The Description"))

	Expect(code).Equals(http.StatusOK)
}

func TestLegalPageHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		Execute(handlers.LegalPage("Terms of Service", "terms.md"))

	Expect(code).Equals(http.StatusOK)
}

func TestLegalPageHandler_Invalid(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		Execute(handlers.LegalPage("Some Page", "somepage.md"))

	Expect(code).Equals(http.StatusNotFound)
}
