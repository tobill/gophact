var path = require('path');
var webpack = require('webpack');

module.exports = {
        
    entry: { 
        app: [
            '@babel/polyfill',
            "./client/index.js",
        ],
    },
            
    output: {
        filename: "bundle.js",
        path: path.resolve(__dirname, "./static/js/"),
        publicPath: '/static/js/'
    },
    devServer: {
        contentBase: './static/js',
        port: 5000,
        open: false,
    },
    // Enable sourcemaps for debugging webpack's output.
    devtool: "source-map",
    // plugins: [
    //     new webpack.ProvidePlugin({
    //         "fetch": "imports-loader?this=>global!exports-loader?global.fetch!whatwg-fetch",
    //         fetch:   'imports?this=>global!exports?global.fetch!whatwg-fetch'
    //         fetch:   'imports-loader?this=>global!exports-loader?global.fetch!whatwg-fetch'
    //     }),
    // ],
    resolve: {
        // Add '.ts' and '.tsx' as resolvable extensions.
        extensions: [".js", ".json"]
    },
    mode: "development",
    module: {
        rules: [
            // All files with a '.ts' or '.tsx' extension will be handled by 'awesome-typescript-loader'.
            {
                test: /.js?$/,
                loader: 'babel-loader',
                exclude: /node_modules/,

            },
            // All output '.js' files will have any sourcemaps re-processed by 'source-map-loader'.
            { enforce: "pre", test: /\.js$/, loader: "source-map-loader" }
        ]
    },
};