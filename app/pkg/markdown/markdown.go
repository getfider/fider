package markdown

import (
	"html/template"
	"strings"

	"github.com/russross/blackfriday"
)

var mdExtns = 0 | blackfriday.EXTENSION_TABLES |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_TITLEBLOCK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_DEFINITION_LISTS |
	blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	blackfriday.EXTENSION_HARD_LINE_BREAK

var htmlExtns = 0 |
	blackfriday.HTML_USE_XHTML |
	blackfriday.HTML_USE_SMARTYPANTS |
	blackfriday.HTML_SKIP_IMAGES |
	blackfriday.HTML_SMARTYPANTS_FRACTIONS |
	blackfriday.HTML_SMARTYPANTS_DASHES |
	blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

// Parse given markdown input into html with all enabled features
func Parse(input string) template.HTML {
	renderer := blackfriday.HtmlRenderer(htmlExtns, "", "")
	output := blackfriday.Markdown([]byte(input), renderer, mdExtns)

	return template.HTML(strings.TrimSpace(string(output)))
}

// PlainText parses given markdown input and return only the text
func PlainText(input string) string {
	renderer := TextRenderer()
	output := blackfriday.Markdown([]byte(input), renderer, mdExtns)
	return strings.TrimSpace(string(output))
}
