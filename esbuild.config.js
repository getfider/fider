/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable no-undef */
let emptyCSS = {
  name: "empty-css-imports",
  setup(build) {
    build.onLoad({ filter: /\.(css|scss)$/ }, () => ({ contents: "" }))
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
    plugins: [emptyCSS],
  })
  .catch(() => process.exit(1))
