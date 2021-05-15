package actions_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/rand"
)

func TestCreateUser_InvalidInput(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		expected []string
		message  string
		action   *actions.CreateUser
	}{
		{
			expected: []string{"name"},
			message:  "Either email or reference is required",
			action:   &actions.CreateUser{},
		},
		{
			expected: []string{"email"},
			action: &actions.CreateUser{
				Name:  "Jon Snow",
				Email: "helloworld",
			},
		},
		{
			expected: []string{"name", "email", "reference"},
			action: &actions.CreateUser{
				Name:      rand.String(101),
				Email:     rand.String(201),
				Reference: rand.String(101),
			},
		},
	}

	for _, testCase := range testCases {
		result := testCase.action.Validate(context.Background(), nil)
		ExpectFailed(result, testCase.expected...)
		if testCase.message != "" {
			for k, v := range result.Errors {
				if v.Field == "" {
					Expect(result.Errors[k].Message).Equals(testCase.message)
				}
			}
		}
	}
}

func TestCreateUser_ValidInput(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		action *actions.CreateUser
	}{
		{
			action: &actions.CreateUser{
				Name:      "John Snow",
				Email:     "jon.snow@got.com",
				Reference: "812747824",
			},
		},
		{
			action: &actions.CreateUser{
				Name:  "John Snow",
				Email: "jon.snow@got.com",
			},
		},
		{
			action: &actions.CreateUser{
				Name:      "John Snow",
				Reference: "812747824",
			},
		},
	}

	for _, testCase := range testCases {
		result := testCase.action.Validate(context.Background(), nil)
		ExpectSuccess(result)
	}
}

func TestChangeUserRole_Unauthorized(t *testing.T) {
	RegisterT(t)

	for _, user := range []*entity.User{
		{ID: 1, Role: enum.RoleVisitor},
		{ID: 1, Role: enum.RoleCollaborator},
		{ID: 2, Role: enum.RoleAdministrator},
	} {
		action := actions.ChangeUserRole{UserID: 2}
		Expect(action.IsAuthorized(context.Background(), user)).IsFalse()
	}
}

func TestChangeUserRole_Authorized(t *testing.T) {
	RegisterT(t)

	user := &entity.User{ID: 2, Role: enum.RoleAdministrator}
	action := actions.ChangeUserRole{UserID: 1}
	Expect(action.IsAuthorized(context.Background(), user)).IsTrue()
}

func TestChangeUserRole_InvalidRole(t *testing.T) {
	RegisterT(t)

	targetUser := &entity.User{Role: enum.RoleVisitor}
	currentUser := &entity.User{Role: enum.RoleAdministrator}

	action := actions.ChangeUserRole{UserID: targetUser.ID, Role: 4}
	action.IsAuthorized(context.Background(), currentUser)
	result := action.Validate(context.Background(), currentUser)
	Expect(result.Err).Equals(app.ErrNotFound)
}

func TestChangeUserRole_InvalidUser(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		return app.ErrNotFound
	})

	currentUser := &entity.User{
		Tenant: &entity.Tenant{ID: 1},
		Role:   enum.RoleAdministrator,
	}

	ctx := context.Background()
	action := actions.ChangeUserRole{UserID: 999, Role: enum.RoleAdministrator}
	action.IsAuthorized(ctx, currentUser)
	result := action.Validate(ctx, currentUser)
	ExpectFailed(result, "userID")
}

func TestChangeUserRole_InvalidUser_Tenant(t *testing.T) {
	RegisterT(t)

	targetUser := &entity.User{
		Tenant: &entity.Tenant{ID: 1},
	}

	currentUser := &entity.User{
		Tenant: &entity.Tenant{ID: 2},
		Role:   enum.RoleAdministrator,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == targetUser.ID {
			q.Result = targetUser
			return nil
		}
		return app.ErrNotFound
	})

	action := actions.ChangeUserRole{UserID: targetUser.ID, Role: enum.RoleAdministrator}
	action.IsAuthorized(context.Background(), currentUser)
	result := action.Validate(context.Background(), currentUser)
	ExpectFailed(result, "userID")
}

func TestChangeUserRole_CurrentUser(t *testing.T) {
	RegisterT(t)

	currentUser := &entity.User{
		Tenant: &entity.Tenant{ID: 2},
		Role:   enum.RoleAdministrator,
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetUserByID) error {
		if q.UserID == currentUser.ID {
			q.Result = currentUser
			return nil
		}
		return app.ErrNotFound
	})

	action := actions.ChangeUserRole{UserID: currentUser.ID, Role: enum.RoleVisitor}
	action.IsAuthorized(context.Background(), currentUser)
	result := action.Validate(context.Background(), currentUser)
	ExpectFailed(result, "userID")
}
