/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */

const fs = require("fs")
const esbuild = require("esbuild")
const babel = require("@babel/core")

// Replace with NPM package when this is resolved: https://github.com/nativew/esbuild-plugin-babel/issues/7
const babelPlugin = (options = {}) => ({
  name: "babel",
  setup(build, { transform } = {}) {
    const { filter = /.*/, namespace = "", config = {} } = options

    const transformContents = ({ args, contents }) => {
      const babelOptions = babel.loadOptions({
        ...config,
        filename: args.path,
        caller: {
          name: "esbuild-plugin-babel",
          supportsStaticESM: true,
        },
      })
      if (!babelOptions) return { contents }

      if (babelOptions.sourceMaps) {
        const filename = path.relative(process.cwd(), args.path)

        babelOptions.sourceFileName = filename
      }

      return new Promise((resolve, reject) => {
        babel.transform(contents, babelOptions, (error, result) => {
          error ? reject(error) : resolve({ contents: result.code })
        })
      })
    }

    if (transform) return transformContents(transform)

    build.onLoad({ filter, namespace }, async (args) => {
      const contents = await fs.promises.readFile(args.path, "utf8")

      return transformContents({ args, contents })
    })
  },
})

const emptyCSS = {
  name: "empty-css-imports",
  setup(build) {
    build.onLoad({ filter: /\.(css|scss)$/ }, () => ({ contents: "" }))
  },
}

const emptySVG = {
  name: "empty-svg-imports",
  setup(build) {
    build.onLoad({ filter: /\.(svg)$/ }, async (args) => {
      let contents = await fs.promises.readFile(args.path, "utf8")
      let buff = Buffer.from(contents).toString("base64")
      return { contents: `data:image/svg+xml;base64,${buff}`, loader: "text" }
    })
  },
}

esbuild
  .build({
    entryPoints: ["./public/ssr.tsx"],
    bundle: true,
    define: {
      "process.env.NODE_ENV": `"${process.env.NODE_ENV || "development"}"`,
    },
    inject: ["./esbuild-shim.js"],
    outfile: "ssr.js",
    plugins: [emptyCSS, emptySVG, babelPlugin()],
  })
  .catch(() => process.exit(1))
