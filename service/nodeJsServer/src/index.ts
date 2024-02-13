import Koa from 'koa'
import { appLogin, appSearch } from './pinterest/pinterest.controller'
require('dotenv').config()

const app = new Koa()
const cors = require('@koa/cors')
const port = process.env.API_PORT

// 爬圖伺服器 pinterest抓圖機
app.use(cors({ origin: process.env.API_CORS }))

// 登入端點
app.use(async (ctx, next) => {
  if (ctx.path === '/login') await appLogin(ctx)
  else await next()
})

// 搜尋圖
app.use(async (ctx, next) => {
  if (ctx.path === '/get') await appSearch(ctx)
  else await next()
})

// 啟動伺服器並監聽指定端口
app.listen(port, async () => {
  console.log(`Server is running on port ${port}`)
})
