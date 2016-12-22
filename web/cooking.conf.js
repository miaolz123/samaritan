var cooking = require('cooking');
var path = require('path');
var CopyWebpackPlugin = require('copy-webpack-plugin');

cooking.set({
  entry: {
    app: './src/index.js',
    vendor: ['react', 'react-dom']
  },
  dist: './dist',
  template: 'src/index.tpl',

  // development
  devServer: {
    hostname: '127.0.0.1',
    port: 8080,
    publicPath: '/'
  },

  // production
  clean: true,
  hash: true,
  chunk: 'vendor',
  publicPath: './',
  assetsPath: 'static',
  sourceMap: true,
  extractCSS: true,
  urlLoaderLimit: 10000,
  postcss: [],

  extends: ['react', 'lint', 'less']
});

cooking.add('resolve.alias', {
  'src': path.join(__dirname, 'src')
});

cooking.add('plugin.copy', new CopyWebpackPlugin([
  {
    from: 'node_modules/monaco-editor/min/vs',
    to: 'vs',
  }
]));

module.exports = cooking.resolve();
