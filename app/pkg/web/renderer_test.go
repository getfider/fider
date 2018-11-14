package web_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log/noop"
	"github.com/getfider/fider/app/pkg/web"
)

func compareRendererResponse(buf *bytes.Buffer, fileName string, ctx *web.Context) {
	bytes, err := ioutil.ReadFile(env.Path(fileName))
	Expect(err).IsNil()
	Expect(strings.Replace(buf.String(), ctx.ContextID(), "CONTEXT_ID", -1)).Equals(string(bytes))
}

func TestRenderer_Basic(t *testing.T) {
	RegisterT(t)

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	renderer := web.NewRenderer(&models.SystemSettings{}, noop.NewLogger())
	renderer.Render(buf, "index.html", web.Props{}, ctx)
	compareRendererResponse(buf, "/app/pkg/web/testdata/basic.html", ctx)
}

func TestRenderer_Tenant(t *testing.T) {
	RegisterT(t)

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	ctx.SetTenant(&models.Tenant{Name: "Game of Thrones"})
	renderer := web.NewRenderer(&models.SystemSettings{}, noop.NewLogger())
	renderer.Render(buf, "index.html", web.Props{}, ctx)
	compareRendererResponse(buf, "/app/pkg/web/testdata/tenant.html", ctx)
}

func TestRenderer_Props(t *testing.T) {
	RegisterT(t)

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	renderer := web.NewRenderer(&models.SystemSettings{}, noop.NewLogger())
	renderer.Render(buf, "index.html", web.Props{
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

	buf := new(bytes.Buffer)
	ctx := newGetContext("https://demo.test.fider.io:3000/", nil)
	ctx.SetUser(&models.User{
		Name:   "Jon Snow",
		Email:  "jon.snow@got.com",
		Status: models.UserActive,
		Role:   models.RoleAdministrator,
	})
	renderer := web.NewRenderer(&models.SystemSettings{}, noop.NewLogger())
	renderer.Render(buf, "index.html", web.Props{
		Title:       "My Page Title",
		Description: "My Page Description",
	}, ctx)

	compareRendererResponse(buf, "/app/pkg/web/testdata/user.html", ctx)
}
