package markdown_test

import (
	"context"
	"strings"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/markdown"
)

func tenantCtx(allowedSchemes string) context.Context {
	return context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{AllowedSchemes: allowedSchemes})
}

func TestFullMarkdownBlocksDangerousSchemes(t *testing.T) {
	RegisterT(t)

	// An unrenderable link keeps its text and loses the anchor, which is how DOMPurify
	// behaves in the browser. An unrenderable image disappears entirely, rather than
	// leaving a broken-image icon behind in every notification email.
	for input, expected := range map[string]string{
		"[x](javascript:alert(1))":              `<p>x</p>`,
		"[x](JaVaScRiPt:alert(1))":              `<p>x</p>`,
		"[x](java&#115;cript:alert(1))":         `<p>x</p>`,
		"[x](&#106;avascript:alert(1))":         `<p>x</p>`,
		"[x](java\tscript:alert(1))":            `<p>x</p>`,
		"[x](vbscript:msgbox(1))":               `<p>x</p>`,
		"[x](data:text/html;base64,PHNjcmlwdD)": `<p>x</p>`,
		"[x](view-source:https://example.com)":  `<p>x</p>`,
		"[x](file:///etc/passwd)":               `<p>x</p>`,
		"![x](javascript:alert(1))":             `<p></p>`,
		"![x](data:text/html,<script>)":         `<p></p>`,

		// Safe destinations must keep working.
		"[x](http://example.com/a)":  `<p><a href="http://example.com/a" rel="nofollow noreferrer">x</a></p>`,
		"[x](https://example.com/a)": `<p><a href="https://example.com/a" rel="nofollow noreferrer">x</a></p>`,
		"[x](mailto:a@b.com)":        `<p><a href="mailto:a@b.com" rel="nofollow noreferrer">x</a></p>`,
		"[x](/posts/1)":              `<p><a href="/posts/1">x</a></p>`,
		"[x](#anchor)":               `<p><a href="#anchor">x</a></p>`,

		// An empty destination is not a link. Rendering it as a bare <a href=""> is what
		// made the Safelink renderer flag unusable here (issue #1233).
		"[link without actual link]()": `<p>link without actual link</p>`,
		"[]()":                         `<p></p>`,
	} {
		Expect(string(markdown.Full(context.Background(), input, true))).Equals(expected)
	}
}

func TestFullMarkdownTenantAllowedSchemes(t *testing.T) {
	RegisterT(t)

	const moneroLink = "[pay](monero:4AdUndXHHZ6cfufTMvppY6JwXNouMBzSkbLYfpAV5Usx3skxNgYeYTRJ5AmD5H3F)"
	configured := tenantCtx("^monero:[48]\n^bitcoin:(1|3|bc1)")

	// A scheme the tenant allows renders in emails and feeds, exactly as it does in the
	// web UI, so the backend does not silently drop links the admin opted into.
	Expect(string(markdown.Full(configured, moneroLink, true))).ContainsSubstring(`href="monero:4AdUnd`)
	Expect(string(markdown.Full(configured, "[pay](bitcoin:bc1qxyz)", true))).ContainsSubstring(`href="bitcoin:bc1qxyz"`)

	// Without that configuration, or for a scheme outside it, the link is not rendered.
	Expect(string(markdown.Full(context.Background(), moneroLink, true))).Equals(`<p>pay</p>`)
	Expect(string(markdown.Full(configured, "[pay](litecoin:LTC1abc)", true))).Equals(`<p>pay</p>`)

	// The allow list is additive only: a pattern that would re-enable a blocked scheme is
	// discarded rather than honoured.
	Expect(string(markdown.Full(tenantCtx("^javascript:"), "[x](javascript:alert(1))", true))).Equals(`<p>x</p>`)
	Expect(string(markdown.Full(tenantCtx("^data:"), "![x](data:text/html,<script>)", true))).Equals(`<p></p>`)

	// An invalid pattern is ignored, leaving the baseline intact rather than panicking.
	invalid := tenantCtx("^mon(ero:")
	Expect(string(markdown.Full(invalid, moneroLink, true))).Equals(`<p>pay</p>`)
	Expect(string(markdown.Full(invalid, "[x](https://example.com/a)", true))).ContainsSubstring(`href="https://example.com/a"`)

	// The instance-wide kill switch wins over tenant configuration.
	env.Config.AllowAllowedSchemes = false
	defer func() { env.Config.AllowAllowedSchemes = true }()
	Expect(string(markdown.Full(configured, moneroLink, true))).Equals(`<p>pay</p>`)
}

func TestFullMarkdownCodeFenceInfoString(t *testing.T) {
	RegisterT(t)

	// The renderer builds class="language-..." out of the fenced-code info string without
	// escaping it, and the parser accepts arbitrary bytes there. Note the form feed: the
	// renderer truncates the class at a tab or space but not at the other characters HTML
	// treats as whitespace, so that variant forges a real second attribute.
	for _, input := range []string{
		"```a\"><img/src=x/onerror=alert(1)>\ncode\n```",
		"```a\"><img\fsrc=x\fonerror=alert(1)>\ncode\n```",
		"```a\"><svg\fonload=alert(1)>\ncode\n```",
		"```a\"></code></pre><script>alert(1)</script>\ncode\n```",
		"```a\"><iframe\fsrc=javascript:alert(1)>\ncode\n```",
	} {
		output := string(markdown.Full(context.Background(), input, true))
		Expect(strings.Contains(output, "<img")).IsFalse()
		Expect(strings.Contains(output, "<svg")).IsFalse()
		Expect(strings.Contains(output, "<iframe")).IsFalse()
		Expect(strings.Contains(output, "<script")).IsFalse()
		Expect(strings.Contains(output, "onerror")).IsFalse()
		Expect(strings.Contains(output, "onload")).IsFalse()
		// The code itself must survive; this is a rendering fix, not a content filter.
		Expect(output).ContainsSubstring("code")
	}

	// A legitimate language is still emitted as a class.
	Expect(string(markdown.Full(context.Background(), "```go\nfmt.Println(1)\n```", true))).
		Equals("<pre><code class=\"language-go\">fmt.Println(1)\n</code></pre>")

	// As is the Titleblock extension's heading class.
	Expect(string(markdown.Full(context.Background(), "% Title\n\nbody", true))).
		ContainsSubstring(`<h1 class="title">Title</h1>`)
}
