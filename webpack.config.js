const path = require("path");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

const publicFolder = path.resolve(__dirname, "public");
const isProduction = process.env.NODE_ENV === "production";

const plugins = [
    new MiniCssExtractPlugin({ 
        filename: "css/[name].[hash].css",
        chunkFilename: "css/[name].[hash].css"
    })
];

if (!isProduction) {
    const CleanObsoleteChunks = require("webpack-clean-obsolete-chunks");
    plugins.push(new CleanObsoleteChunks());
}

module.exports = {
    mode: process.env.NODE_ENV || "development",
    entry: "./public/index.tsx",
    output: {
        path: __dirname + "/dist",
        filename: "js/bundle.[chunkhash].js",
        publicPath: "/assets/",
    },
    devtool: "source-map",
    resolve: {
        extensions: [".ts", ".tsx", ".js"],
        alias: {
            "@fider": publicFolder
        }
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
                use: "ts-loader",
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
    plugins,
    stats: {
        children: false,
        modules: false,
    }
};