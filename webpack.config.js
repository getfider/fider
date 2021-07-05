/* eslint-disable no-undef */
/* eslint-disable @typescript-eslint/no-var-requires */
const path = require("path")
const glob = require("glob")

const ForkTsCheckerWebpackPlugin = require("fork-ts-checker-webpack-plugin")
const PurgecssPlugin = require("purgecss-webpack-plugin")
const SpriteLoaderPlugin = require("svg-sprite-loader/plugin")
const MiniCssExtractPlugin = require("mini-css-extract-plugin")
const BundleAnalyzerPlugin = require("webpack-bundle-analyzer").BundleAnalyzerPlugin
const publicFolder = path.resolve(__dirname, "public")
const localeFolder = path.resolve(__dirname, "locale")

const isProduction = process.env.NODE_ENV === "production"

const plugins = [
  new MiniCssExtractPlugin({
    filename: "css/[name].[contenthash].css",
    chunkFilename: "css/[name].[contenthash].css",
  }),
  new ForkTsCheckerWebpackPlugin(),
  new BundleAnalyzerPlugin({
    analyzerMode: "disabled", // To visualize the treemap of dependencies, change "disabled" to "static" and remove statsOptions
    generateStatsFile: true,
    statsFilename: "assets.json",
    statsOptions: {
      assets: false,
      children: false,
      chunks: false,
      entrypoints: true,
      chunkGroups: true,
      modules: false,
    },
  }),
  new PurgecssPlugin({
    paths: glob.sync(`./public/**/*.{html,tsx}`, { nodir: true }),
    defaultExtractor: (content) => content.match(/[^<>"'`\s]*[^<>"'`\s:]/g) || [],
    safelist: [/--/, /__/, /data-/],
  }),
  new SpriteLoaderPlugin({ plainSprite: true }),
]

// On Development Mode, we allow Assets to be up to 14 times bigger than on Production Mode
const maxSizeFactor = isProduction ? 1 : 14

module.exports = {
  mode: process.env.NODE_ENV || "development",
  entry: {
    main: "./public/index.tsx",
  },
  output: {
    path: __dirname + "/dist",
    filename: "js/[name].[contenthash].js",
    publicPath: "/assets/",
    clean: true,
  },
  devtool: "source-map",
  resolve: {
    extensions: [".mjs", ".ts", ".tsx", ".js", ".svg"],
    alias: {
      "@fider": publicFolder,
      "@locale": localeFolder,
    },
  },
  performance: {
    maxEntrypointSize: 368640 * maxSizeFactor, // 360 KiB. Should ideally be ~240 KiB
    maxAssetSize: 194560 * maxSizeFactor, // 190 KiB
    hints: "error",
  },
  module: {
    rules: [
      {
        test: /\.(css|scss)$/,
        use: [MiniCssExtractPlugin.loader, "css-loader", "sass-loader"],
      },
      {
        test: /\.(ts|tsx)$/,
        include: [publicFolder, localeFolder],
        loader: "babel-loader",
      },
      {
        test: /\.(json)$/,
        include: localeFolder,
        loader: "@lingui/loader",
        type: "javascript/auto",
      },
      {
        test: /\.svg$/,
        include: publicFolder,
        loader: "svg-sprite-loader",
        options: {
          esModule: false,
          extract: true,
          outputPath: "icons/",
          spriteFilename: "sprite.[hash].svg",
          publicPath: "/assets/icons/",
        },
      },
    ],
  },
  optimization: {
    moduleIds: "deterministic",
    runtimeChunk: "single",
    splitChunks: {
      cacheGroups: {
        common: {
          chunks: "all",
          name: "common",
          test: /[\\/]public[\\/](components|services|models)[\\/]/,
        },
        markdown: {
          chunks: "all",
          name: "markdown",
          test: /dompurify|marked/,
        },
        vendor: {
          chunks: "all",
          name: "vendor",
          test: /(react($|\/)|react-dom|tslib|react-textarea-autosize|@lingui\/core)/,
        },
      },
    },
  },
  plugins,
  stats: {
    assets: true,
    children: false,
    entrypoints: false,
    modules: false,
  },
}
