package markdown

import "strings"

// Escape returns input with its markdown syntax neutralised, so that it renders
// as the same literal text when placed inline in a markdown string, e.g. a user
// name inside "**%s**". It targets CommonMark/GFM as rendered by the frontend
// (marked), which is what displays in-app notification titles. This package's
// own renderer honours fewer backslash escapes, so don't use it for content
// rendered by Full or PlainText.
//
// Every ASCII punctuation character is backslash-escaped, which CommonMark and
// GFM always render as the literal character. This covers emphasis, links,
// images, code spans, strikethrough, tables, entities ("&") and Fider mentions
// ("@[...]"), and also stops GFM from auto-linking bare URLs, "www." hosts and
// email addresses, because "://", "www." and "@" no longer appear contiguously.
//
// "<" is the exception: the frontend entity-encodes "<" before parsing, so inline
// HTML and "<...>" autolinks are already shown as plain text, and a "\<" would
// be displayed as "&lt;" rather than "<".
//
// Line breaks are replaced with spaces so the value cannot start a block
// element such as a heading, list or code block.
func Escape(input string) string {
	var sb strings.Builder
	sb.Grow(len(input) * 2)
	for _, r := range input {
		switch {
		case r == '\r' || r == '\n':
			sb.WriteByte(' ')
		case r == '<':
			sb.WriteRune(r)
		case r < 0x80 && isASCIIPunctuation(byte(r)):
			sb.WriteByte('\\')
			sb.WriteRune(r)
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func isASCIIPunctuation(c byte) bool {
	return (c >= '!' && c <= '/') || (c >= ':' && c <= '@') || (c >= '[' && c <= '`') || (c >= '{' && c <= '~')
}
