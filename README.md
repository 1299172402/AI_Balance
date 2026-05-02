# AI Balance

> 🤖 AI 贡献者：开始工作前请先查阅 `.ai/` 目录下的规范与备忘。一般情况下 .github/copilot-instructions.md 的提示词会自动在 vscode 中注入上下文。

一个轻量级的 Go HTTP 服务项目，提供 AI platform 的 token 余额查询接口。

---

## 📌 已经接入的 AI 平台

- DeepSeek
- XiaoMi MiMo

## 🚀 快速开始

### 本地运行

```bash
git clone https://github.com/1299172402/AI_Balance.git
cd AI_Balance

make run-native
# 或: go run main.go
```

服务默认在 `http://localhost:8080` 启动。

### Docker 运行

从 GitHub Container Registry 拉取并运行：

```bash
# 拉取镜像
docker pull ghcr.io/1299172402/ai_balance:latest

# 运行容器（前台，Ctrl+C 退出）
docker run --rm -p 8080:8080 --env-file .env ghcr.io/1299172402/ai_balance
```

> 国内用户可使用南京大学镜像加速：
> ```bash
> docker pull ghcr.nju.edu.cn/1299172402/ai_balance:latest
> docker run --rm -p 8080:8080 --env-file .env ghcr.nju.edu.cn/1299172402/ai_balance
> ```

也可通过 Makefile 一键构建并运行：

```bash
make build-docker       # 构建镜像（默认源）
make run-docker         # 运行容器
```
