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

func AppLogin(w http.ResponseWriter, r *http.Request) {
	// 初始化 Pinterest 服務
	_, err := robot.InitPinterestService()
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

	// 構建回應
	response := map[string]interface{}{
		"ok": payload.Ok,
		"code": func() int {
			if payload.Ok {
				return 200
			}
			return 99
		}(),
		"data": payload,
		"message": func() string {
			if payload.Ok {
				return "Got data"
			}
			return "No data"
		}(),
	}

	// 將回應轉換為 JSON 格式並發送
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
