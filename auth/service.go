package auth

import (
	"database/sql"
	"errors"
)

const (
	//OAuthFacebookProvider is const for 'facebook'
	OAuthFacebookProvider = "facebook"
	//OAuthGoogleProvider is const for 'google'
	OAuthGoogleProvider = "google"
)

//User represents an user inside our application
type User struct {
	ID        int64
	Name      string
	Email     string
	Providers []*UserProvider
}

//UserProvider represents the relashionship between an User and an Authentication provide
type UserProvider struct {
	Name string
	UID  string
}

//ErrUserNotFound is "User not found"
var ErrUserNotFound = errors.New("User not found")

// Service is used for auth operations
type Service interface {
	GetByEmail(email string) (*User, error)
	Register(user *User) error
}

// PostgresService is used for auth operations using a Postgres database
type PostgresService struct {
	DB *sql.DB
}

// GetByEmail returns a user based on given email
func (svc PostgresService) GetByEmail(email string) (*User, error) {
	user := &User{}
	row := svc.DB.QueryRow("SELECT id, name, email FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// Register creates a new user based on given information
func (svc PostgresService) Register(user *User) error {
	tx, err := svc.DB.Begin()
	if err != nil {
		return err
	}

	err = tx.QueryRow("INSERT INTO users (name, email) VALUES($1, $2) returning id;", user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO user_providers (user_id, provider, provider_uid) VALUES($1, $2, $3);", user.ID, user.Providers[0].Name, user.Providers[0].UID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
