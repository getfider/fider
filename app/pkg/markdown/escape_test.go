package markdown_test

import (
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/markdown"
)

func TestEscape(t *testing.T) {
	RegisterT(t)

	for input, expected := range map[string]string{
		"Jon Snow":                      "Jon Snow",
		"Zoë Ünïcode ✓":                 "Zoë Ünïcode ✓",
		"[Click](https://evil.example)": `\[Click\]\(https\:\/\/evil\.example\)`,
		"**bold** _it_ ~~del~~ `code`":  `\*\*bold\*\* \_it\_ \~\~del\~\~ \` + "`" + `code\` + "`",
		"![img](x.png)":                 `\!\[img\]\(x\.png\)`,
		"www.evil.example":              `www\.evil\.example`,
		"jon@evil.example":              `jon\@evil\.example`,
		"@[Jon Snow]":                   `\@\[Jon Snow\]`,
		"Tom & Jerry <3 <b>x</b>":       `Tom \& Jerry <3 <b\>x<\/b\>`,
		`a\b | c # d`:                   `a\\b \| c \# d`,
		"line1\nline2\r\n# heading":     `line1 line2  \# heading`,
	} {
		Expect(markdown.Escape(input)).Equals(expected)
	}
}
