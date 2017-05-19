const path = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const CleanObsoleteChunks = require('webpack-clean-obsolete-chunks');

module.exports = {
    entry: "./public/index.tsx",
    output: {
        path: __dirname + '/dist',
        filename: "js/bundle.[chunkhash].js",
        publicPath: "/assets/"
    },
    devtool: "source-map",
    resolve: {
        extensions: ['.ts', '.tsx', '.js']
    },
    module: {
        rules: [
            { 
                test: /\.scss$/, 
                loader: ExtractTextPlugin.extract("css-loader!sass-loader")
            },
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
        new CleanObsoleteChunks(),
        new ExtractTextPlugin({ 
            filename: 'css/[name].[hash].css', 
            disable: false, 
            allChunks: true 
        })
    ],
    externals: {
        'react': 'React',
        'react-dom' : 'ReactDOM',
        "jquery": "jQuery"
    }
};