import { initPinterestService, searchOnPinterest, loginToPinterest } from './pinterest.robot'

export const appLogin = async (ctx) => {
  try {
    await initPinterestService()
    const isLogin = await loginToPinterest()

    // 返回登入成功的響應
    ctx.body = {
      ok: isLogin,
      code: isLogin ? 200 : 99,
      data: isLogin,
      message: isLogin ? 'Login success' : 'Login failed',
    }
  } catch (error) {
    // 登入失敗
    ctx.body = {
      ok: false,
      code: 500,
      data: null,
      message: 'Error logging in: ' + error.message,
    }
  }
}

export const appSearch = async (ctx) => {
  const { key } = ctx.query

  try {
    const payload = await searchOnPinterest(key)

    ctx.body = {
      ok: payload.ok,
      code: payload.ok ? 200 : 99,
      data: payload,
      message: payload.ok ? 'Got data' : 'No data',
    }
  } catch (error) {
    // 如果出現錯誤，返回錯誤的響應
    ctx.body = {
      ok: false,
      code: 99,
      data: null,
      message: 'Error search: ' + error.message,
    }
  }
}
