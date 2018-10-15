package handlers_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestGetMentions(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()

	code, query := server.
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		AddParam("query", "Jon").
		ExecuteAsJSON(handlers.GetMentions())

	Expect(code).Equals(http.StatusOK)
	Expect(query).IsNotNil()

}
