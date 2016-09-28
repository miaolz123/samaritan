var cooking = require('cooking');
var path = require('path');

cooking.set({
  entry: {
    app: './src/app.js',
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
  publicPath: './dist/',
  assetsPath: 'static',
  sourceMap: true,
  extractCSS: true,
  urlLoaderLimit: 10000,

  extends: ['react', 'lint']
});

cooking.add('resolve.alias', {
  'src': path.join(__dirname, 'src')
});

module.exports = cooking.resolve();
