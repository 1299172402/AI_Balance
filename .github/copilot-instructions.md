# AI Balance — 仓库指令

## 📋 项目概览

轻量级 Go HTTP 服务，提供 AI 平台 Token 余额查询接口。监听 `:8080`。

## 🛠 快速命令

| 操作 | 命令 |
|------|------|
| 本地运行 | `go run main.go` |
| 编译二进制 | `go build -o ai-balance .` |
| 构建 Docker 镜像 | `make build-docker` |
| 构建 Docker（国内源） | `make build-docker-in-china` |
| 运行 Docker 容器 | `make run-docker` |

## 📁 项目结构

```
├── main.go                         # 入口，路由 / 和 /ping
├── Dockerfile                      # 多阶段构建（golang → scratch）
├── Makefile                        # 构建/运行快捷命令
├── go.mod                          # Go 模块定义
├── .github/
│   ├── copilot-instructions.md     # ← 本文件
│   └── workflows/docker-publish.yml  # CI/CD → ghcr.io
├── .ai/
│   ├── commit.md                   # Commit 规范 + AI 身份注册表
│   └── memory/                     # 开发备忘
│       ├── go_environment.md       # Go 环境与命令
│       ├── docker_deploy.md        # Docker 本地构建
│       └── ci_cd.md                # GitHub Actions 工作流
```

## 📐 Docker 构建要点

- 多阶段构建：stage 1 用 `golang:1.26-alpine` 编译，stage 2 用 `scratch` 运行
- Go 编译为纯静态二进制（`CGO_ENABLED=0`），无需外部运行时，最终镜像 ~15MB
- 国内构建：传入 `--build-arg APK_MIRROR=mirrors.tuna.tsinghua.edu.cn`

## 🔁 CI/CD

- GitHub Actions 自动构建 → 推送到 `ghcr.io/1299172402/AI_Balance`
- 触发：push `main` / 创建 `v*` tag / PR to `main`
- 标签：`latest`、`v1.2.3`、`sha-xxxxxx`

## ✅ 提交规范

每次 `git commit` **必须**遵守以下规则（详见 `.ai/commit.md`）：

1. 格式：`<type>: <subject>` + Description（中文）
2. **必须**添加 `Co-authored-by:`，身份从 `.ai/commit.md` 已注册表格中选择
3. 首次见到的 AI **必须先注册自己**到表格中

## 📖 开发备忘

深度信息在 `.ai/memory/` 下，按主题分开。新增知识请更新对应文件。
