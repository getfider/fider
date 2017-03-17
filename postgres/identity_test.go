package postgres_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/WeCanHearYou/wechy/identity"
	"github.com/WeCanHearYou/wechy/postgres"
	. "github.com/onsi/gomega"
)

func TestUserService_GetByEmail_Error(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := mock.ExpectQuery("SELECT id, name, email FROM users WHERE email = \\$1")
	query.WithArgs("jon.stark@got.com")
	query.WillReturnError(sql.ErrNoRows)

	svc := &postgres.UserService{DB: db}
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

	svc := &postgres.UserService{DB: db}
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

	query := mock.ExpectQuery("INSERT INTO users \\(name, email, created_on\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id")
	query.WithArgs("Jon Snow", "jon.snow@got.com", AnyTime{})
	query.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(300)))

	exec := mock.ExpectExec("INSERT INTO user_providers \\(user_id, provider, provider_uid\\) VALUES \\(\\$1, \\$2, \\$3\\)")
	exec.WithArgs(int64(300), identity.OAuthFacebookProvider, "123123123")
	exec.WillReturnResult(driver.ResultNoRows)

	mock.ExpectCommit()

	svc := &postgres.UserService{DB: db}
	user := &identity.User{
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
		Providers: []*identity.UserProvider{
			{
				UID:  "123123123",
				Name: identity.OAuthFacebookProvider,
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

	query := mock.ExpectQuery("INSERT INTO users \\(name, email, created_on\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id")
	query.WithArgs("Jon Snow", "jon.snow@got.com", AnyTime{})
	query.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(300)))

	exec := mock.ExpectExec("INSERT INTO user_providers \\(user_id, provider, provider_uid\\) VALUES \\(\\$1, \\$2, \\$3\\)")
	exec.WithArgs(int64(300), identity.OAuthFacebookProvider, "123123123")
	exec.WillReturnResult(driver.ResultNoRows)

	exec = mock.ExpectExec("INSERT INTO user_providers \\(user_id, provider, provider_uid\\) VALUES \\(\\$1, \\$2, \\$3\\)")
	exec.WithArgs(int64(300), identity.OAuthGoogleProvider, "456456456")
	exec.WillReturnResult(driver.ResultNoRows)

	mock.ExpectCommit()

	svc := &postgres.UserService{DB: db}
	user := &identity.User{
		Name:  "Jon Snow",
		Email: "jon.snow@got.com",
		Providers: []*identity.UserProvider{
			{
				UID:  "123123123",
				Name: identity.OAuthFacebookProvider,
			},
			{
				UID:  "456456456",
				Name: identity.OAuthGoogleProvider,
			},
		},
	}
	err := svc.Register(user)

	Expect(err).To(BeNil())
	Expect(user.ID).To(Equal(int64(300)))
}

func TestTenantService_GetByDomain_Error(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1").WithArgs("mydomain").WillReturnError(sql.ErrNoRows)

	svc := &postgres.TenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	Expect(tenant).To(BeNil())
	Expect(err).NotTo(BeNil())
}

func TestTenantService_GetByDomain_Subdomain(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "subdomain"}).AddRow(int64(234), "My Domain Inc.", "mydomain")
	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1 OR cname = \\$2").WithArgs("mydomain", "mydomain").WillReturnRows(rows)

	svc := &postgres.TenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain")

	Expect(tenant.ID).To(Equal(int64(234)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Domain).To(Equal("mydomain.test.canhearyou.com"))
	Expect(err).To(BeNil())
}

func TestTenantService_GetByDomain_FullDomain(t *testing.T) {
	RegisterTestingT(t)

	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "subdomain"}).AddRow(int64(234), "My Domain Inc.", "mydomain")
	mock.ExpectQuery("SELECT id, name, subdomain FROM tenants WHERE subdomain = \\$1 OR cname = \\$2").WithArgs("mydomain", "mydomain.anydomain.com").WillReturnRows(rows)

	svc := &postgres.TenantService{DB: db}
	tenant, err := svc.GetByDomain("mydomain.anydomain.com")

	Expect(tenant.ID).To(Equal(int64(234)))
	Expect(tenant.Name).To(Equal("My Domain Inc."))
	Expect(tenant.Domain).To(Equal("mydomain.test.canhearyou.com"))
	Expect(err).To(BeNil())
}
