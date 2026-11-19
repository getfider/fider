package apiv1

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// maxListUsersLimit is the largest page size accepted by ListUsers
const maxListUsersLimit = 100

// ListUsers returns paginated registered users
func ListUsers() web.HandlerFunc {
	return func(c *web.Context) error {
		page, _ := c.QueryParamAsInt("page")
		if page <= 0 {
			page = 1
		}

		limit, _ := c.QueryParamAsInt("limit")
		if limit <= 0 {
			limit = 10
		} else if limit > maxListUsersLimit {
			limit = maxListUsersLimit
		}

		searchUsers := &query.SearchUsers{
			Query: c.QueryParam("query"),
			Roles: c.QueryParamAsArray("roles"),
			Page:  page,
			Limit: limit,
		}

		if err := bus.Dispatch(c, searchUsers); err != nil {
			return c.Failure(err)
		}

		// Create an array of UserWithEmail structs to include email in JSON response
		allUsersWithEmail := make([]entity.UserWithEmail, len(searchUsers.Result))
		for i, user := range searchUsers.Result {
			allUsersWithEmail[i] = entity.UserWithEmail{
				User: user,
			}
		}

		totalPages := (searchUsers.TotalCount + limit - 1) / limit

		return c.Ok(web.Map{
			"users":      allUsersWithEmail,
			"totalCount": searchUsers.TotalCount,
			"totalPages": totalPages,
			"page":       page,
			"limit":      limit,
		})
	}
}

func ListTaggableUsers() web.HandlerFunc {
	return func(c *web.Context) error {
		allUsers := &query.GetAllUsersNames{}
		if err := bus.Dispatch(c, allUsers); err != nil {
			return c.Failure(err)
		}
		return c.Ok(allUsers.Result)
	}
}

// CreateUser is used to create new users
func CreateUser() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.CreateUser)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		var user *entity.User

		getByReference := &query.GetUserByProvider{Provider: "reference", UID: action.Reference}
		err := bus.Dispatch(c, getByReference)
		user = getByReference.Result

		if err != nil && errors.Cause(err) == app.ErrNotFound {
			if action.Email != "" {
				getByEmail := &query.GetUserByEmail{Email: action.Email}
				err = bus.Dispatch(c, getByEmail)
				user = getByEmail.Result
			}
			if err != nil && errors.Cause(err) == app.ErrNotFound {
				user = &entity.User{
					Tenant: c.Tenant(),
					Name:   action.Name,
					Email:  action.Email,
					Role:   enum.RoleVisitor,
				}
				err = bus.Dispatch(c, &cmd.RegisterUser{User: user})
			}
		}

		if err != nil {
			return c.Failure(err)
		}

		if action.Reference != "" && !user.HasProvider("reference") {
			if err := bus.Dispatch(c, &cmd.RegisterUserProvider{
				UserID:       user.ID,
				ProviderName: "reference",
				ProviderUID:  action.Reference,
			}); err != nil {
				return c.Failure(err)
			}
		}

		return c.Ok(web.Map{
			"id": user.ID,
		})
	}
}

// DeleteUser erases a user's personal data
func DeleteUser() web.HandlerFunc {
	return func(c *web.Context) error {
		userID, err := c.ParamAsInt("userID")
		if err != nil {
			return c.NotFound()
		}

		getUser := &query.GetUserByID{UserID: userID, TenantID: c.Tenant().ID}
		if err := bus.Dispatch(c, getUser); err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				return c.NotFound()
			}
			return c.Failure(err)
		}

		if err := bus.Dispatch(c, &cmd.DeleteUser{UserID: userID}); err != nil {
			return c.Failure(err)
		}

		// Demoting removes them from the userlist
		if env.Config.UserList.Enabled {
			c.Enqueue(tasks.UserListAddOrRemoveUser(userID, enum.RoleVisitor))
		}

		return c.Ok(web.Map{})
	}
}
