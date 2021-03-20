let plugin = {
  name: 'empty-css-imports',
  setup(build) {
    build.onLoad({ filter: /\.scss$/ }, () => ({ contents: '' }))
  },
}

require('esbuild').build({
  entryPoints: ['./public/ssr.tsx'],
  bundle: true,
  define: {
    "process.env.NODE_ENV": `"${process.env.NODE_ENV || 'development'}"`
  },
  inject: ['./global-shim.js'],
  outfile: 'ssr.js',
  plugins: [plugin],
}).catch(() => process.exit(1))
