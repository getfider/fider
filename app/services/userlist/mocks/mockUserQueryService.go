package userlist_mock

import (
	"context"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
)

type Service struct{}

func (s Service) Name() string {
	return "PostgreSQL"
}

func (s Service) Category() string {
	return "sqlstore"
}

func (s Service) Enabled() bool {
	return true
}

func (s Service) Init() {
	bus.AddHandler(GetUserByID)
}

func GetUserByID(ctx context.Context, q *query.GetUserByID) error {
	q.Result = &entity.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Tenant:    &entity.Tenant{ID: 1, Name: "Example Tenant"},
		Role:      enum.RoleAdministrator,
		Providers: []*entity.UserProvider{},
		Status:    enum.UserActive,
	}
	return nil
}
