package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
		Data:    nil, // Go 的 nil 会被序列化为 JSON 的 null
	})
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ping", pingHandler)
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
