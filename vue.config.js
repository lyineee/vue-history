module.exports = {
  transpileDependencies: [
    'vuetify'
  ],
  publicPath: process.env.NODE_ENV === 'production'
    ? '/history/'
    : '/',
  devServer:{
    before: require('./mock/index'),
  },
  chainWebpack: config => {
    config
      .plugin('html')
      .tap(args => {
        args[0].title = "历史记录";
        return args
      })
  }
}
