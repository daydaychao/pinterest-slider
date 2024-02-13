package robot

import (
	"encoding/base64"
	"fmt"
	"goServer/logColor"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
)

type Payload struct {
	Ok           bool           `json:"ok"`
	Playwright   PlaywrightData `json:"playwright"`
	ResourceData []ImageData    `json:"resource_data"`
}

type PlaywrightData struct {
	BrowserInstance playwright.Browser `json:"browserInstance"`
	Page            playwright.Page    `json:"page"`
	IsLogin         bool               `json:"isLogin"`
}

type ImageData struct {
	URL         string `json:"url"`
	Description string `json:"description"`
	AutoAltText string `json:"autoAltText"`
}

var browserInstance playwright.Browser
var currentPage playwright.Page

func logError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
		return
	}
}

// playwright 初始化  browser
func InitService() (playwright.Browser, error) {

	if browserInstance == nil {
		fmt.Println(logColor.Green("[Pinterest机器人服务] 启动.."))
		pw, err := playwright.Run()
		if err != nil {
			logError("无法启动 playwright: %v", err)
		}
		browser, err := pw.Chromium.Launch()
		if err != nil {
			logError("无法启动浏览器: %v", err)
		}
		browserInstance = browser

		initPage()
	}
	return browserInstance, nil
}

// playwright 初始化  page
func initPage() {
	// 创建页面
	page, err := browserInstance.NewPage()
	if err != nil {
		logError("无法创建页面: %v", err)
	}
	currentPage = page
	page.SetViewportSize(800, 600)
}

// playwright 截圖
func TakeScreenshot() (bool, string) {
	fmt.Println(logColor.Green("[Pinterest机器人服务] 截圖"))

	// 创建通道用于接收截图路径和结果
	screenChannel := make(chan struct {
		Path    string
		Success bool
	})

	// 检查 currentPage 是否为nil
	if currentPage != nil {
		// 进行屏幕截图
		go func() {
			// 执行截图操作，并将结果发送到通道中
			_, err := currentPage.Screenshot(playwright.PageScreenshotOptions{
				FullPage: playwright.Bool(false),
				Path:     playwright.String("./screen/api_screen.png"),
			})
			if err != nil {
				// 如果截图失败，则打印错误信息
				log.Println("截图失败")
				screenChannel <- struct {
					Path    string
					Success bool
				}{Path: "", Success: false}
				return
			}

			// 将截图转换为 base64 格式
			data, err := os.ReadFile("./screen/api_screen.png")
			if err != nil {
				log.Println("读取截图文件失败")
				screenChannel <- struct {
					Path    string
					Success bool
				}{Path: "", Success: false}
				return
			}
			base64Image := base64.StdEncoding.EncodeToString(data)

			// 将 base64 编码的图片数据和成功标志发送到通道中
			screenChannel <- struct {
				Path    string
				Success bool
			}{Path: base64Image, Success: true}
		}()
	}

	// 接收截图路径和结果
	result := <-screenChannel
	if result.Success {
		log.Println("截图成功")
		return true, result.Path
	} else {
		log.Println("截图失败")
		return false, ""
	}
}

// Pinterest 机器人 login
func LoginToPinterest() (bool, error) {

	if browserInstance == nil {
		InitService()
	}

	if currentPage == nil {
		initPage()
	}

	// 导航到 Pinterest 页面
	if _, err := currentPage.Goto("https://www.pinterest.com/"); err != nil {
		logError("无法导航到 Pinterest 页面: %v", err)
	}
	fmt.Println("goto pinterest ok")

	// 点击登录按钮
	if err := currentPage.Locator("[data-test-id='simple-login-button']").Click(); err != nil {
		logError("点击登录按钮时发生错误: %v", err)
	}
	fmt.Println("login btn ok")

	// 等待登录模态框
	if err := currentPage.Locator("[data-test-id='login-modal-default']").WaitFor(); err != nil {
		logError("等待登录模态框时发生错误: %v", err)
	}
	fmt.Println("login modal ok")

	// Access environment variables
	username := os.Getenv("INFO_USERNAME")
	password := os.Getenv("INFO_PASSWORD")

	fmt.Println("Login user:", username+password)

	// Username password login
	if err := currentPage.Locator("input#email").Fill(username); err != nil {
		logError("填写用户名时发生错误: %v", err)
	}
	fmt.Println("email ok")

	if err := currentPage.Locator("input#password").Fill(password); err != nil {
		logError("填写密码时发生错误: %v", err)
	}
	fmt.Println("password ok")

	if err := currentPage.Locator(".red.SignupButton").Click(); err != nil {
		logError("点击登录按钮时发生错误: %v", err)
	}
	fmt.Println("red btn login ok")

	// 截图并保存到文件
	// if _, err = page.Screenshot(playwright.PageScreenshotOptions{
	// 	Path: playwright.String("/screenshot/waiting_login.png"),
	// }); err != nil {
	// 	logError("截圖时发生错误: %v", err)
	// }
	//fmt.Println("Screenshot waiting_login")

	// 等待登录成功
	if err := currentPage.Locator("#homefeedGridFadeInTransitionContainer").WaitFor(); err != nil {
		logError("等待登录成功时发生错误: %v", err)
	}
	fmt.Println("Timeout waiting_login")

	isLogin, err := CheckIsLoginByCookie(currentPage)
	fmt.Println("isLogin:", isLogin)
	if err != nil {
		logError("擷取check login發生錯誤: %v", err)
	}

	return isLogin, nil
}

