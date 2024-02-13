package robot

import (
	"fmt"
	"log"
	"os"

	"github.com/playwright-community/playwright-go"
)

var browserInstance playwright.Browser
var currentPage playwright.Page

func logError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
		return
	}
}

// initPinterestService 初始化 Pinterest 服务
func InitPinterestService() (playwright.Browser, error) {

	if browserInstance == nil {
		fmt.Println("[Pinterest机器人服务] 启动..")
		pw, err := playwright.Run()
		if err != nil {
			logError("无法启动 Playwright: %v", err)
		}
		browser, err := pw.Chromium.Launch()
		if err != nil {
			logError("无法启动浏览器: %v", err)
		}
		browserInstance = browser
	}
	return browserInstance, nil
}

// RunPinterestRobot 运行 Pinterest 机器人
func LoginToPinterest() (bool, error) {
	// 初始化 Pinterest 服务
	browser, err := InitPinterestService()
	if err != nil {
		logError("初始化 Pinterest 服务失败: %v", err)
	}
	fmt.Println("init ok")

	// 创建页面
	page, err := browser.NewPage()
	if err != nil {
		logError("无法创建页面: %v", err)
	}

	// 导航到 Pinterest 页面
	if _, err := page.Goto("https://www.pinterest.com/"); err != nil {
		logError("无法导航到 Pinterest 页面: %v", err)
	}
	fmt.Println("goto pinterest ok")

	// 点击登录按钮
	if err := page.Locator("[data-test-id='simple-login-button']").Click(); err != nil {
		logError("点击登录按钮时发生错误: %v", err)
	}
	fmt.Println("login btn ok")

	// 等待登录模态框
	if err := page.Locator("[data-test-id='login-modal-default']").WaitFor(); err != nil {
		logError("等待登录模态框时发生错误: %v", err)
	}
	fmt.Println("login modal ok")

	// Access environment variables
	username := os.Getenv("INFO_USERNAME")
	password := os.Getenv("INFO_PASSWORD")

	fmt.Println("Login user:", username+password)

	// Username password login
	if err := page.Locator("input#email").Fill(username); err != nil {
		logError("填写用户名时发生错误: %v", err)
	}
	fmt.Println("email ok")

	if err := page.Locator("input#password").Fill(password); err != nil {
		logError("填写密码时发生错误: %v", err)
	}
	fmt.Println("password ok")

	if err := page.Locator(".red.SignupButton").Click(); err != nil {
		logError("点击登录按钮时发生错误: %v", err)
	}
	fmt.Println("red btn login ok")

	// 截图并保存到文件
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("/screenshot/waiting_login.png"),
	}); err != nil {
		logError("截圖时发生错误: %v", err)
	}
	fmt.Println("Screenshot waiting_login")

	// 等待登录成功
	if err := page.Locator("#homefeedGridFadeInTransitionContainer").WaitFor(); err != nil {
		logError("等待登录成功时发生错误: %v", err)
	}
	fmt.Println("Timeout waiting_login")

	isLogin, err := CheckIsLoginByCookie(page)
	fmt.Println("isLogin:", isLogin)
	if err != nil {
		logError("擷取check login發生錯誤: %v", err)
	}

	return isLogin, nil
}

type Payload struct {
	Ok           bool
	Playwright   PlaywrightData
	ResourceData []ImageData
}

type PlaywrightData struct {
	BrowserInstance playwright.Browser
	Page            playwright.Page
	IsLogin         bool
}

type ImageData struct {
	URL         string
	Description string
	AutoAltText string
}

func SearchOnPinterest(keyword string) (Payload, error) {
	println("test SearchOnPinterest %s", keyword)

	isLogin, err := CheckIsLoginByCookie(currentPage)
	if err != nil {
		logError("Error checking login status: %v", err)
	}

	payload := Payload{
		Ok: true,
		Playwright: PlaywrightData{
			BrowserInstance: browserInstance,
			Page:            currentPage,
			IsLogin:         isLogin,
		},
		ResourceData: []ImageData{},
	}
	return payload, nil

}

// Cookie 模拟 Cookie 类型
type Cookie struct {
	name  string
	value string
}

// Color 模拟 Color 类型
type Color string

const (
	// 定义颜色常量
	red   Color = "red"
	green Color = "green"
	close Color = "close"
)

// checkIsLoginByCookie 模拟 checkIsLoginByCookie 函数
func CheckIsLoginByCookie(page playwright.Page) (bool, error) {
	if page == nil {
		log.Println(red + "checkCookie失敗 沒有page" + close)
		return false, fmt.Errorf("沒有page資料, 沒有cookie可判斷")
	}

	// 模拟获取 cookie
	cookies := []*Cookie{
		{"_auth", "1"},
		{"other_cookie", "value"},
	}

	// 检查是否存在名为 "_auth" 的 cookie
	var isLoggedIn bool
	for _, cookie := range cookies {
		if cookie.name == "_auth" && cookie.value == "1" {
			isLoggedIn = true
			break
		}
	}

	// 模拟截图
	log.Println("截图：/screenshot/checkIsLogin.png")

	if isLoggedIn {
		log.Println(green + "登錄中" + close)
		return true, nil
	}

	log.Println(red + "未登錄" + close)
	// 模拟关闭页面
	log.Println("page close")
	return false, nil
}
