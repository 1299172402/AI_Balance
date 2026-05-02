# CI/CD — GitHub Actions

> 云端自动构建与发布工作流，与本地开发环境分离。

---

## 工作流文件

- `.github/workflows/docker-publish.yml`

## 触发条件

| 事件 | 行为 |
|------|------|
| 推送 `main` 分支 | 构建并推送到 `ghcr.io` |
| 创建 `v*` 标签 | 构建并推送，额外生成 semver 标签 |
| PR → `main` | 仅构建验证，不推送 |

## 产物

- **Registry**: ghcr.io（GitHub Container Registry）
- **镜像**: `ghcr.io/1299172402/AI_Balance`
- **自动标签**: `latest`（main）、`v1.2.3`（semver）、`sha-xxxxxx`

## 认证

- 使用 `secrets.GITHUB_TOKEN` 自动登录，无需手动配置密钥
- 仓库默认该 token 已有 `packages: write` 权限

## 拉取已构建的镜像

```bash
docker pull ghcr.io/1299172402/AI_Balance:latest
```

> 首次拉取需先登录：
> `echo $GITHUB_TOKEN | docker login ghcr.io -u 1299172402 --password-stdin`