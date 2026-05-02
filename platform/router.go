package platform

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// writeJSON 成功响应
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// writeError 错误响应
func writeError(w http.ResponseWriter, httpStatus int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(Response{
		Code:    -1,
		Message: msg,
	})
}

// parseMonthYear 从查询参数解析 month 和 year，缺省则用当前月份
func parseMonthYear(r *http.Request) (month, year int) {
	now := time.Now()
	year = now.Year()
	month = int(now.Month())

	if m := r.URL.Query().Get("month"); m != "" {
		if v, err := strconv.Atoi(m); err == nil && v >= 1 && v <= 12 {
			month = v
		}
	}
	if y := r.URL.Query().Get("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil && v >= 2020 {
			year = v
		}
	}
	return
}

// RegisterRoutes 注册平台通用路由
//
//	GET /{name}/balance
//	GET /{name}/usage/tokens?month=&year=
//	GET /{name}/usage/cost?month=&year=
func RegisterRoutes(mux *http.ServeMux, client PlatformClient) {
	prefix := client.Name()

	// GET /{platform}/balance
	mux.HandleFunc("GET /"+prefix+"/balance", func(w http.ResponseWriter, r *http.Request) {
		data, err := client.GetBalance()
		if err != nil {
			log.Printf("[%s] 查询余额失败: %v", prefix, err)
			writeError(w, http.StatusInternalServerError, "查询余额失败: "+err.Error())
			return
		}
		writeJSON(w, data)
	})

	// GET /{platform}/usage/tokens?month=&year=
	mux.HandleFunc("GET /"+prefix+"/usage/tokens", func(w http.ResponseWriter, r *http.Request) {
		m, y := parseMonthYear(r)
		data, err := client.GetTokenUsage(m, y)
		if err != nil {
			log.Printf("[%s] 查询 Token 用量失败: %v", prefix, err)
			writeError(w, http.StatusInternalServerError, "查询 Token 用量失败: "+err.Error())
			return
		}
		writeJSON(w, data)
	})

	// GET /{platform}/usage/cost?month=&year=
	mux.HandleFunc("GET /"+prefix+"/usage/cost", func(w http.ResponseWriter, r *http.Request) {
		m, y := parseMonthYear(r)
		data, err := client.GetCostUsage(m, y)
		if err != nil {
			log.Printf("[%s] 查询费用失败: %v", prefix, err)
			writeError(w, http.StatusInternalServerError, "查询费用失败: "+err.Error())
			return
		}
		writeJSON(w, data)
	})
}
