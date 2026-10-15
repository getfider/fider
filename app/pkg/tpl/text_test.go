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

	result, err := tpl.Execute(tmpl, map[string]any{"description": "**Tom & Jerry** integration"})
	Expect(err).IsNil()
	Expect(result).Equals("Tom & Jerry integration")
}

func TestTextTemplate_Markdown_QuotedAsJSON(t *testing.T) {
	RegisterT(t)

	tmpl, err := tpl.GetTextTemplate("test", `{"m": {{ quote (markdown .description) }}}`)
	Expect(err).IsNil()

	result, err := tpl.Execute(tmpl, map[string]any{"description": "Tom & Jerry\n\nUse `List<String>` > \"Map\""})
	Expect(err).IsNil()

	var payload map[string]string
	Expect(json.Unmarshal([]byte(result), &payload)).IsNil()
	Expect(payload["m"]).Equals("Tom & Jerry\nUse List<String> > “Map”")
}
