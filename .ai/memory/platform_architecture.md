# Platform 抽象架构

## 目录结构

```
platform/           # 通用抽象层
├── types.go        # 公共类型 + PlatformClient 接口
└── router.go       # 通用路由模板
deepseek/           # DeepSeek 具体实现
└── client.go       # 实现 PlatformClient 接口
```

## PlatformClient 接口

```go
type PlatformClient interface {
    Name() string                          // 平台名，用作 URL 前缀
    GetBalance() (*Balance, error)
    GetTokenUsage(month, year int) (*TokenUsageResp, error)
    GetCostUsage(month, year int) (*CostUsageResp, error)
}
```

## 路由注册

main.go 中调用 `platform.RegisterRoutes(mux, client)` 即可注册三个端点：
- `GET /{name}/balance`
- `GET /{name}/usage/tokens?month=&year=`
- `GET /{name}/usage/cost?month=&year=`

## 返回值结构

所有响应包在 `{ code, message, data }` 中，data 结构：
- balance: `{ balance, currency }`
- usage/tokens: `{ month, by_model: [{ model, total: {input_cache_hit,...}, days: [{date,...}] }] }`
- usage/cost: `{ month, by_model: [{ model, total: {cost}, days: [{date, cost}] }] }`

## OpenAPI 规范

- `openapi.json` — OpenAPI 3.1 规范文件，通过 `go:embed` 嵌入
- 端点 `GET /openapi.json` 返回该文件，可供 Scalar/Swagger UI 加载

## 接入新平台

1. 新建目录如 `xiaomi/`
2. 实现 `platform.PlatformClient` 接口
3. 在 main.go 中注册即可
