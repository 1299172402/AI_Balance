package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/1299172402/AI_Balance/deepseek"
	"github.com/1299172402/AI_Balance/platform"
)

// 默认为 nil，由 openapi_dev.go 在 dev 编译时赋值
var openAPIHandler func(w http.ResponseWriter, r *http.Request)

// Response 通用响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Code:    404,
		Message: "not found, try /ping",
		Data:    nil,
	})
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Code:    0,
		Message: "pong",
		Data:    nil,
	})
}

func main() {
	mux := http.NewServeMux()

	// 注册基础路由
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("GET /ping", pingHandler)

	// 仅 dev 编译时注册 /openapi.json
	if openAPIHandler != nil {
		mux.HandleFunc("GET /openapi.json", openAPIHandler)
	}

	// 注册平台子路由（deepseek）
	deepseekClient := deepseek.NewClient()
	platform.RegisterRoutes(mux, deepseekClient)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
