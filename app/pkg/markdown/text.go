package markdown

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	htmlrenderer "github.com/gomarkdown/markdown/html"
	mdparser "github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"regexp"
	"strings"
)

var textRenderer = htmlrenderer.NewRenderer(htmlrenderer.RendererOptions{
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
		case *ast.Code:
			_, _ = io.WriteString(w, fmt.Sprintf("`%s`", node.Literal))
			return ast.GoToNext, true
		}
		return ast.GoToNext, false
	},
})

// The policy strips all HTML tags from the input text.
var strictPolicy = bluemonday.StrictPolicy()

// The regular expression finds duplicate newlines.
var regexNewlines = regexp.MustCompile(`\n+`)

// PlainText parses given markdown input and return only the text
func PlainText(input string) string {
	// Apparently a parser cannot be reused.
	// https://github.com/gomarkdown/markdown/issues/229
	parser := mdparser.NewWithExtensions(mdExtns)
	output := markdown.ToHTML([]byte(input), parser, textRenderer)
	sanitizedOutput := strictPolicy.Sanitize(string(output))
	sanitizedOutput = regexNewlines.ReplaceAllString(sanitizedOutput, "\n")
	return strings.TrimSpace(sanitizedOutput)
}
