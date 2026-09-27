package markdown

import (
	"bytes"
	"context"
	stdhtml "html"
	"regexp"
	"strings"
	"sync"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/microcosm-cc/bluemonday"
)

// baselineURLPattern is the set of link and image destinations we accept without any
// tenant configuration. It mirrors the schemes DOMPurify permits in the browser (which is
// what public/services/markdown.ts relies on), plus relative references, with one
// hardening: no ASCII control characters or spaces anywhere in the value, so that
// "java\tscript:alert(1)" cannot slip through the relative-reference branch.
//
// The relative branch is expressed positively - "no ':' before the first '/', '?' or '#'" -
// rather than as a negated scheme match, so a scheme we have not thought about is rejected
// rather than accepted.
const baselineURLPattern = `(?:` +
	`(?i:https?|mailto|ftps?|tel|sms|callto|cid|xmpp):[^\x00-\x20\x7f]*` +
	`|[^:/?#\x00-\x20\x7f]*(?:[/?#][^\x00-\x20\x7f]*)?` +
	`)`

// blockedSchemes are never renderable, whatever the tenant has configured. This mirrors
// the unconditional /^javascript/i block in public/services/markdown.ts, extended to the
// rest of the script-capable and local-resource family.
var blockedSchemes = map[string]struct{}{
	"javascript":  {},
	"data":        {},
	"vbscript":    {},
	"livescript":  {},
	"mocha":       {},
	"view-source": {},
	"file":        {},
	"blob":        {},
	"jar":         {},
	"about":       {},
}

// dangerousProbes keep the tenant allow list purely additive. bluemonday can only
// allow-on-match - it has no deny list - so before folding an admin-supplied pattern into
// the policy we check it cannot accept anything from this set. A tenant administrator
// cannot re-open the hole by configuring "^javascript".
var dangerousProbes = []string{
	"javascript:alert(1)",
	"JaVaScRiPt:alert(1)",
	"data:text/html,<script>alert(1)</script>",
	"DATA:text/html;base64,AAAA",
	"vbscript:msgbox(1)",
	"livescript:x",
	"view-source:https://example.com/",
	"file:///etc/passwd",
}

