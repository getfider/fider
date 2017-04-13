package handlers_test

import (
	"testing"

	"net/http"

	"net/url"

	"github.com/WeCanHearYou/wechy/app/handlers"
	"github.com/WeCanHearYou/wechy/app/mock"
	"github.com/WeCanHearYou/wechy/app/pkg/dbx"
	"github.com/WeCanHearYou/wechy/app/pkg/oauth"
	"github.com/WeCanHearYou/wechy/app/storage/postgres"
	. "github.com/onsi/gomega"
)

var (
	db            *dbx.Database
	oauthService  oauth.Service
	TenantStorage *postgres.TenantStorage
	UserStorage   *postgres.UserStorage
)

func setup() {
	db, _ = dbx.New()

	oauthService = &OAuthService{}
	TenantStorage = &postgres.TenantStorage{DB: db}
	UserStorage = &postgres.UserStorage{DB: db}
}

//OAuthService implements a mocked OAuthService
type OAuthService struct{}

//GetAuthURL returns authentication url for given provider
func (p *OAuthService) GetAuthURL(provider string, redirect string) string {
	return "http://orange.test.canherayou.com/oauth/token?provider=" + provider + "&redirect=" + redirect
}

//GetProfile returns user profile based on provider and code
func (p *OAuthService) GetProfile(provider string, code string) (*oauth.UserProfile, error) {
	if provider == "facebook" && code == "123" {
		return &oauth.UserProfile{
			ID:    "FB1234",
			Name:  "Jon Snow",
			Email: "jon.snow@got.com",
		}, nil
	}

	if provider == "google" && code == "123" {
		return &oauth.UserProfile{
			ID:    "GO1234",
			Name:  "Jon Snow",
			Email: "jon.snow@got.com",
		}, nil
	}

	if provider == "facebook" && code == "456" {
		return &oauth.UserProfile{
			ID:    "FB5678",
			Name:  "Some Facebook Guy",
			Email: "some.guy@facebook.com",
		}, nil
	}

	return nil, nil
}

func getHandlers() handlers.OAuthHandlers {
	setup()
	return handlers.OAuth(TenantStorage, oauthService, UserStorage)
}

func TestLoginHandler(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	code, response := server.ExecuteRaw(getHandlers().Login(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canherayou.com/oauth/token?provider=facebook&redirect="))
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=abc")
	code, _ := server.ExecuteRaw(getHandlers().Callback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusInternalServerError))
}

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.canhearyou.com")
	code, response := server.ExecuteRaw(getHandlers().Callback(oauth.FacebookProvider))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canhearyou.com"))
}

func TestCallbackHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://demo.test.canherayou.com/oauth/callback?state=http://demo.test.canhearyou.com&code=123")
	code, response := server.ExecuteRaw(getHandlers().Callback(oauth.FacebookProvider))

	Expect(db.Count("SELECT id FROM users WHERE email = 'jon.snow@got.com'")).To(Equal(1))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://demo.test.canhearyou.com?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjozMDAsInVzZXIvbmFtZSI6IkpvbiBTbm93IiwidXNlci9lbWFpbCI6Impvbi5zbm93QGdvdC5jb20ifQ.6_dLZrulH37ymBtqy-l7bhCti9hBv0lgEhH8tLm07CI"))
}

func TestCallbackHandler_NewUser(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://orange.test.canhearyou.com&code=456")
	code, response := server.ExecuteRaw(getHandlers().Callback(oauth.FacebookProvider))

	Expect(db.QueryInt("SELECT tenant_id FROM users WHERE email = 'some.guy@facebook.com'")).To(Equal(400))
	Expect(db.Exists("SELECT * FROM user_providers WHERE provider_uid = 'FB5678'")).To(BeTrue())

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://orange.test.canhearyou.com?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjoxLCJ1c2VyL25hbWUiOiJTb21lIEZhY2Vib29rIEd1eSIsInVzZXIvZW1haWwiOiJzb21lLmd1eUBmYWNlYm9vay5jb20ifQ.PGavs5a6HRotRozfXfNP39JPb0vSus_8LL9MAOeLGDs"))
}

func TestCallbackHandler_ExistingUser_NewProvider(t *testing.T) {
	RegisterTestingT(t)

	server := mock.NewServer()
	server.Context.Request().URL, _ = url.Parse("http://login.test.canherayou.com/oauth/callback?state=http://demo.test.canhearyou.com&code=123")
	code, response := server.ExecuteRaw(getHandlers().Callback(oauth.GoogleProvider))

	Expect(db.Count("SELECT id FROM users WHERE email = 'jon.snow@got.com'")).To(Equal(1))
	Expect(db.QueryString("SELECT provider_uid FROM user_providers WHERE user_id = 300 and provider = 'facebook'")).To(Equal("FB1234"))
	Expect(db.QueryString("SELECT provider_uid FROM user_providers WHERE user_id = 300 and provider = 'google'")).To(Equal("GO1234"))

	Expect(code).To(Equal(http.StatusTemporaryRedirect))
	Expect(response.Header.Get("Location")).To(Equal("http://demo.test.canhearyou.com?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyL2lkIjozMDAsInVzZXIvbmFtZSI6IkpvbiBTbm93IiwidXNlci9lbWFpbCI6Impvbi5zbm93QGdvdC5jb20ifQ.6_dLZrulH37ymBtqy-l7bhCti9hBv0lgEhH8tLm07CI"))
}
