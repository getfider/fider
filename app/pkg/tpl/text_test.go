package tpl_test

import (
	"encoding/json"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/tpl"
)

func TestTextTemplate_Markdown(t *testing.T) {
	RegisterT(t)

	tmpl, err := tpl.GetTextTemplate("test", "{{ markdown .description }}")
	Expect(err).IsNil()

	for input, expected := range map[string]string{
		"**Tom & Jerry** integration":     "Tom &amp; Jerry integration",
		"Support `<details>` blocks":      "Support &lt;details&gt; blocks",
		"Hey <!channel> look":             "Hey &lt;!channel&gt; look",
		"Hey &lt;!channel&gt; look":       "Hey &lt;!channel&gt; look",
		"Hey `<!channel>` look":           "Hey &lt;!channel&gt; look",
		"Hi <img src=x onerror=alert(1)>": "Hi &lt;img src=x onerror=alert(1)&gt;",
	} {
		result, err := tpl.Execute(tmpl, map[string]any{"description": input})
		Expect(err).IsNil()
		Expect(result).Equals(expected)
	}
}

func TestTextTemplate_Markdown_QuotedAsJSON(t *testing.T) {
	RegisterT(t)

	tmpl, err := tpl.GetTextTemplate("test", `{"m": {{ quote (markdown .description) }}}`)
	Expect(err).IsNil()

	result, err := tpl.Execute(tmpl, map[string]any{"description": "Tom & Jerry\n\nUse `List<String>` > <!channel>"})
	Expect(err).IsNil()

	var payload map[string]string
	Expect(json.Unmarshal([]byte(result), &payload)).IsNil()
	Expect(payload["m"]).Equals("Tom &amp; Jerry\nUse List&lt;String&gt; &gt; &lt;!channel&gt;")
}