var schemeRegex = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.\-]*:`)

// policy decides which destinations may be rendered, and sanitizes the final HTML.
type policy struct {
	allowURL *regexp.Regexp
	html     *bluemonday.Policy
}

// policyCache is keyed by the raw AllowedSchemes string, so it is bounded by the number of
// distinct tenant configurations rather than by the number of tenants - in practice almost
// every entry is the empty string.
var policyCache sync.Map

func policyFor(ctx context.Context) *policy {
	allowedSchemes := allowedSchemesFrom(ctx)
	if cached, ok := policyCache.Load(allowedSchemes); ok {
		return cached.(*policy)
	}
	p := compilePolicy(ctx, allowedSchemes)
	policyCache.Store(allowedSchemes, p)
	return p
}

// allowedSchemesFrom reads the tenant's extra schemes off the context. A missing tenant
// fails closed: baseline schemes only. The ALLOW_ALLOWED_SCHEMES kill switch is honoured
// at read time as well as at write time, because
// app/services/sqlstore/postgres/tenant.go only blanks the column on save.
func allowedSchemesFrom(ctx context.Context) string {
	if ctx == nil || !env.Config.AllowAllowedSchemes {
		return ""
	}
	if tenant, ok := ctx.Value(app.TenantCtxKey).(*entity.Tenant); ok && tenant != nil {
		return tenant.AllowedSchemes
	}
	return ""
}

func compilePolicy(ctx context.Context, allowedSchemes string) *policy {
	alternatives := []string{baselineURLPattern}

	for _, line := range strings.Split(allowedSchemes, "\n") {
		pattern := strings.TrimSpace(line)
		if pattern == "" {
			continue
		}

		// The browser compiles these with the "i" flag (public/hooks/use-fider.ts), so match that.
		compiled, err := regexp.Compile(`(?i:` + pattern + `)`)
		if err != nil {
			log.Warnf(ctx, "Ignoring invalid allowed scheme pattern @{Pattern}", dto.Props{
				"Pattern": pattern,
				"Error":   err.Error(),
			})
			continue
		}

		if matchesDangerousProbe(compiled) {
			log.Warnf(ctx, "Ignoring allowed scheme pattern @{Pattern} because it would permit a blocked scheme", dto.Props{
				"Pattern": pattern,
			})
			continue
		}

		// Tenant patterns are prefix patterns such as "^monero:[48]". Keep those prefix
		// semantics, but require the remainder to be free of control characters and spaces.
		alternatives = append(alternatives, `(?i:`+pattern+`)[^\x00-\x20\x7f]*`)
	}

	allowURL, err := regexp.Compile(`\A(?:` + strings.Join(alternatives, "|") + `)\z`)
	if err != nil {
		// Composition of already-compiled alternatives should not fail. Fail closed anyway.
		log.Warnf(ctx, "Falling back to baseline URL policy, could not compose allowed schemes", dto.Props{
			"Error": err.Error(),
		})
		allowURL = regexp.MustCompile(`\A(?:` + baselineURLPattern + `)\z`)
	}

	return &policy{allowURL: allowURL, html: newHTMLPolicy(allowURL)}
}

func matchesDangerousProbe(re *regexp.Regexp) bool {
	for _, probe := range dangerousProbes {
		if re.MatchString(probe) {
			return true
		}
	}
	return false
}

// allowsDestination reports whether a link or image destination may be rendered.
//
// The destination is normalized the way a browser would before deciding: the renderer's
// EscLink entity-decodes the destination on its way into the attribute, so
// "java&#115;cript:alert(1)" arrives here encoded but reaches the client decoded, and the
// ASCII whitespace and control characters that clients ignore inside a URL are stripped.
func (p *policy) allowsDestination(raw []byte) bool {
	destination := stdhtml.UnescapeString(string(raw))
	destination = strings.Map(func(r rune) rune {
		if r <= 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, destination)

	if destination == "" {
		return false
	}

	if scheme := schemeRegex.FindString(destination); scheme != "" {
		name := strings.ToLower(strings.TrimSuffix(scheme, ":"))
		if _, blocked := blockedSchemes[name]; blocked {
			return false
		}
	}

	return p.allowURL.MatchString(destination)
}

// codeLangName is the only shape of fenced-code info string we let through to the
// renderer, which builds class="language-..." out of it without escaping.
var codeLangName = regexp.MustCompile(`\A[\w.+#-]+\z`)

// safeCodeBlockInfo reduces a fenced-code info string to a plain language name, or to
// nothing at all when it is not one. It mirrors the renderer's own truncation at the first
// tab or space, so what we validate is what the renderer would emit - but note that the
// renderer does not truncate at the other characters HTML treats as whitespace, such as
// form feed and carriage return, which is what makes forging a second attribute possible.
func safeCodeBlockInfo(info []byte) []byte {
	if end := bytes.IndexAny(info, "\t "); end >= 0 {
		info = info[:end]
	}
	if len(info) == 0 || !codeLangName.Match(info) {
		return nil
	}
	return info
}

var (
	// codeLangClass is the sanitizer-side counterpart of codeLangName.
	codeLangClass = regexp.MustCompile(`\Alanguage-[\w.+#-]+\z`)
	// headingClass covers the Titleblock extension's <h1 class="title">.
	headingClass = regexp.MustCompile(`\A(?:title|special|title special)\z`)
	// linkRel covers the NofollowLinks / NoreferrerLinks renderer flags.
	linkRel = regexp.MustCompile(`\A(?:nofollow|noreferrer|noopener)(?: (?:nofollow|noreferrer|noopener))*\z`)
)

// newHTMLPolicy builds a bluemonday policy that permits exactly what gomarkdown can emit
// for the extensions and flags configured in markdown.go, and nothing else. It is the
// backstop that catches markup forged by the renderer itself rather than supplied as a
// link destination.
//
// Deliberately not built on UGCPolicy: that does not allow rel on anchors (so
// rel="nofollow noreferrer" would be rewritten to rel="nofollow"), it enables
// RequireParseableURLs which round-trips every URL through url.URL.String(), and it
// constrains alt and title with bluemonday.Paragraph, silently dropping alt text
// containing a colon, ampersand or quote.
func newHTMLPolicy(allowURL *regexp.Regexp) *bluemonday.Policy {
	p := bluemonday.NewPolicy()

	p.AllowElements(
		"p", "br", "hr",
		"h1", "h2", "h3", "h4", "h5", "h6",
		"strong", "em", "del", "code", "pre", "blockquote", "tt", "sup", "sub",
		"ul", "ol", "li", "dl", "dt", "dd",
		"table", "thead", "tbody", "tfoot", "tr", "th", "td", "caption",
	)

	p.AllowAttrs("class").Matching(codeLangClass).OnElements("code")
	p.AllowAttrs("class").Matching(headingClass).OnElements("h1", "h2", "h3", "h4", "h5", "h6")

	// allowURL is anchored, so it also rejects any value containing a space or control character.
	p.AllowAttrs("href").Matching(allowURL).OnElements("a")
	p.AllowAttrs("src").Matching(allowURL).OnElements("img")

	p.AllowAttrs("rel").Matching(linkRel).OnElements("a")
	p.AllowAttrs("alt").OnElements("img")
	p.AllowAttrs("title").OnElements("a", "img")

	p.AllowAttrs("align").Matching(bluemonday.CellAlign).OnElements("th", "td")
	p.AllowAttrs("colspan").Matching(bluemonday.Integer).OnElements("th", "td")

	return p
}
