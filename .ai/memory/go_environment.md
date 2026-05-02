# Go Environment

---

## 模块初始化

```bash
go mod init github.com/1299172402/AI_Balance
```

- 创建 `go.mod`，指定模块路径为 `github.com/1299172402/AI_Balance`
- **本项目是 HTTP 服务（`package main`），非库模块**，无需被外部 Go 项目 import

## 依赖管理

```bash
# 下载/整理依赖，更新 go.mod & go.sum
go mod tidy
```

