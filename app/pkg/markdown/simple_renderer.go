package markdown

import (
	"bytes"

	"github.com/russross/blackfriday"
)

// Ensure we implement the Blackfriday Markdown Renderer interface
var _ blackfriday.Renderer = (*renderer)(nil)

// renderer renders Markdown to plain-text meant for listings and excerpts,
// and implements the blackfriday.Renderer interface.
//
// Many of the methods are stubs with no output to prevent output of HTML markup.
type proxy struct {
	renderer blackfriday.Renderer
}

// SimpleRenderer renders most HTML tags, except title
func SimpleRenderer(flags int) proxy {
	return proxy{
		renderer: blackfriday.HtmlRenderer(flags, "", ""),
	}
}

// Blocklevel callbacks

// BlockCode is the code tag callback.
func (p proxy) BlockCode(out *bytes.Buffer, text []byte, land string) {
	p.renderer.BlockCode(out, text, land)
}

// BlockQuote is the quote tag callback.
func (p proxy) BlockQuote(out *bytes.Buffer, text []byte) {
	p.renderer.BlockQuote(out, text)
}

// BlockHtml is the HTML tag callback.
func (p proxy) BlockHtml(out *bytes.Buffer, text []byte) {
	p.renderer.BlockHtml(out, text)
}

// Header is the header tag callback.
func (p proxy) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	p.Paragraph(out, text)
}

// HRule is the horizontal rule tag callback.
func (p proxy) HRule(out *bytes.Buffer) {
	p.renderer.HRule(out)
}

// List is the list tag callback.
func (p proxy) List(out *bytes.Buffer, text func() bool, flags int) {
	p.renderer.List(out, text, flags)
}

// ListItem is the list item tag callback.
func (p proxy) ListItem(out *bytes.Buffer, text []byte, flags int) {
	p.renderer.ListItem(out, text, flags)
}

// Paragraph is the paragraph tag callback.  This renders simple paragraph text
// into plain text, such that summaries can be easily generated.
func (p proxy) Paragraph(out *bytes.Buffer, text func() bool) {
	p.renderer.Paragraph(out, text)
}

// Table is the table tag callback.
func (p proxy) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	p.renderer.Table(out, header, body, columnData)
}

// TableRow is the table row tag callback.
func (p proxy) TableRow(out *bytes.Buffer, text []byte) {
	p.renderer.TableRow(out, text)
}

// TableHeaderCell is the table header cell tag callback.
func (p proxy) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {
	p.renderer.TableHeaderCell(out, text, flags)
}

// TableCell is the table cell tag callback.
func (p proxy) TableCell(out *bytes.Buffer, text []byte, flags int) {
	p.renderer.TableCell(out, text, flags)
}

// Footnotes is the foot notes tag callback.
func (p proxy) Footnotes(out *bytes.Buffer, text func() bool) {
	p.renderer.Footnotes(out, text)
}

// FootnoteItem is the footnote item tag callback.
func (p proxy) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	p.renderer.FootnoteItem(out, name, text, flags)
}

// TitleBlock is the title tag callback.
func (p proxy) TitleBlock(out *bytes.Buffer, text []byte) {
	p.renderer.TitleBlock(out, text)
}

// Spanlevel callbacks

// AutoLink is the autolink tag callback.
func (p proxy) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	p.renderer.AutoLink(out, link, kind)
}

// CodeSpan is the code span tag callback.  Outputs a simple Markdown version
// of the code span.
func (p proxy) CodeSpan(out *bytes.Buffer, text []byte) {
	p.renderer.CodeSpan(out, text)
}

// DoubleEmphasis is the double emphasis tag callback.  Outputs a simple
// plain-text version of the input.
func (p proxy) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	p.renderer.DoubleEmphasis(out, text)
}

// Emphasis is the emphasis tag callback.  Outputs a simple plain-text
// version of the input.
func (p proxy) Emphasis(out *bytes.Buffer, text []byte) {
	p.renderer.Emphasis(out, text)
}

// Image is the image tag callback.
func (p proxy) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	p.renderer.Image(out, link, title, alt)
}

// LineBreak is the line break tag callback.
func (p proxy) LineBreak(out *bytes.Buffer) {
	p.renderer.LineBreak(out)
}

// Link is the link tag callback.  Outputs a simple plain-text version
// of the input.
func (p proxy) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	p.renderer.Link(out, link, title, content)
}

// RawHtmlTag is the raw HTML tag callback.
func (p proxy) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	p.renderer.RawHtmlTag(out, tag)
}

// TripleEmphasis is the triple emphasis tag callback.  Outputs a simple plain-text
// version of the input.
func (p proxy) TripleEmphasis(out *bytes.Buffer, text []byte) {
	p.renderer.TripleEmphasis(out, text)
}

// StrikeThrough is the strikethrough tag callback.
func (p proxy) StrikeThrough(out *bytes.Buffer, text []byte) {
	p.renderer.StrikeThrough(out, text)
}

// FootnoteRef is the footnote ref tag callback.
func (p proxy) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	p.renderer.FootnoteRef(out, ref, id)
}

// Lowlevel callbacks

// Entity callback.  Outputs a simple plain-text version of the input.
func (p proxy) Entity(out *bytes.Buffer, entity []byte) {
	p.renderer.Entity(out, entity)
}

// NormalText callback.  Outputs a simple plain-text version of the input.
func (p proxy) NormalText(out *bytes.Buffer, text []byte) {
	p.renderer.NormalText(out, text)
}

// Header and footer

// DocumentHeader callback.
func (p proxy) DocumentHeader(out *bytes.Buffer) {
	p.renderer.DocumentHeader(out)
}

// DocumentFooter callback.
func (p proxy) DocumentFooter(out *bytes.Buffer) {
	p.renderer.DocumentFooter(out)
}

// GetFlags returns zero.
func (p proxy) GetFlags() int {
	return p.renderer.GetFlags()
}
