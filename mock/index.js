const Mock = require('mockjs')
const utils = require('./utils')

module.exports = function (app) {
    app.get('/history/user', (rep, res)=>{
        var result = []
        for (let index = 0; index < 20; index++) {
            var json = utils.getJsonFile('./mock_his.json')
            result[index]= json
        }
        res.json(Mock.mock(result))
    })
}