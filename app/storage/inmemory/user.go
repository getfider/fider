package inmemory

import (
	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/models"
)

var users = []*models.User{}
var lastID = 0

// UserStorage is used for user operations
type UserStorage struct {
}

// GetByID returns a user based on given id
func (u *UserStorage) GetByID(userID int) (*models.User, error) {
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}
	return nil, app.ErrNotFound
}

// GetByEmail returns a user based on given email
func (u *UserStorage) GetByEmail(tenantID int, email string) (*models.User, error) {
	for _, user := range users {
		if user.Email == email && user.Tenant.ID == tenantID {
			return user, nil
		}
	}
	return nil, app.ErrNotFound
}

// Register creates a new user based on given information
func (u *UserStorage) Register(user *models.User) error {
	if user.ID == 0 {
		lastID = lastID + 1
		user.ID = lastID
	}
	users = append(users, user)
	return nil
}

// RegisterProvider adds given provider to userID
func (u *UserStorage) RegisterProvider(userID int, provider *models.UserProvider) error {
	for _, user := range users {
		if user.ID == userID {
			user.Providers = append(user.Providers, provider)
			return nil
		}
	}
	return app.ErrNotFound
}
