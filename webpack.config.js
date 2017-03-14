const path = require('path');
const CleanWebpackPlugin = require('clean-webpack-plugin');
var AutoCleanBuildPlugin = require('webpack-auto-clean-build-plugin');

module.exports = {
    entry: "./public/index.ts",
    output: {
        path: __dirname + '/dist',
        filename: "js/bundle.[chunkhash].js",
        publicPath: "/assets/"
    },
    resolve: {
        extensions: ['.ts', '.tsx', '.js']
    },
    module: {
        loaders: [
            { test: /\.css$/, loader: "style-loader!css-loader" },
            { test: /\.(ts|tsx)$/, loader: "ts-loader" },
            {
                test: /\.(eot|svg|ttf|woff|woff2)$/,
                loader: 'file-loader?name=fonts/[name].[hash].[ext]'
            },
            {
                test: /\.(png|gif|jpg|jpeg)$/,
                loader: 'file-loader?name=images/[name].[hash].[ext]'
            }
        ]
    },
    plugins: [
        new AutoCleanBuildPlugin()
    ],
    watch: true
};