package main

import (
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

var ignoreCssPlugin = api.Plugin{
	Name: "empty-css-imports",
	Setup: func(build api.PluginBuild) {
		build.OnLoad(api.OnLoadOptions{Filter: `\.(css|scss)$`},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				contents := ""
				return api.OnLoadResult{
					Contents: &contents,
				}, nil
			})
	},
}

func main() {
	nodeEnv := os.Getenv("NODE_ENV")
	if nodeEnv == "" {
		nodeEnv = "development"
	}

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{"../../public/ssr.tsx"},
		Outfile:     "../../ssr.js",
		Bundle:      true,
		Write:       true,
		Define: map[string]string{
			"process.env.NODE_ENV": `"` + nodeEnv + `"`,
		},
		Inject:  []string{"./esbuild-shim.js"},
		Plugins: []api.Plugin{ignoreCssPlugin},
	})

	if len(result.Errors) > 0 {
		for _, m := range result.Errors {
			println(m.Text)
		}
		os.Exit(1)
	}
}
