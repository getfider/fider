package tpl_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/getfider/fider/app/models/dto"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/tpl"
)

func TestGetTemplate_LayoutIncludesCustomFileAtEndOfBody(t *testing.T) {
	RegisterT(t)

	bf := new(bytes.Buffer)
	tmpl := tpl.GetTemplate("/views/base.html", "/views/index.html", "/app/pkg/tpl/testdata/custom.html")
	err := tpl.Render(context.Background(), tmpl, bf, dto.Props{
		"public": dto.Props{
			"contextID": "NONCE",
			"settings":  dto.Props{"googleAnalytics": ""},
		},
		"private": dto.Props{},
	})

	Expect(err).IsNil()
	Expect(bf.String()).ContainsSubstring(`<script nonce="NONCE" src="https://stats.example.com/t.js"></script>` + "\n\n</body>")
}

func TestGetTemplate_ExtraFilesAreCachedSeparately(t *testing.T) {
	RegisterT(t)

	const base, echo, custom = "app/pkg/tpl/testdata/base.html", "app/pkg/tpl/testdata/echo.html", "app/pkg/tpl/testdata/custom.html"

	Expect(tpl.GetTemplate(base, echo, custom).Lookup("custom.html") != nil).IsTrue()
	Expect(tpl.GetTemplate(base, echo).Lookup("custom.html") == nil).IsTrue()
}

func TestGetTemplate_Render(t *testing.T) {
	RegisterT(t)

	bf := new(bytes.Buffer)
	tmpl := tpl.GetTemplate("app/pkg/tpl/testdata/base.html", "app/pkg/tpl/testdata/echo.html")
	err := tpl.Render(context.Background(), tmpl, bf, dto.Props{
		"name": "John",
	})

	Expect(err).IsNil()
	Expect(bf.String()).Equals(`<html>
  <head>This goes on the head.</head>
  <body>
  Hello, John!
</body>
</html>`)
}
