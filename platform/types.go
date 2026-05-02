package platform

// Balance 余额
type Balance struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

// TokenTotal Token 汇总
type TokenTotal struct {
	InputCacheHit  int64 `json:"input_cache_hit"`
	InputCacheMiss int64 `json:"input_cache_miss"`
	Output         int64 `json:"output"`
}

// DayToken 某天的 Token 用量
type DayToken struct {
	Date           string `json:"date"`
	InputCacheHit  int64  `json:"input_cache_hit"`
	InputCacheMiss int64  `json:"input_cache_miss"`
	Output         int64  `json:"output"`
}

// ModelToken 单个模型的 Token 消耗（汇总 + 每日明细）
type ModelToken struct {
	Model string     `json:"model"`
	Total TokenTotal `json:"total"`
	Days  []DayToken `json:"days"`
}

// TokenUsageResp Token 用量响应 data 部分
type TokenUsageResp struct {
	Month   string       `json:"month"`
	ByModel []ModelToken `json:"by_model"`
}

// CostTotal 费用汇总
type CostTotal struct {
	Cost float64 `json:"cost"`
}

// DayCost 某天的费用
type DayCost struct {
	Date string  `json:"date"`
	Cost float64 `json:"cost"`
}

// ModelCost 单个模型的费用（汇总 + 每日明细）
type ModelCost struct {
	Model string    `json:"model"`
	Total CostTotal `json:"total"`
	Days  []DayCost `json:"days"`
}

// CostUsageResp 费用响应 data 部分
type CostUsageResp struct {
	Month   string      `json:"month"`
	ByModel []ModelCost `json:"by_model"`
}

// PlatformClient 各平台需实现的接口
type PlatformClient interface {
	// Name 返回平台名称，用作 URL 路径前缀
	Name() string
	// GetBalance 获取余额
	GetBalance() (*Balance, error)
	// GetTokenUsage 获取指定月份的 Token 用量（汇总 + 每日明细）
	GetTokenUsage(month, year int) (*TokenUsageResp, error)
	// GetCostUsage 获取指定月份的费用（汇总 + 每日明细）
	GetCostUsage(month, year int) (*CostUsageResp, error)
}