// Pinterest 机器人 查字
func SearchOnPinterest(keyword string) (Payload, error) {
	println("SearchOnPinterest %c", keyword)

	if browserInstance == nil {
		InitService()
		LoginToPinterest()
	}
	if currentPage == nil {
		logError("沒有currentPage: %v", nil)
	}

	isLogin, err := CheckIsLoginByCookie(currentPage)
	if err != nil {
		logError("Error checking login status: %v", err)
	}

	// 未登入
	if !isLogin {
		log.Println(logColor.Red("沒有login不能輸入文字"))

		payload := Payload{
			Ok: false,
			Playwright: PlaywrightData{
				BrowserInstance: browserInstance,
				Page:            currentPage,
				IsLogin:         isLogin,
			},
			ResourceData: []ImageData{},
		}
		return payload, nil
	}

	// 填充搜索框
	if err := currentPage.Locator("input[name=\"searchBoxInput\"]").Fill(keyword); err != nil {
		return Payload{}, fmt.Errorf("填充搜索框时发生错误: %v", err)
	}
	log.Println("input ok")

	// 模拟按下 Enter 键
	if err := currentPage.Keyboard().Press("Enter"); err != nil {
		return Payload{}, fmt.Errorf("模拟按下 Enter 键时发生错误: %v", err)
	}
	log.Println("keyword enter")

	// 执行搜索操作并返回结果
	payload, err := getBaseSearchResource(currentPage)
	if err != nil {
		return Payload{}, fmt.Errorf("执行搜索操作时发生错误: %v", err)
	}

	log.Println("payload", payload)

	// 截图并保存到文件
	if _, err = currentPage.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("/screenshot/input_keyword.png"),
	}); err != nil {
		logError("截圖时发生错误: %v", err)
	}
	fmt.Println("Screenshot input_keyword")

	// 已登入
	return payload, nil

}

// 檢查有沒有login
func CheckIsLoginByCookie(page playwright.Page) (bool, error) {
	if page == nil {
		log.Println(logColor.Red("checkCookie失敗 沒有page"))
		return false, nil
	}

	// 获取页面的cookies
	cookies, err := page.Context().Cookies()
	if err != nil {
		log.Println(logColor.Red("无法获取页面的cookies: %v", err))
	}

	// Convert cookies to []*playwright.Cookie
	var cookiePointers []*playwright.Cookie
	for _, cookie := range cookies {
		cookiePointers = append(cookiePointers, &cookie)
	}

	isLoggedIn := checkCookieHasAuth(cookiePointers)

	if isLoggedIn {
		log.Println(logColor.Green("登錄中"))
		return true, nil
	}

	// 模拟截图
	log.Println("截图：/screenshot/checkIsLogin.png")

	log.Println(logColor.Red("未登錄"))
	return false, nil

}

// tool 检查是否存在名为 "_auth" 的 cookie
func checkCookieHasAuth(cookies []*playwright.Cookie) bool {
	for _, cookie := range cookies {
		if cookie.Name == "_auth" && cookie.Value == "1" {
			return true
		}
	}
	return false
}

// tool 查字 res
func getBaseSearchResource(page playwright.Page) (Payload, error) {
	log.Println("getBaseSearchResource")

	// Channel response
	responseChan := make(chan playwright.Response)

	// On response
	page.OnResponse(func(response playwright.Response) {
		if strings.HasPrefix(response.URL(), "https://www.pinterest.com/resource/BaseSearchResource/get/") {
			responseChan <- response
		}
	})

	// Get response
	response := <-responseChan
	if response == nil {
		log.Println(logColor.Red("response为空"))
		// return Payload{}, errors.New("response为空")
	}
	log.Println("get response DONE")

	// 解析 JSON 数据

	var res map[string]interface{}
	if response != nil {
		if err := response.JSON(&res); err != nil {
			log.Println(logColor.Red("JSON傳換發生錯誤"))
			// return Payload{}, fmt.Errorf("解析 JSON 数据时发生错误: %v", err)
		}
		log.Println("get res parseJson DONE")

		// 提取图片数据
		imgDataList := make([]ImageData, 0)
		results := res["resource_response"].(map[string]interface{})["data"].(map[string]interface{})["results"].([]interface{})
		for _, result := range results {
			resultData := result.(map[string]interface{})
			imgURL := resultData["images"].(map[string]interface{})["orig"].(map[string]interface{})["url"].(string)
			description := resultData["description"].(string)
			autoAltText := resultData["auto_alt_text"].(string)
			imgDataList = append(imgDataList, ImageData{URL: imgURL, Description: description, AutoAltText: autoAltText})
		}
		log.Println("imgDataList", imgDataList)

		// 构造返回值
		isLogin, _ := CheckIsLoginByCookie(page)
		payload := Payload{
			Ok: true,
			Playwright: PlaywrightData{
				BrowserInstance: browserInstance,
				Page:            currentPage,
				IsLogin:         isLogin,
			},
			ResourceData: imgDataList,
		}
		return payload, nil
	}

	// 构造返回值
	isLogin, _ := CheckIsLoginByCookie(page)
	payload := Payload{
		Ok: true,
		Playwright: PlaywrightData{
			BrowserInstance: browserInstance,
			Page:            currentPage,
			IsLogin:         isLogin,
		},
		ResourceData: make([]ImageData, 0),
	}
	return payload, nil
}
