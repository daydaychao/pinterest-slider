package controller

import (
	"encoding/json"
	"fmt"
	"goServer/robot"
	"net/http"
)

func SendError(w http.ResponseWriter, errMsg string) {
	http.Error(w, fmt.Sprintf("%s", errMsg), http.StatusInternalServerError)
}

func AppInit(w http.ResponseWriter, r *http.Request) {
	// 初始化 Pinterest 服務
	isInit, err := robot.InitService()
	if err != nil {
		SendError(w, "Error initialing service:"+err.Error())
		return
	}

	// 構建回應
	var code int
	var msg string
	if isInit != nil {
		code = 200
		msg = "Init OK"
	} else {
		code = 99
		msg = "Init error"
	}
	response := map[string]interface{}{
		"ok":      isInit,
		"code":    code,
		"data":    isInit,
		"message": msg,
	}

	// 將回應轉換為 JSON 格式並發送
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AppLogin(w http.ResponseWriter, r *http.Request) {
	// 初始化 Pinterest 服務
	_, err := robot.InitService()
	if err != nil {
		SendError(w, "Error initialing service:"+err.Error())
		return
	}

	// 登入 Pinterest
	isLogin, err := robot.LoginToPinterest()
	if err != nil {
		SendError(w, "Error logging into Pinterest:"+err.Error())
		return
	}

	// 構建回應
	response := map[string]interface{}{
		"ok": isLogin,
		"code": func() int {
			if isLogin {
				return 200
			}
			return 99
		}(),
		"data": isLogin,
		"message": func() string {
			if isLogin {
				return "Login success"
			}
			return "Login failed"
		}(),
	}

	// 將回應轉換為 JSON 格式並發送
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AppSearch(w http.ResponseWriter, r *http.Request) {
	// 從查詢參數中獲取關鍵字
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "Url Param 'key' is missing", http.StatusBadRequest)
		return
	}
	key := keys[0]

	// 在 Pinterest 上搜索
	payload, err := robot.SearchOnPinterest(key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching on Pinterest: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Println("got payload", payload)

	var code int
	var msg string
	if payload.Ok {
		code = 200
		msg = "Got data"
	} else {
		code = 99
		msg = "No data"
	}

	// 構建回應
	response := map[string]interface{}{
		"ok":      payload.Ok,
		"code":    code,
		"data":    payload,
		"message": msg,
	}
	// 將回應轉換為 JSON 格式並發送
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AppScreenshot(w http.ResponseWriter, r *http.Request) {

	// 截圖
	payload, screenshot := robot.TakeScreenshot()
	if screenshot == "" {
		http.Error(w, "Error screenshot on Pinterest", http.StatusInternalServerError)
		return
	}

	var code int
	var msg string
	if payload {
		code = 200
		msg = "Got data"
	} else {
		code = 99
		msg = "No data"
	}

	// 構建回應
	response := map[string]interface{}{
		"ok":      payload,
		"code":    code,
		"data":    map[string]interface{}{"screenshot": screenshot},
		"message": msg,
	}
	// 將回應轉換為 JSON 格式並發送
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
