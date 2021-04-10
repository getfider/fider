/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */

const fs = require("fs")

let emptyCSS = {
  name: "empty-css-imports",
  setup(build) {
    build.onLoad({ filter: /\.(css|scss)$/ }, () => ({ contents: "" }))
  },
}

let emptySVG = {
  name: "empty-svg-imports",
  setup(build) {
    build.onLoad({ filter: /\.(svg)$/ }, async (args) => {
      let contents = await fs.promises.readFile(args.path, "utf8")
      let buff = Buffer.from(contents).toString("base64")
      return { contents: `data:image/svg+xml;base64,${buff}`, loader: "text" }
    })
  },
}

require("esbuild")
  .build({
    entryPoints: ["./public/ssr.tsx"],
    bundle: true,
    define: {
      "process.env.NODE_ENV": `"${process.env.NODE_ENV || "development"}"`,
    },
    inject: ["./esbuild-shim.js"],
    outfile: "ssr.js",
    plugins: [emptyCSS, emptySVG],
  })
  .catch(() => process.exit(1))
