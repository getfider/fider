/* eslint-disable @typescript-eslint/no-var-requires */
require("isomorphic-fetch")

process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0"

require("ts-node").register({
  transpileOnly: true,
  compilerOptions: {
    target: "es6",
    strict: true,
    module: "commonjs",
  },
})
