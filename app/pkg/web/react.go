package web

import (
	"encoding/json"
	"os"

	"github.com/getfider/fider/app/pkg/errors"
	"rogchap.com/v8go"
)

type ReactRenderer struct {
	scriptPath    string
	scriptContent []byte
	isolate       *v8go.Isolate
}

func NewReactRenderer(scriptPath string) *ReactRenderer {
	bytes, _ := os.ReadFile(scriptPath)
	isolate, err := v8go.NewIsolate()
	if err != nil {
		return &ReactRenderer{isolate: nil, scriptPath: scriptPath}
	}

	return &ReactRenderer{isolate: isolate, scriptPath: scriptPath, scriptContent: bytes}
}

func (r *ReactRenderer) Render(urlPath string, props Map) (string, error) {
	if len(r.scriptContent) == 0 {
		return "", nil
	}

	v8ctx, _ := v8go.NewContext(r.isolate)
	_, err := v8ctx.RunScript(string(r.scriptContent), r.scriptPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse script")
	}

	jsonArg, err := json.Marshal(props)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal props")
	}

	val, err := v8ctx.RunScript(`ssrRender("`+urlPath+`", `+string(jsonArg)+`)`, r.scriptPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute ssrRender")
	}

	return val.String(), nil
}
