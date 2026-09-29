package web_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/getfider/fider/app/models/entity"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func TestReactRenderer_FileNotFound(t *testing.T) {
	RegisterT(t)

	r, err := web.NewReactRenderer("unknown.js")
	Expect(err).IsNotNil()
	Expect(r).IsNil()
}

func TestReactRenderer_EmptyFile(t *testing.T) {
	RegisterT(t)

	r, err := web.NewReactRenderer("/app/pkg/web/testdata/empty.js")
	Expect(err).IsNil()
	Expect(r).IsNotNil()

	u, _ := url.Parse("https://github.com")
	html, err := r.Render(u, web.Map{})
	Expect(html).Equals("")
	Expect(err).IsNil()
}

func TestReactRenderer_RenderEmptyHomeHTML(t *testing.T) {
	RegisterT(t)

	r, err := web.NewReactRenderer("ssr.js")
	Expect(err).IsNil()

	u, _ := url.Parse("https://demo.test.fider.io")
	html, err := r.Render(u, web.Map{
		"page": "Home/Home.page",
		"tenant": &entity.Tenant{
			Locale: "en",
		},
		"settings": web.Map{
			"locale": "en",
		},
		"props": web.Map{
			"posts":          make([]web.Map, 0),
			"tags":           make([]web.Map, 0),
			"countPerStatus": web.Map{},
		},
	})
	Expect(html).ContainsSubstring(`<div class="c-dev-banner">DEV</div>`)
	Expect(html).ContainsSubstring(`<button class="p-home__add-idea-btn">`)
	Expect(html).ContainsSubstring(`Enter your suggestion here...`)
	Expect(html).ContainsSubstring(`What can we do better? This is the place for you to vote, discuss and share ideas.`)
	Expect(html).ContainsSubstring(`No posts have been created yet.`)
	Expect(html).ContainsSubstring(`Powered by Fider`)
	Expect(err).IsNil()
}

func TestReactRenderer_RenderWithMaliciousURL(t *testing.T) {
	RegisterT(t)

	r, err := web.NewReactRenderer("ssr.js")
	Expect(err).IsNil()

	// Attempt to break out of the JS string literal and inject code via the URL.
	// The URL contains a literal double-quote followed by JS that would, if interpolated
	// unsafely, execute server-side and emit attacker-controlled content into the SSR
	// output. After the fix the URL must be JSON-encoded into the ssrRender call so the
	// injection becomes inert.
	u := &url.URL{
		Scheme:   "https",
		Host:     "demo.test.fider.io",
		Path:     "/",
		RawQuery: `q=",{"page":"Error/Error404.page","settings":{},"tenant":{"locale":"en"}});` + "`<h1>RCE-PROOF</h1>`" + `//`,
	}
	html, err := r.Render(u, web.Map{
		"page": "Home/Home.page",
		"tenant": &entity.Tenant{
			Locale: "en",
		},
		"settings": web.Map{
			"locale": "en",
		},
		"props": web.Map{
			"posts":          make([]web.Map, 0),
			"tags":           make([]web.Map, 0),
			"countPerStatus": web.Map{},
		},
	})
	Expect(err).IsNil()
	Expect(strings.Contains(html, `<h1>RCE-PROOF</h1>`)).IsFalse()
}

func TestReactRenderer_RenderEmptyHomeHTML_Portuguese(t *testing.T) {
	RegisterT(t)

	r, err := web.NewReactRenderer("ssr.js")
	Expect(err).IsNil()

	u, _ := url.Parse("https://demo.test.fider.io")
	html, err := r.Render(u, web.Map{
		"page": "Home/Home.page",
		"tenant": &entity.Tenant{
			Locale: "pt-BR",
		},
		"settings": web.Map{
			"locale": "pt-BR",
		},
		"props": web.Map{
			"posts":          make([]web.Map, 0),
			"tags":           make([]web.Map, 0),
			"countPerStatus": web.Map{},
		},
	})
	Expect(html).ContainsSubstring(`<div class="c-dev-banner">DEV</div>`)
	Expect(html).ContainsSubstring(`<button class="p-home__add-idea-btn">`)
	Expect(html).ContainsSubstring(`Insira sua sugestão aqui...`)
	Expect(html).ContainsSubstring(`O que podemos fazer melhor? Este é o lugar para você votar, discutir e compartilhar ideias.`)
	Expect(html).ContainsSubstring(`Nenhuma postagem foi criada ainda.`)
	Expect(html).ContainsSubstring(`Powered by Fider`)
	Expect(err).IsNil()
}

func renderHomeWithVotes(t *testing.T, locale string) string {
	r, err := web.NewReactRenderer("ssr.js")
	Expect(err).IsNil()

	posts := make([]web.Map, 0)
	for i, votes := range []int{1, 3, 5} {
		posts = append(posts, web.Map{
			"id":            i + 1,
			"number":        i + 1,
			"slug":          "post",
			"title":         "Post",
			"description":   "",
			"status":        "open",
			"votesCount":    votes,
			"commentsCount": 0,
			"tags":          []string{},
			"user":          web.Map{"id": 1, "name": "Jon Snow"},
		})
	}

	u, _ := url.Parse("https://demo.test.fider.io")
	html, err := r.Render(u, web.Map{
		"page":     "Home/Home.page",
		"tenant":   &entity.Tenant{Locale: locale},
		"settings": web.Map{"locale": locale},
		"props": web.Map{
			"posts":          posts,
			"tags":           make([]web.Map, 0),
			"countPerStatus": web.Map{"open": len(posts)},
		},
	})
	Expect(err).IsNil()
	return html
}

// Vote counts use plural messages, which need Intl.PluralRules (missing in v8go)
func TestReactRenderer_RenderHomeWithPosts_Plurals(t *testing.T) {
	RegisterT(t)

	html := renderHomeWithVotes(t, "en")
	Expect(html).ContainsSubstring(`<span class="text-gray-700">Vote</span>`)
	Expect(html).ContainsSubstring(`<span class="text-gray-700">Votes</span>`)
}

func TestReactRenderer_RenderHomeWithPosts_Plurals_Polish(t *testing.T) {
	RegisterT(t)

	// Polish has distinct forms for 1 (one), 3 (few) and 5 (many)
	html := renderHomeWithVotes(t, "pl")
	Expect(html).ContainsSubstring(`<span class="text-gray-700">Głos</span>`)
	Expect(html).ContainsSubstring(`<span class="text-gray-700">Głosy</span>`)
	Expect(html).ContainsSubstring(`<span class="text-gray-700">Głosów</span>`)
}
