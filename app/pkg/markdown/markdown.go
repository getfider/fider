package markdown

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"html/template"
	"io"
	"strings"

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
	htmlrenderer.SmartypantsLatexDashes

var fullRenderer = htmlrenderer.NewRenderer(htmlrenderer.RendererOptions{
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
		}
		return ast.GoToNext, false
	},
})

// Full turns a markdown into HTML using all rules
func Full(input string) template.HTML {
	// Apparently a parser cannot be reused.
	// https://github.com/gomarkdown/markdown/issues/229
	parser := mdparser.NewWithExtensions(mdExtns)
	output := markdown.ToHTML([]byte(input), parser, fullRenderer)
	return template.HTML(strings.TrimSpace(string(output)))
}
