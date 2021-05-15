package apiv1_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestListUsersHandler(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetAllUsers) error {
		q.Result = []*entity.User{
			{ID: 1, Name: "User 1"},
			{ID: 2, Name: "User 2"},
		}
		return nil
	})

	server := mock.NewServer()

	status, query := server.
		AsUser(mock.JonSnow).
		ExecuteAsJSON(apiv1.ListUsers())

	Expect(status).Equals(http.StatusOK)
	Expect(query.IsArray()).IsTrue()
	Expect(query.ArrayLength()).Equals(2)
}

func TestCreateUser_ExistingEmail(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == mock.AryaStark.Email {
			q.Result = mock.AryaStark
			return nil
		}
		return app.ErrNotFound
	})

	var newProvider *cmd.RegisterUserProvider
	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUserProvider) error {
		newProvider = c
		return nil
	})

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
	Expect(id).Equals(mock.AryaStark.ID)
	Expect(newProvider.UserID).Equals(id)
	Expect(newProvider.ProviderName).Equals("reference")
	Expect(newProvider.ProviderUID).Equals("AA564645")
}

func TestCreateUser_ExistingEmail_NoReference(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		if q.Email == mock.AryaStark.Email {
			q.Result = mock.AryaStark
			return nil
		}
		return app.ErrNotFound
	})

	status, query := server.
		AsUser(mock.JonSnow).
		OnTenant(mock.DemoTenant).
		ExecutePostAsJSON(apiv1.CreateUser(),
			`{
				"name": "Arya",
				"email": "arya.stark@got.com"
			}`)

	Expect(status).Equals(http.StatusOK)
	Expect(query.Int32("id")).Equals(mock.AryaStark.ID)
}

func TestCreateUser_NewUser(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()

	var newUser *entity.User

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByProvider) error {
		if newUser != nil {
			if q.Provider == newUser.Providers[0].Name && q.UID == newUser.Providers[0].UID {
				q.Result = newUser
				return nil
			}
		}
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByEmail) error {
		return app.ErrNotFound
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUser) error {
		if newUser != nil {
			panic("newUser is already set")
		}

		c.User.ID = 1
		newUser = c.User
		return nil
	})

	bus.AddHandler(func(ctx context.Context, c *cmd.RegisterUserProvider) error {
		if c.UserID == newUser.ID {
			newUser.Providers = append(newUser.Providers, &entity.UserProvider{
				Name: c.ProviderName,
				UID:  c.ProviderUID,
			})
		}
		return nil
	})

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
	Expect(newUser.ID).Equals(userID)
	Expect(newUser.Name).Equals("Martin")
	Expect(newUser.Email).Equals("martin@company.com")
	Expect(newUser.Tenant).Equals(mock.DemoTenant)
	Expect(newUser.Role).Equals(enum.RoleVisitor)
	Expect(newUser.Providers).HasLen(1)
	Expect(newUser.Providers[0].UID).Equals("89014714")
	Expect(newUser.Providers[0].Name).Equals("reference")

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
}
