package web

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"sync"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"rogchap.com/v8go"
)

type ReactRenderer struct {
	scriptPath    string
	scriptContent []byte
	pool          *sync.Pool
}

func newIsolatePool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			isolate, err := v8go.NewIsolate()
			if err != nil {
				return errors.Wrap(err, "unable to initialize v8 isolate.")
			}

			runtime.SetFinalizer(isolate, func(iso *v8go.Isolate) {
				if iso != nil {
					iso.Dispose()
				}
			})

			return isolate
		},
	}
}

func NewReactRenderer(scriptPath string) (*ReactRenderer, error) {
	pool := newIsolatePool()
	bytes, err := os.ReadFile(env.Path(scriptPath))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read SSR script.")
	}

	return &ReactRenderer{pool: pool, scriptPath: scriptPath, scriptContent: bytes}, nil
}

func (r *ReactRenderer) Render(u *url.URL, props Map) (string, error) {
	if len(r.scriptContent) == 0 {
		return "", nil
	}

	item := r.pool.Get()
	isolate, ok := item.(*v8go.Isolate)
	if !ok {
		return "", item.(error)
	}
	defer r.pool.Put(isolate)

	v8ctx, err := v8go.NewContext(isolate)
	if err != nil {
		return "", errors.Wrap(err, "unable to initialize v8 context.")
	}
	defer v8ctx.Close()

	_, err = v8ctx.RunScript(string(r.scriptContent), r.scriptPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute SSR script.")
	}

	jsonArg, err := json.Marshal(props)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal props")
	}

	renderCmd := fmt.Sprintf(`ssrRender("%s", "%s", %s)`, u.String(), u.Path, string(jsonArg))
	val, err := v8ctx.RunScript(renderCmd, r.scriptPath)
	if err != nil {
		if jsErr, ok := err.(*v8go.JSError); ok {
			err = fmt.Errorf("%v", jsErr.StackTrace)
		}
		return "", errors.Wrap(err, "failed to execute ssrRender")
	}

	return val.String(), nil
}
