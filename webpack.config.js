const path = require("path");

const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');

const publicFolder = path.resolve(__dirname, "public");
const isProduction = process.env.NODE_ENV === "production";

const plugins = [
  new MiniCssExtractPlugin({ 
    filename: "css/[name].[contenthash].css",
    chunkFilename: "css/[name].[contenthash].css"
  }),
  new ForkTsCheckerWebpackPlugin({ checkSyntacticErrors: true })
];

if (!isProduction) {
  const CleanObsoleteChunks = require("webpack-clean-obsolete-chunks");
  plugins.push(new CleanObsoleteChunks());
}

module.exports = {
  mode: process.env.NODE_ENV || "development",
  entry: {
    main: "./public/index.tsx",
    vendor: [ "react", "react-dom", "tslib", "markdown-it", "react-textarea-autosize", "react-toastify", "react-loadable" ]
  },
  output: {
    path: __dirname + "/dist",
    filename: "js/[name].[contenthash].js",
    publicPath: "/assets/",
  },
  devtool: "source-map",
  resolve: {
    extensions: [".mjs", ".ts", ".tsx", ".js"],
    alias: {
      "@fider": publicFolder
    }
  },
  performance: {
    maxEntrypointSize: 250000,
    maxAssetSize: 250000,
    hints: 'error'
  },
  module: {
    rules: [
      { 
        test: /\.(css|scss)$/,
        use: [
          MiniCssExtractPlugin.loader,
          "css-loader",
          "sass-loader"                  
        ]
      },
      { 
        test: /\.(ts|tsx)$/,
        include: publicFolder,
        loader: "ts-loader",
        options: {
          transpileOnly: true
        }
      },
      {
        test: /\.(eot|svg|ttf|woff|woff2)$/,
        use: "file-loader?name=fonts/[name].[hash].[ext]"
      },
      {
        test: /\.(png|gif|jpg|jpeg)$/,
        use: "file-loader?name=images/[name].[hash].[ext]"
      }
    ]
  },
  optimization: {
    splitChunks: {
      cacheGroups: {
        vendor: {
          chunks: 'initial',
          name: 'vendor',
          test: 'vendor'
        },
      }
    }
  },
  plugins,
  stats: {
    children: false,
    modules: false,
  }
};