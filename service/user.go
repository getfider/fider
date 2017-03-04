package service

import (
	"database/sql"

	"github.com/WeCanHearYou/wchy/model"
)

// UserService is used for user operations
type UserService interface {
	GetByEmail(email string) (*model.User, error)
	Register(user *model.User) error
}

// PostgresUserService is used for user operations using a Postgres database
type PostgresUserService struct {
	DB *sql.DB
}

// GetByEmail returns a user based on given email
func (svc PostgresUserService) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	row := svc.DB.QueryRow("SELECT id, name, email FROM users WHERE email = $1", email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, ErrNotFound
	}

	return user, nil
}

// Register creates a new user based on given information
func (svc PostgresUserService) Register(user *model.User) error {
	tx, err := svc.DB.Begin()
	if err != nil {
		return err
	}

	if err = tx.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&user.ID); err != nil {
		tx.Rollback()
		return err
	}

	for _, provider := range user.Providers {
		if _, err = tx.Exec("INSERT INTO user_providers (user_id, provider, provider_uid) VALUES ($1, $2, $3)", user.ID, provider.Name, provider.UID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
