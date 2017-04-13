package identity_test

import (
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"
)

var (
	db            *dbx.Database
	oauthService  *OAuthService
	TenantStorage *postgres.TenantStorage
	UserStorage   *postgres.UserStorage
)

func setup() {
	db, _ = dbx.New()

	oauthService = &OAuthService{}
	TenantStorage = &postgres.TenantStorage{DB: db}
	UserStorage = &postgres.UserStorage{DB: db}
}
