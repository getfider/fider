/* eslint-disable @typescript-eslint/no-var-requires */
require("isomorphic-fetch")

require("ts-node").register({
  transpileOnly: true,
  compilerOptions: {
    target: "es6",
    strict: true,
    module: "commonjs",
  },
})
