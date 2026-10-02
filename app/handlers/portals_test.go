package handlers_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/mock"
)

func TestPortalDirectory_Disabled(t *testing.T) {
	RegisterT(t)

	server := mock.NewServer()
	env.Config.PortalDirectoryEnabled = false

	code, _ := server.
		WithURL("http://test.fider.io/").
		Execute(handlers.PortalDirectory())

	Expect(code).Equals(http.StatusNotFound)
}

func TestPortalDirectory_ListsPublicPortals(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPublicTenants) error {
		q.Result = []*entity.Tenant{
			{ID: 2, Name: "Avengers", Subdomain: "avengers", Status: enum.TenantActive, LogoBlobKey: "logos/avengers.png"},
			{ID: 1, Name: "Demonstration", Subdomain: "demo", Status: enum.TenantActive},
		}
		return nil
	})

	server := mock.NewServer()
	env.Config.PortalDirectoryEnabled = true
	defer func() { env.Config.PortalDirectoryEnabled = false }()

	code, props := server.
		WithURL("http://test.fider.io/").
		ExecuteAsPage(handlers.PortalDirectory())

	Expect(code).Equals(http.StatusOK)

	portals, ok := props.Data["portals"].([]any)
	Expect(ok).IsTrue()
	Expect(portals).HasLen(2)

	avengers := portals[0].(map[string]any)
	Expect(avengers["name"]).Equals("Avengers")
	Expect(avengers["url"]).Equals("http://avengers.test.fider.io")
	Expect(avengers["host"]).Equals("avengers.test.fider.io")
	Expect(avengers["logoURL"]).Equals("http://avengers.test.fider.io/static/images/logos/avengers.png?size=200")

	demo := portals[1].(map[string]any)
	Expect(demo["name"]).Equals("Demonstration")
	Expect(demo["url"]).Equals("http://demo.test.fider.io")
	Expect(demo["host"]).Equals("demo.test.fider.io")
	Expect(demo["logoURL"]).IsNil()
}

func TestPortalDirectory_NoPortals(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPublicTenants) error {
		q.Result = []*entity.Tenant{}
		return nil
	})

	server := mock.NewServer()
	env.Config.PortalDirectoryEnabled = true
	defer func() { env.Config.PortalDirectoryEnabled = false }()

	code, props := server.
		WithURL("http://test.fider.io/").
		ExecuteAsPage(handlers.PortalDirectory())

	Expect(code).Equals(http.StatusOK)

	portals, ok := props.Data["portals"].([]any)
	Expect(ok).IsTrue()
	Expect(portals).HasLen(0)
}

func TestPortalDirectory_ServerRendersForCrawlers(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.GetPublicTenants) error {
		q.Result = []*entity.Tenant{
			{ID: 2, Name: "Avengers", Subdomain: "avengers", Status: enum.TenantActive},
		}
		return nil
	})

	server := mock.NewServer()
	env.Config.PortalDirectoryEnabled = true
	defer func() { env.Config.PortalDirectoryEnabled = false }()

	code, response := server.
		WithURL("http://test.fider.io/").
		AddHeader("User-Agent", "Googlebot/2.1").
		Execute(handlers.PortalDirectory())

	Expect(code).Equals(http.StatusOK)

	// The card markup only appears when the page was rendered server-side: the client bundle
	// never runs in a test. This is what a crawler indexing the root domain receives, and it
	// only works while the page stays registered in public/ssr.tsx.
	Expect(response.Body.String()).ContainsSubstring("c-portal-card")
}
