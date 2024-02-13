import { chromium, Browser } from 'playwright'
import type { Page } from 'playwright'
import { Color } from '../constants'

// Pinterest機器人 這裡的邏輯是為了pinterest訂製的

let browserInstance: Browser | undefined
let currentPage: Page | undefined

// init 機器人 有機器人才能抓圖做事
export const initPinterestService = async () => {
  if (!browserInstance) {
    console.log(Color.yellow + '[Pinterest機器人服務] 啟動..' + Color.close)
    browserInstance = await chromium.launch()
  }
  return browserInstance
}
// close 機器人
export const closePinterestService = async () => {
  if (browserInstance) {
    await browserInstance.close()
    console.log(Color.yellow + '[Pinterest機器人服務] 關閉..' + Color.close)
  }
}

// api: login
export const loginToPinterest = async (): Promise<boolean> => {
  if (!browserInstance) {
    throw new Error('Browser instance not initialized')
  }

  if (currentPage && (await checkIsLoginByCookie(currentPage))) {
    console.log('收到Login指示,但是機器人已經登入了')
    return true
  }

  const page = await browserInstance.newPage()
  currentPage = page
  console.log(Color.yellow + '[playwright] login..' + Color.close)

  await page.goto('https://www.pinterest.com/')
  console.log('goto OK')

  await page.click('[data-test-id="simple-login-button"]')
  console.log('login btn OK')

  await page.screenshot({ path: '/screenshot/1_index.png' })

  // 跳彈login表單
  await page.waitForSelector('[data-test-id="login-modal-default"]')
  console.log('wait login modal OK')

  await page.fill('input#email', process.env.INFO_USERNAME)
  console.log('input email OK')

  await page.fill('input#password', process.env.INFO_PASSWORD)
  console.log('input password OK')

  await page.click('.red.SignupButton') // 點擊送出
  console.log('btn login OK')

  await page.screenshot({ path: '/screenshot/2_login_modal.png' })

  // 已登入畫面
  await page.waitForSelector('#homefeedGridFadeInTransitionContainer')
  console.log('waitForSelector OK')

  const isLogin = await checkIsLoginByCookie(page)
  return isLogin
}

// api: get?key=keyword
export const searchOnPinterest = async (keyword: string, page?) => {
  if (!browserInstance) {
    await initPinterestService()
    await loginToPinterest()
  }
  if (!currentPage) throw new Error('Page instance not initialized')
  if (!page) page = currentPage

  //console.log('searchOnPinterest', page)
  const isLogin = await checkIsLoginByCookie(page)
  if (!isLogin) {
    console.log(Color.red + '發生錯誤, 沒有login不能輸入文字' + Color.close)
    currentPage = null
    page.close()
    return {
      ok: false,
      playwright: {
        browserInstance: browserInstance,
        page: currentPage,
        isLogin: await checkIsLoginByCookie(page),
      },
      resource_data: [],
    }
  }

  // await page.goto(`https://www.pinterest.com/search/pins/?q=${keyword}`)
  await page.fill('input[name="searchBoxInput"]', keyword) // 輸入關鍵字
  await page.keyboard.press('Enter') // 模擬按下 Enter 鍵
  console.log('searchOnPinterest', keyword)
  const payload = await getBaseSearchResource(page) // 執行搜索操作
  return payload
}

// check login
const checkIsLoginByCookie = async (page: Page): Promise<boolean> => {
  if (!page) {
    console.log(Color.red + 'checkCookie失敗 沒有page' + Color.close)
    throw new Error('沒有page資料, 沒有cookie可判斷')
  }

  const cookies = await page?.context().cookies()
  

  // 檢查是否存在名為 "_auth" 的 cookie
  const isLoggedIn = cookies.some((cookie) => {
    return cookie.name === '_auth' && cookie.value == '1'
  })

  await page.screenshot({ path: '/screenshot/checkIsLogin.png' })

  if (isLoggedIn) {
    console.log(Color.green + '登錄中' + Color.close)
    return true
  } else {
    console.log(Color.red + '未登錄' + Color.close)
    page.close()
    currentPage = null
    console.log('page close')
    return false
  }
}

// wait response
const getBaseSearchResource = async (page: Page) => {
  const response = await page.waitForResponse((response) => {
    return response.url().startsWith('https://www.pinterest.com/resource/BaseSearchResource/get/')
  })

  const imgDataList = []
  const res = await response.json()
  console.log(Color.white + '[playwright] response' + Color.close)
  console.log(res)
  if (res.resource.name == 'BaseSearchResource') {
    const results = res.resource_response.data.results
    results.map((imgData) => {
      const url = imgData.images?.orig.url
      const description = imgData?.description
      const autoAltText = imgData?.auto_alt_text
      imgDataList.push({ url, description, autoAltText })
    })
  }

  const payload = {
    ok: true,
    playwright: {
      browserInstance: browserInstance,
      page: currentPage,
      isLogin: await checkIsLoginByCookie(page),
    },
    resource_data: imgDataList,
  }
  return payload
}
