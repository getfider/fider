package web

import (
	"encoding/json"
	"os"

	"github.com/getfider/fider/app/pkg/errors"
	"rogchap.com/v8go"
)

type ReactRenderer struct {
	scriptContent []byte
	isolate       *v8go.Isolate
}

func NewReactRenderer() *ReactRenderer {
	bytes, _ := os.ReadFile("ssr.js")
	isolate, err := v8go.NewIsolate()
	if err != nil {
		return &ReactRenderer{isolate: nil}
	}

	return &ReactRenderer{isolate: isolate, scriptContent: bytes}
}

func (r *ReactRenderer) IsEnabled() bool {
	return len(r.scriptContent) > 0
}

func (r *ReactRenderer) Render(urlPath string, props Map) (string, error) {
	v8ctx, _ := v8go.NewContext(r.isolate)
	_, err := v8ctx.RunScript(string(r.scriptContent), "ssr.js")
	if err != nil {
		return "", errors.Wrap(err, "unable to run ssr.js")
	}

	jsonArg, err := json.Marshal(props)
	if err != nil {
		return "", errors.Wrap(err, "unable to masrhal props")
	}

	val, err = v8ctx.RunScript(`ssrRender("`+urlPath+`", `+string(jsonArg)+`)`, "ssr.js")
	if err != nil {
		return "", errors.Wrap(err, "failed to execute ssrRender")
	}

	return val.String(), nil
}
