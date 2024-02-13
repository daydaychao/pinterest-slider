package main

import (
	"fmt"
	"goServer/controller"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

// 跨來源資源共用中間件
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 檢查來源是否允許
		origin := r.Header.Get("Origin")
		if origin != "" {
			// 設定允許的來源為任意的 localhost 區域
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		// 如果是 OPTIONS 請求，直接回應成功
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 傳遞給下一個處理函式
		next.ServeHTTP(w, r)
	})
}

func main() {

	// env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := ":3088"

	// 創建路由器
	mux := http.NewServeMux()

	// 添加路由與對應的處理函式
	mux.HandleFunc("/login", controller.AppLogin)
	mux.HandleFunc("/get", controller.AppSearch)

	fmt.Printf("Server is running on port %s\n", port)

	// 全域設置 CORS 中間件，並將路由器與中間件結合
	http.ListenAndServe(port, cors(mux))
}
