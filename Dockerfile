# ============================= 构建阶段 =============================
# 使用官方Golang镜像作为构建环境
FROM golang:1.23rc1-alpine AS builder

# 设置工作目录
WORKDIR /app

# 首先复制依赖管理文件（利用Docker层缓存优化）
COPY go.mod go.sum ./

# 下载依赖（-mod=readonly确保不会修改go.mod）
RUN go mod download

# 复制所有源代码到容器
COPY . .

# 构建可执行文件（禁用CGO以确保静态编译）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/chat-app ./cmd/server/main.go

# ============================= 运行时阶段 =============================
# 使用更小的Alpine镜像作为运行时基础
FROM alpine:3.17

# 安装CA证书（用于需要HTTPS请求的场景）
RUN apk --no-cache add ca-certificates tzdata && \
    # 创建专用用户和用户组
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup -h /app && \
    # 设置时区为上海（可按需修改）
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 切换到应用目录
WORKDIR /app

RUN mkdir -p /app/configs && chown appuser:appgroup /app/configs

# 从构建阶段拷贝可执行文件
COPY --from=builder --chown=appuser:appgroup /app/chat-app .

# 从构建阶段拷贝配置文件（如有需要）
COPY --from=builder --chown=appuser:appgroup /app/configs/config.yaml ./configs/

# 声明使用的端口
EXPOSE 8080

# 切换到非root用户
USER appuser

# 启动应用程序（参数可通过环境变量覆盖）
CMD ["./chat-app"]