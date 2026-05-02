# Docker & 部署

> 构建、运行、CI/CD 相关备忘。

---

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make build-docker` | 构建镜像（默认源） |
| `make build-docker-in-china` | 构建镜像（清华 apk 源） |
| `make run-docker` | 运行容器（前台，`Ctrl+C` 退出） |
| `make run-native` | 本地直接 `go run main.go` |

## 镜像拉取

| 源 | 地址 |
|------|------|
| GitHub Container Registry | `ghcr.io/1299172402/ai_balance:latest` |
| 南京大学镜像（国内加速） | `ghcr.nju.edu.cn/1299172402/ai_balance:latest` |

运行容器：

```bash
docker run --rm -p 8080:8080 ghcr.io/1299172402/ai_balance
```

国内用户：

```bash
docker run --rm -p 8080:8080 ghcr.nju.edu.cn/1299172402/ai_balance
```

## Dockerfile 要点

- **多阶段构建**：Stage 1 用 `golang:1.26-alpine` 编译，Stage 2 用 `scratch` 运行
- **scratch 可行原因**：Go 编译为静态二进制，运行时、GC、net/http 全在二进制内，不依赖外部运行时
- **SSL 证书**：从 builder 阶段复制 `/etc/ssl/certs/ca-certificates.crt`，供 `net/http` 发 HTTPS 请求
- **国内镜像**：`--build-arg APK_MIRROR=mirrors.tuna.tsinghua.edu.cn` 替换 apk 源
