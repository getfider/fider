let plugin = {
  name: 'empty-css-imports',
  setup(build) {
    build.onLoad({ filter: /\.scss$/ }, () => ({ contents: '' }))
  },
}

require('esbuild').build({
  entryPoints: ['./public/ssr.tsx'],
  bundle: true,
  platform: 'node',
  define: {
    "process.env.NODE_ENV": "'production'"
  },
  outfile: 'ssr.js',
  plugins: [plugin],
}).catch(() => process.exit(1))
