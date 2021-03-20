package web

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"rogchap.com/v8go"
)

type ReactRenderer struct {
	scriptPath    string
	scriptContent []byte
	ctx           *v8go.Context
}

func NewReactRenderer(scriptPath string) *ReactRenderer {
	bytes, _ := os.ReadFile(env.Path(scriptPath))
	isolate, err := v8go.NewIsolate()
	if err != nil {
		return &ReactRenderer{scriptPath: scriptPath}
	}

	v8ctx, _ := v8go.NewContext(isolate)
	_, err = v8ctx.RunScript(string(bytes), scriptPath)
	if err != nil {
		return &ReactRenderer{scriptPath: scriptPath}
	}

	return &ReactRenderer{ctx: v8ctx, scriptPath: scriptPath, scriptContent: bytes}
}

func (r *ReactRenderer) Render(u *url.URL, props Map) (string, error) {
	if len(r.scriptContent) == 0 {
		return "", nil
	}

	jsonArg, err := json.Marshal(props)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal props")
	}

	val, err := r.ctx.RunScript(`ssrRender("`+u.String()+`", "`+u.Path+`", `+string(jsonArg)+`)`, r.scriptPath)
	if err != nil {
		if jsErr, ok := err.(*v8go.JSError); ok {
			err = fmt.Errorf("%v", jsErr.StackTrace)
		}
		return "", errors.Wrap(err, "failed to execute ssrRender")
	}

	return val.String(), nil
}
