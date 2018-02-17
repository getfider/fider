package markdown_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/markdown"
	. "github.com/onsi/gomega"
)

func TestParseMarkdown(t *testing.T) {
	RegisterTestingT(t)

	for input, expected := range map[string]string{
		"# Hello World": `<h1>Hello World</h1>
`,
		"![](http://example.com/hello.jpg)": `<p></p>
`,
		"Go to http://example.com/hello.jpg": `<p>Go to <a href="http://example.com/hello.jpg">http://example.com/hello.jpg</a></p>
`,
		`Can you try this?

- **Option 1**
- *Option 2*
- ~~Option 3~~
`: `<p>Can you try this?</p>

<ul>
<li><strong>Option 1</strong></li>
<li><em>Option 2</em></li>
<li><del>Option 3</del></li>
</ul>
`,
	} {
		output := markdown.Parse(input)
		Expect(output).To(Equal(expected))
	}
}
