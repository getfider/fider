package identity_test

import (
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/postgres"
)

var (
	db            *dbx.Database
	oauthService  *OAuthService
	tenantService *postgres.TenantService
	userService   *postgres.UserService
)

func setup() {
	db, _ = dbx.New()

	oauthService = &OAuthService{}
	tenantService = &postgres.TenantService{DB: db}
	userService = &postgres.UserService{DB: db}
}
