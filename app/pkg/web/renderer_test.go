package web_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/enum"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

func compareRendererResponse(buf *bytes.Buffer, fileName string, ctx *web.Context) {
	// ioutil.WriteFile(env.Path(fileName), []byte(strings.Replace(buf.String(), ctx.ContextID(), "CONTEXT_ID", -1)), 0744)
	bytes, err := ioutil.ReadFile(env.Path(fileName))
	Expect(err).IsNil()
	Expect(strings.Replace(buf.String(), ctx.ContextID(), "CONTEXT_ID", -1)).Equals(string(bytes))
}

func TestRenderer_Basic(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveOAuthProviders) error {
		return nil
	})

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	renderer := web.NewRenderer(&models.SystemSettings{})
	renderer.Render(buf, http.StatusOK, "index.html", web.Props{}, ctx)
	compareRendererResponse(buf, "/app/pkg/web/testdata/basic.html", ctx)
}

func TestRenderer_WithChunkPreload(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveOAuthProviders) error {
		return nil
	})

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	renderer := web.NewRenderer(&models.SystemSettings{})
	renderer.Render(buf, http.StatusOK, "index.html", web.Props{ChunkName: "Test.page"}, ctx)
	compareRendererResponse(buf, "/app/pkg/web/testdata/chunk.html", ctx)
}

func TestRenderer_Tenant(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveOAuthProviders) error {
		return nil
	})

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	ctx.SetTenant(&models.Tenant{Name: "Game of Thrones"})
	renderer := web.NewRenderer(&models.SystemSettings{})
	renderer.Render(buf, http.StatusOK, "index.html", web.Props{}, ctx)
	compareRendererResponse(buf, "/app/pkg/web/testdata/tenant.html", ctx)
}

func TestRenderer_WithCanonicalURL(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveOAuthProviders) error {
		return nil
	})

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	ctx.SetCanonicalURL("http://feedback.demo.org")
	renderer := web.NewRenderer(&models.SystemSettings{})
	renderer.Render(buf, http.StatusOK, "index.html", web.Props{}, ctx)
	compareRendererResponse(buf, "/app/pkg/web/testdata/canonical.html", ctx)
}

func TestRenderer_Props(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveOAuthProviders) error {
		return nil
	})

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	renderer := web.NewRenderer(&models.SystemSettings{})
	renderer.Render(buf, http.StatusOK, "index.html", web.Props{
		Title:       "My Page Title",
		Description: "My Page Description",
		Data: web.Map{
			"number": 2,
			"array":  []string{"1", "2"},
			"object": web.Map{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}, ctx)

	compareRendererResponse(buf, "/app/pkg/web/testdata/props.html", ctx)
}

func TestRenderer_AuthenticatedUser(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveOAuthProviders) error {
		return nil
	})

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	ctx.SetUser(&models.User{
		ID:         5,
		Name:       "Jon Snow",
		Email:      "jon.snow@got.com",
		Status:     enum.UserActive,
		Role:       enum.RoleAdministrator,
		AvatarType: enum.AvatarTypeGravatar,
		AvatarURL:  "https://demo.test.fider.io:3000/avatars/gravatar/5/Jon%20Snow",
	})
	renderer := web.NewRenderer(&models.SystemSettings{})
	renderer.Render(buf, http.StatusOK, "index.html", web.Props{
		Title:       "My Page Title",
		Description: "My Page Description",
	}, ctx)

	compareRendererResponse(buf, "/app/pkg/web/testdata/user.html", ctx)
}

func TestRenderer_WithOAuth(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.ListActiveOAuthProviders) error {
		q.Result = []*dto.OAuthProviderOption{
			&dto.OAuthProviderOption{
				Provider:         app.GoogleProvider,
				DisplayName:      "Google",
				ClientID:         "1234",
				URL:              "https://demo.test.fider.io:3000/oauth/google",
				CallbackURL:      "https://demo.test.fider.io:3000/oauth/google/callback",
				LogoBlobKey:      "google.png",
				IsCustomProvider: false,
				IsEnabled:        true,
			},
		}
		return nil
	})

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	renderer := web.NewRenderer(&models.SystemSettings{})
	renderer.Render(buf, http.StatusOK, "index.html", web.Props{}, ctx)
	compareRendererResponse(buf, "/app/pkg/web/testdata/oauth.html", ctx)
}

func TestRenderer_NonOK(t *testing.T) {
	RegisterT(t)

	// it should not dispatch query.ListActiveOAuthProviders
	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	renderer := web.NewRenderer(&models.SystemSettings{})
	renderer.Render(buf, http.StatusNotFound, "index.html", web.Props{}, ctx)
	renderer.Render(buf, http.StatusBadRequest, "index.html", web.Props{}, ctx)
	renderer.Render(buf, http.StatusTemporaryRedirect, "index.html", web.Props{}, ctx)
}
