package markdown

import (
	"html/template"

	"github.com/russross/blackfriday"
)

//const full = md('commonmark', { html: false, breaks: true, linkify: true }).enable('linkify');
//const simple = md('commonmark', { html: false, breaks: true, linkify: true }).enable('linkify').disable('heading').disable('image');

// Parse given markdown input into html with all enabled features
func Parse(input string) template.HTML {
	mdExtns := 0
	mdExtns |= blackfriday.EXTENSION_TABLES
	mdExtns |= blackfriday.EXTENSION_AUTOLINK
	mdExtns |= blackfriday.EXTENSION_FENCED_CODE
	mdExtns |= blackfriday.EXTENSION_TITLEBLOCK
	mdExtns |= blackfriday.EXTENSION_STRIKETHROUGH
	mdExtns |= blackfriday.EXTENSION_DEFINITION_LISTS
	mdExtns |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	mdExtns |= blackfriday.EXTENSION_HARD_LINE_BREAK

	htmlExtns := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SKIP_IMAGES |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	renderer := blackfriday.HtmlRenderer(htmlExtns, "", "")
	output := blackfriday.Markdown([]byte(input), renderer, mdExtns)
	return template.HTML(output)
}
