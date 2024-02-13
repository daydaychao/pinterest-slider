import { initPinterestService, loginToPinterest } from '../pinterest.robot'

describe('Login functionality', () => {
  beforeAll(async () => {
    const browser = await initPinterestService()
    expect(browser).not.toBeNull() // 驗證 browser 實例不為空
  })

  test('Login test', async () => {
    const isLogin = await loginToPinterest()
    expect(isLogin).toBeTruthy() // 驗證 Login 為true
  }, 10000)
})
