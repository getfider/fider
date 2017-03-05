package service_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/WeCanHearYou/wchy/app/auth"
	"github.com/WeCanHearYou/wchy/app/model"
	"github.com/WeCanHearYou/wchy/app/service"
	. "github.com/onsi/gomega"
)

func TestUserService_GetByEmail_Error(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := mock.ExpectQuery("SELECT id, name, email FROM users WHERE email = \\$1")
	query.WithArgs("jon.stark@got.com")
	query.WillReturnError(sql.ErrNoRows)

	svc := &service.PostgresUserService{DB: db}
	user, err := svc.GetByEmail("jon.stark@got.com")

	Expect(user).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestUserService_GetByEmail(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(int64(234), "Jon Snow", "jon.snow@got.com")
	query := mock.ExpectQuery("SELECT id, name, email FROM users WHERE email = \\$1")
	query.WithArgs("jon.snow@got.com")
	query.WillReturnRows(rows)

	svc := &service.PostgresUserService{DB: db}
	user, err := svc.GetByEmail("jon.snow@got.com")

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int64(234)))
	Expect(user.Name).To(Equal("Jon Snow"))
	Expect(user.Email).To(Equal("jon.snow@got.com"))
}

func TestUserService_Register(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectBegin()

	query := mock.ExpectQuery("INSERT INTO users \\(name, email\\) VALUES \\(\\$1, \\$2\\) RETURNING id")
	query.WithArgs("Jon Snow", "jon.snow@got.com")
	query.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(300)))

	exec := mock.ExpectExec("INSERT INTO user_providers \\(user_id, provider, provider_uid\\) VALUES \\(\\$1, \\$2, \\$3\\)")
	exec.WithArgs(int64(300), auth.OAuthFacebookProvider, "123123123")
	exec.WillReturnResult(driver.ResultNoRows)

	mock.ExpectCommit()

	svc := &service.PostgresUserService{DB: db}
	user := &model.User{
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
		Providers: []*model.UserProvider{
			{
				UID:  "123123123",
				Name: auth.OAuthFacebookProvider,
			},
		},
	}
	err := svc.Register(user)

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int64(300)))
}

func TestUserService_Register_MultipleProviders(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectBegin()

	query := mock.ExpectQuery("INSERT INTO users \\(name, email\\) VALUES \\(\\$1, \\$2\\) RETURNING id")
	query.WithArgs("Jon Snow", "jon.snow@got.com")
	query.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(300)))

	exec := mock.ExpectExec("INSERT INTO user_providers \\(user_id, provider, provider_uid\\) VALUES \\(\\$1, \\$2, \\$3\\)")
	exec.WithArgs(int64(300), auth.OAuthFacebookProvider, "123123123")
	exec.WillReturnResult(driver.ResultNoRows)

	exec = mock.ExpectExec("INSERT INTO user_providers \\(user_id, provider, provider_uid\\) VALUES \\(\\$1, \\$2, \\$3\\)")
	exec.WithArgs(int64(300), auth.OAuthGoogleProvider, "456456456")
	exec.WillReturnResult(driver.ResultNoRows)

	mock.ExpectCommit()

	svc := &service.PostgresUserService{DB: db}
	user := &model.User{
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
		Providers: []*model.UserProvider{
			{
				UID:  "123123123",
				Name: auth.OAuthFacebookProvider,
			},
			{
				UID:  "456456456",
				Name: auth.OAuthGoogleProvider,
			},
		},
	}
	err := svc.Register(user)

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int64(300)))
}
