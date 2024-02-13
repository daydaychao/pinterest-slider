const playwright = require('playwright')

;(async () => {
  console.log('%c[playwright] 啟動..', 'color: green')
  const browser = await playwright.chromium.launch()
  const page = (await browser.newContext()).newPage() //cookie不會共用

  // 導航到 Pinterest 網站
  await page.goto('https://www.pinterest.com/')

  // 點擊登入按鈕
  await page.click('data-test-id="simple-login-button"')

  await page.fill('input#email', 'chiachaichao@gmail.com')
  await page.fill('input#password', '1234567890d')
  await page.click('button.red.SignupButton.active[type="submit"]') // 點擊送出

  // 等待登入完成
  await page.waitForNavigation()

  // 獲取登入後的 JSON 資料
  const response = await page.waitForResponse((response) => {
    return response.url().startsWith('https://www.pinterest.com/resource/UserResource/get')
  })

  const data = await response.json()

  // 在控制台輸出 JSON 資料
  console.log(data)
  console.log('%c[playwright] 關閉', 'color: green')

  await browser.close()
})()
