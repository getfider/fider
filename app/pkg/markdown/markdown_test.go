package markdown_test

import (
	"html/template"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/markdown"
)

func TestFullMarkdown(t *testing.T) {
	RegisterT(t)

	for input, expected := range map[string]string{
		"# Hello World":                      `<h1>Hello World</h1>`,
		"![](http://example.com/hello.jpg)":  `<p><img src="http://example.com/hello.jpg" alt="" /></p>`,
		"Go to http://example.com/hello.jpg": `<p>Go to <a href="http://example.com/hello.jpg">http://example.com/hello.jpg</a></p>`,
		`-123
-456
-789`: `<p>-123<br />
-456<br />
-789</p>`,
		`
- **Option 1**
- *Option 2*
- ~~Option 3~~`: `<ul>
<li><strong>Option 1</strong><br /></li>
<li><em>Option 2</em><br /></li>
<li><del>Option 3</del><br /></li>
</ul>`,
		`Please add:
– SEND_SMS
– RECEIVE_SMS
– READ_PHONE_STATE
This will allow to send and receive SMS and get the IMEI No. in our app.

Thanks!`: `<p>Please add:<br />
– SEND_SMS<br />
– RECEIVE_SMS<br />
– READ_PHONE_STATE<br />
This will allow to send and receive SMS and get the IMEI No. in our app.</p>

<p>Thanks!</p>`,
	} {
		output := markdown.Full(input)
		Expect(output).Equals(template.HTML(expected))
	}
}

func TestPlainTextMarkdown(t *testing.T) {
	RegisterT(t)

	for input, expected := range map[string]string{
		"**Hello World**":                           `Hello World`,
		"[My Link](http://example.com/)":            `My Link`,
		"![My Image](http://example.com/hello.jpg)": ``,
		"Go to http://example.com/hello.jpg":        `Go to http://example.com/hello.jpg`,
		"~~Option 3~~":                              `Option 3`,
		`-123
-456
-789`: `-123
-456
-789`,
		"# Hello World": `Hello World`,
		`# Hello World
How are you?`: `Hello World
How are you?`,
		`Hello World

How are you?`: `Hello World
How are you?`,
		"### Hello World":         `Hello World`,
		"Check this out: `HEEEY`": "Check this out: `HEEEY`",
	} {
		output := markdown.PlainText(input)
		Expect(output).Equals(expected)
	}
}

func TestSimpleMarkdown(t *testing.T) {
	RegisterT(t)

	for input, expected := range map[string]string{
		"**Hello World**":                           `<p><strong>Hello World</strong></p>`,
		"[My Link](http://example.com/)":            `<p><a href="http://example.com/">My Link</a></p>`,
		"![My Image](http://example.com/hello.jpg)": `<p></p>`,
		"Go to http://example.com/hello.jpg":        `<p>Go to <a href="http://example.com/hello.jpg">http://example.com/hello.jpg</a></p>`,
		"~~Option 3~~":                              `<p><del>Option 3</del></p>`,
		"# Hello World":                             `<p>Hello World</p>`,
		"### Hello World":                           `<p>Hello World</p>`,
		"Check this out: `HEEEY`":                   "<p>Check this out: <code>HEEEY</code></p>",
		`-123
-456
-789`: `<p>-123<br />
-456<br />
-789</p>`,
		`# Hello World
How are you?`: `<p>Hello World</p>

<p>How are you?</p>`,
		`Hello World

How are you?`: `<p>Hello World</p>

<p>How are you?</p>`,
	} {
		output := markdown.Simple(input)
		Expect(output).Equals(template.HTML(expected))
	}
}
