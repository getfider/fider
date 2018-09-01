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

func TestCreateUser_ExistingEmail(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)

	status, query := server.
		AsUser(mock.JonSnow).
		OnTenant(mock.DemoTenant).
		ExecutePostAsJSON(apiv1.CreateUser(),
			`{
				"name": "Arya",
				"email": "arya.stark@got.com",
				"reference": "AA564645"
			}`)

	Expect(status).Equals(http.StatusOK)
	id := query.Int32("id")
	user, err := services.Users.GetByID(id)
	Expect(err).IsNil()
	Expect(user.ID).Equals(mock.AryaStark.ID)
	Expect(user.Name).Equals("Arya Stark")
	Expect(user.Email).Equals("arya.stark@got.com")
	Expect(user.Tenant).Equals(mock.DemoTenant)
	Expect(user.Providers).HasLen(1)
	Expect(user.Providers[0].Name).Equals("reference")
	Expect(user.Providers[0].UID).Equals("AA564645")
}

func TestCreateUser_ExistingEmail_NoReference(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)

	status, query := server.
		AsUser(mock.JonSnow).
		OnTenant(mock.DemoTenant).
		ExecutePostAsJSON(apiv1.CreateUser(),
			`{
				"name": "Arya",
				"email": "arya.stark@got.com"
			}`)

	Expect(status).Equals(http.StatusOK)
	id := query.Int32("id")
	user, err := services.Users.GetByID(id)
	Expect(err).IsNil()
	Expect(user.ID).Equals(mock.AryaStark.ID)
	Expect(user.Name).Equals("Arya Stark")
	Expect(user.Email).Equals("arya.stark@got.com")
	Expect(user.Tenant).Equals(mock.DemoTenant)
	Expect(user.Providers).HasLen(0)
}

func TestCreateUser_NewUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.SetCurrentTenant(mock.DemoTenant)
	services.SetCurrentUser(mock.JonSnow)

	status, query := server.
		AsUser(mock.JonSnow).
		OnTenant(mock.DemoTenant).
		ExecutePostAsJSON(apiv1.CreateUser(),
			`{
				"name": "Martin",
				"email": "martin@company.com",
				"reference": "89014714"
			}`)

	Expect(status).Equals(http.StatusOK)
	userID := query.Int32("id")
	user, err := services.Users.GetByID(userID)
	Expect(err).IsNil()
	Expect(user.ID).Equals(userID)
	Expect(user.Name).Equals("Martin")
	Expect(user.Email).Equals("martin@company.com")
	Expect(user.Tenant).Equals(mock.DemoTenant)
	Expect(user.Providers).HasLen(1)
	Expect(user.Providers[0].Name).Equals("reference")
	Expect(user.Providers[0].UID).Equals("89014714")

	// Try to recreate another user with same reference
	status, query = server.
		AsUser(mock.JonSnow).
		OnTenant(mock.DemoTenant).
		ExecutePostAsJSON(apiv1.CreateUser(),
			`{
				"name": "The Other Martin",
				"reference": "89014714"
			}`)

	Expect(status).Equals(http.StatusOK)
	theOtherUserID := query.Int32("id")
	Expect(theOtherUserID).Equals(userID)
	user, err = services.Users.GetByID(theOtherUserID)
	Expect(err).IsNil()
	Expect(user.Name).Equals("Martin")
	Expect(user.Email).Equals("martin@company.com")
	Expect(user.Providers).HasLen(1)
	Expect(user.Providers[0].Name).Equals("reference")
	Expect(user.Providers[0].UID).Equals("89014714")
}
