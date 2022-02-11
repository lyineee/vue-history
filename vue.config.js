const webpack = require("webpack")
module.exports = {
  productionSourceMap: false,
  transpileDependencies: [
    'vuetify'
  ],
  publicPath: process.env.NODE_ENV === 'production'
    ? './'
    : './',
  devServer: {
    before: require('./mock/index'),
  },
  chainWebpack: config => {
    config
      .plugin('html')
      .tap(args => {
        args[0].title = "历史记录";
        return args
      })
    if (process.env.NODE_ENV == "production") {
      config
        .plugin('replace')
        .use(new webpack.NormalModuleReplacementPlugin(
          /environment\.ts/,
          'environment.prod.ts'
        ))
    }
  }
}
