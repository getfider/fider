package apiv1_test

import (
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers/apiv1"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestListUsersHandler(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)

	status, query := server.
		AsUser(mock.JonSnow).
		ExecuteAsJSON(apiv1.ListUsers())

	Expect(status).Equals(http.StatusOK)
	Expect(query.IsArray()).IsTrue()
	Expect(query.ArrayLength()).Equals(2)
}
