# ---- Stage 1: 编译 ----
FROM golang:1.26-alpine AS builder

# 中国大陆用户：构建时传 --build-arg APK_MIRROR=mirrors.tuna.tsinghua.edu.cn
ARG APK_MIRROR=dl-cdn.alpinelinux.org

# 如果传了国内镜像源，替换 apk 源
RUN if [ "$APK_MIRROR" != "dl-cdn.alpinelinux.org" ]; then \
        sed -i "s|dl-cdn.alpinelinux.org|${APK_MIRROR}|g" /etc/apk/repositories; \
    fi

RUN apk add --no-cache ca-certificates
WORKDIR /src

COPY go.mod* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/ai-balance ./cmd/ai-balance

# ---- Stage 2: 运行 ----
# scratch 是 Docker 最精简的空镜像，0 文件
FROM scratch

# SSL 根证书 —— Go 用 net/http 发 HTTPS 请求时需要它
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# 只带编译好的二进制
COPY --from=builder /app/ai-balance /ai-balance

EXPOSE 8080
ENTRYPOINT ["/ai-balance"]
