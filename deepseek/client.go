package deepseek

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/1299172402/AI_Balance/platform"
)

// ---------------------------------------------------------------------------
// 配置与客户端
// ---------------------------------------------------------------------------

const (
	defaultBaseURL = "https://platform.deepseek.com"
	defaultTimeout = 30 * time.Second
)

// Client DeepSeek 平台 API 客户端
//
// 使用从 platform.deepseek.com 浏览器登录获取的 Token 鉴权。
// Token 在浏览器开发者工具 → 应用 → Cookie 中可找到（.thumbcache_xxx 或 Authorization 请求头）。
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient 从环境变量创建客户端
//
// 环境变量:
//   - DEEPSEEK_PLATFORM_TOKEN  : 浏览器登录后的 Token（必填）
//   - DEEPSEEK_BASE_URL        : API 基础地址（可选，默认 https://platform.deepseek.com）
func NewClient() *Client {
	return &Client{
		BaseURL:    getEnv("DEEPSEEK_BASE_URL", defaultBaseURL),
		Token:      os.Getenv("DEEPSEEK_PLATFORM_TOKEN"),
		HTTPClient: &http.Client{Timeout: defaultTimeout},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ---------------------------------------------------------------------------
// 底层 HTTP 请求
// ---------------------------------------------------------------------------

// doRequest 发起 HTTP 请求，通过 biz_data 层并返回其原始 JSON
func (c *Client) doRequest(method, path string) ([]byte, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("DEEPSEEK_PLATFORM_TOKEN 未设置")
	}

	req, err := http.NewRequest(method, c.BaseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 浏览器同源请求所需 header（绕过 WAF）
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://platform.deepseek.com/usage")
	req.Header.Set("X-App-Version", "20240425.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API HTTP %d: %s", resp.StatusCode, string(body))
	}

	// 统一解包: 所有 DeepSeek 平台 API 返回 { code, msg, data: { biz_code, biz_msg, biz_data } }
	var wrapper struct {
		Code int `json:"code"`
		Data struct {
			BizCode int             `json:"biz_code"`
			BizData json.RawMessage `json:"biz_data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("解析响应包装层失败: %w", err)
	}
	if wrapper.Code != 0 {
		return nil, fmt.Errorf("API 错误 code=%d", wrapper.Code)
	}
	if wrapper.Data.BizCode != 0 {
		return nil, fmt.Errorf("业务错误 biz_code=%d", wrapper.Data.BizCode)
	}

	return wrapper.Data.BizData, nil
}

// ---------------------------------------------------------------------------
// 响应结构体 — 用户 / 汇总
// ---------------------------------------------------------------------------

// UserSummary 用户汇总（来自 /api/v0/users/get_user_summary）
type UserSummary struct {
	CurrentToken                  int64            `json:"current_token"`
	MonthlyUsage                  string           `json:"monthly_usage"`
	TotalUsage                    int64            `json:"total_usage"`
	NormalWallets                 []Wallet         `json:"normal_wallets"`
	BonusWallets                  []Wallet         `json:"bonus_wallets"`
	TotalAvailableTokenEstimation string           `json:"total_available_token_estimation"`
	MonthlyCosts                  []CurrencyAmount `json:"monthly_costs"`
	MonthlyTokenUsage             string           `json:"monthly_token_usage"`
}

// Wallet 钱包
type Wallet struct {
	Currency        string `json:"currency"`
	Balance         string `json:"balance"`
	TokenEstimation string `json:"token_estimation"`
}

// CurrencyAmount 带币种的金额
type CurrencyAmount struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// ---------------------------------------------------------------------------
// 响应结构体 — Usage（用量 / 金额：amount 和 cost 共用同一结构）
// ---------------------------------------------------------------------------

// UsageReport 用量报告（total + 每日明细）
type UsageReport struct {
	Total []ModelUsage `json:"total"`
	Days  []DayUsage   `json:"days"`
}

// ModelUsage 单个模型的用量汇总
type ModelUsage struct {
	Model string      `json:"model"`
	Usage []UsageItem `json:"usage"`
}

// UsageItem 一个用量指标
type UsageItem struct {
	Type   string `json:"type"`
	Amount string `json:"amount"`
}

// DayUsage 某一天的用量
type DayUsage struct {
	Date string       `json:"date"`
	Data []ModelUsage `json:"data"`
}

// ---------------------------------------------------------------------------
// API 方法 — 用户信息
// ---------------------------------------------------------------------------

// GetUserSummary 获取用户汇总（钱包余额、本月用量、预估可用的 Token 数等）
//
//	GET /api/v0/users/get_user_summary
func (c *Client) GetUserSummary() (*UserSummary, error) {
	raw, err := c.doRequest("GET", "/api/v0/users/get_user_summary")
	if err != nil {
		return nil, err
	}
	var s UserSummary
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, fmt.Errorf("解析用户汇总失败: %w", err)
	}
	return &s, nil
}

// ---------------------------------------------------------------------------
// API 方法 — 用量（Token 数）
// ---------------------------------------------------------------------------

// GetUsageAmount 获取指定月份的 Token 用量（按模型、按天）
//
//	GET /api/v0/usage/amount?month=5&year=2026
func (c *Client) GetUsageAmount(month, year int) (*UsageReport, error) {
	path := fmt.Sprintf("/api/v0/usage/amount?month=%d&year=%d", month, year)
	raw, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	var r UsageReport
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("解析用量数据失败: %w", err)
	}
	return &r, nil
}

// ---------------------------------------------------------------------------
// API 方法 — 消费金额
// ---------------------------------------------------------------------------

// GetUsageCost 获取指定月份的消费金额（按模型、按天）
//
//	GET /api/v0/usage/cost?month=5&year=2026
//
// 注意: cost 接口的 biz_data 是数组 [{total, days}]，取第一个元素
func (c *Client) GetUsageCost(month, year int) (*UsageReport, error) {
	path := fmt.Sprintf("/api/v0/usage/cost?month=%d&year=%d", month, year)
	raw, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	var arr []UsageReport
	if err := json.Unmarshal(raw, &arr); err != nil {
		return nil, fmt.Errorf("解析消费数据失败: %w", err)
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("消费数据为空")
	}
	return &arr[0], nil
}

// ---------------------------------------------------------------------------
// PlatformClient 接口实现
// ---------------------------------------------------------------------------

// Name 返回平台名称，用作 URL 路径前缀
func (c *Client) Name() string {
	return "deepseek"
}

// GetBalance 获取余额
func (c *Client) GetBalance() (*platform.Balance, error) {
	s, err := c.GetUserSummary()
	if err != nil {
		return nil, err
	}

	total := 0.0
	currency := "CNY"
	for _, w := range s.NormalWallets {
		if b, err := strconv.ParseFloat(w.Balance, 64); err == nil {
			total += b
		}
		if w.Currency != "" {
			currency = w.Currency
		}
	}
	for _, w := range s.BonusWallets {
		if b, err := strconv.ParseFloat(w.Balance, 64); err == nil {
			total += b
		}
	}

	return &platform.Balance{
		Balance:  total,
		Currency: currency,
	}, nil
}

// GetTokenUsage 获取指定月份的 Token 用量（汇总 + 每日明细）
func (c *Client) GetTokenUsage(month, year int) (*platform.TokenUsageResp, error) {
	r, err := c.GetUsageAmount(month, year)
	if err != nil {
		return nil, err
	}

	resp := &platform.TokenUsageResp{
		Month:   fmt.Sprintf("%d-%02d", year, month),
		ByModel: make([]platform.ModelToken, 0),
	}

	modelIdx := make(map[string]int)

	// 汇总（total）
	for _, m := range r.Total {
		mt := platform.ModelToken{
			Model: m.Model,
			Total: platform.TokenTotal{},
			Days:  make([]platform.DayToken, 0),
		}
		for _, u := range m.Usage {
			amt, _ := strconv.ParseInt(u.Amount, 10, 64)
			switch u.Type {
			case "PROMPT_CACHE_HIT_TOKEN":
				mt.Total.InputCacheHit = amt
			case "PROMPT_CACHE_MISS_TOKEN":
				mt.Total.InputCacheMiss = amt
			case "RESPONSE_TOKEN":
				mt.Total.Output = amt
			}
		}
		modelIdx[m.Model] = len(resp.ByModel)
		resp.ByModel = append(resp.ByModel, mt)
	}

	// 每日明细（days）
	for _, d := range r.Days {
		for _, m := range d.Data {
			idx, ok := modelIdx[m.Model]
			if !ok {
				continue
			}
			dt := platform.DayToken{Date: d.Date}
			for _, u := range m.Usage {
				amt, _ := strconv.ParseInt(u.Amount, 10, 64)
				switch u.Type {
				case "PROMPT_CACHE_HIT_TOKEN":
					dt.InputCacheHit = amt
				case "PROMPT_CACHE_MISS_TOKEN":
					dt.InputCacheMiss = amt
				case "RESPONSE_TOKEN":
					dt.Output = amt
				}
			}
			resp.ByModel[idx].Days = append(resp.ByModel[idx].Days, dt)
		}
	}

	return resp, nil
}

// GetCostUsage 获取指定月份的费用（汇总 + 每日明细）
func (c *Client) GetCostUsage(month, year int) (*platform.CostUsageResp, error) {
	r, err := c.GetUsageCost(month, year)
	if err != nil {
		return nil, err
	}

	resp := &platform.CostUsageResp{
		Month:   fmt.Sprintf("%d-%02d", year, month),
		ByModel: make([]platform.ModelCost, 0),
	}

	modelIdx := make(map[string]int)

	// 汇总（total）
	for _, m := range r.Total {
		mc := platform.ModelCost{
			Model: m.Model,
			Total: platform.CostTotal{},
			Days:  make([]platform.DayCost, 0),
		}
		for _, u := range m.Usage {
			amt, _ := strconv.ParseFloat(u.Amount, 64)
			mc.Total.Cost += amt
		}
		modelIdx[m.Model] = len(resp.ByModel)
		resp.ByModel = append(resp.ByModel, mc)
	}

	// 每日明细（days）
	for _, d := range r.Days {
		for _, m := range d.Data {
			idx, ok := modelIdx[m.Model]
			if !ok {
				continue
			}
			dc := platform.DayCost{Date: d.Date}
			for _, u := range m.Usage {
				amt, _ := strconv.ParseFloat(u.Amount, 64)
				dc.Cost += amt
			}
			resp.ByModel[idx].Days = append(resp.ByModel[idx].Days, dc)
		}
	}

	return resp, nil
}
