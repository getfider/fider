package markdown

import (
	"context"
	"html/template"
	"io"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"

	htmlrenderer "github.com/gomarkdown/markdown/html"
	mdparser "github.com/gomarkdown/markdown/parser"
)

var mdExtns = 0 |
	mdparser.Tables |
	mdparser.Autolink |
	mdparser.FencedCode |
	mdparser.Titleblock |
	mdparser.Strikethrough |
	mdparser.DefinitionLists |
	mdparser.NoIntraEmphasis |
	mdparser.HardLineBreak

var htmlFlags = 0 |
	htmlrenderer.UseXHTML |
	htmlrenderer.Smartypants |
	htmlrenderer.SmartypantsFractions |
	htmlrenderer.SmartypantsDashes |
	htmlrenderer.SmartypantsLatexDashes |
	htmlrenderer.NofollowLinks |
	htmlrenderer.NoreferrerLinks

func createRenderer(handleImages bool, pol *policy) *htmlrenderer.Renderer {
	return htmlrenderer.NewRenderer(htmlrenderer.RendererOptions{
		Flags: htmlFlags,
		RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
			switch node := node.(type) {
			case *ast.HTMLSpan:
				htmlrenderer.EscapeHTML(w, node.Literal)
				return ast.GoToNext, true
			case *ast.HTMLBlock:
				_, _ = io.WriteString(w, "\n")
				htmlrenderer.EscapeHTML(w, node.Literal)
				_, _ = io.WriteString(w, "\n")
				return ast.GoToNext, true
			case *ast.CodeBlock:
				// The renderer interpolates the fenced-code info string straight into
				// class="language-..." without escaping it, and the parser accepts
				// arbitrary bytes there. An info string containing a quote can therefore
				// close the attribute and forge sibling attributes and elements. Reduce
				// the info string to a plain language name before the renderer sees it.
				node.Info = safeCodeBlockInfo(node.Info)
			case *ast.Link:
				// The renderer has no scheme filtering of its own unless the Safelink flag
				// is set, and that flag cannot be used here: it routes through
				// parser.IsSafeURL, which slices the destination before checking its
				// length and so panics on an empty link (see issue #1233). Gate the
				// destination ourselves instead. Writing nothing for the anchor while
				// still walking into the children drops the tags but keeps the link text.
				if !pol.allowsDestination(node.Destination) {
					return ast.GoToNext, true
				}
			case *ast.Image:
				// Skip images entirely when handleImages is false, or when the
				// destination is not renderable. Note that Safelink would not have
				// covered images in any case: imageEnter never consults it.
				if !handleImages || !pol.allowsDestination(node.Destination) {
					return ast.SkipChildren, true
				}
			}
			return ast.GoToNext, false
		},
	})
}

// Full turns a markdown into HTML using all rules.
//
// The context carries the tenant, whose AllowedSchemes setting extends the set of link
// destinations that may be rendered. A context without a tenant falls back to the baseline
// schemes, so forgetting to pass one fails closed.
func Full(ctx context.Context, input string, handleImages bool) template.HTML {
	pol := policyFor(ctx)

	// Apparently a parser cannot be reused.
	// https://github.com/gomarkdown/markdown/issues/229
	parser := mdparser.NewWithExtensions(mdExtns)

	renderer := createRenderer(handleImages, pol)

	output := markdown.ToHTML([]byte(input), parser, renderer)

	// Sanitize the rendered HTML as well as gating destinations above. The renderer can
	// forge markup on its own: it interpolates the fenced-code info string into
	// class="language-..." without escaping it, so the output is not trustworthy simply
	// because the input was parsed as markdown.
	return template.HTML(strings.TrimSpace(pol.html.Sanitize(string(output))))
}
