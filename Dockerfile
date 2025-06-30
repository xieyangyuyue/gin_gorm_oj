# 阶段1：构建静态链接的二进制文件
FROM golang:1.24-alpine AS builder

# 设置国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 安装时区和证书
RUN apk add --no-cache tzdata ca-certificates

ENV GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# 优先复制依赖文件（优化缓存层）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN go build -ldflags="-w -s" -trimpath -o /app/myapp

# 阶段2：创建最小化运行环境
FROM scratch
WORKDIR /app
# 复制时区文件
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# 复制CA证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# 复制二进制文件
COPY --from=builder /app/myapp /app/myapp

# 复制配置文件
COPY --from=builder /build/config /app/config

# 设置环境变量
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["/app/myapp", "-config", "/app/config/config.ini"]

# FROM busybox:1.36 AS runtime
# WORKDIR /app
#
# # 复制所有文件
# COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /app/myapp /app/myapp
# COPY --from=builder /build/config /app/config
#
# # 设置非root用户
# RUN adduser -D -u 1000 appuser && \
#     chown -R appuser:appuser /app
# USER appuser
#
# EXPOSE 8080
# CMD ["/app/myapp", "-config", "/app/config/config.ini"]


# # 阶段2：创建最小化运行环境
# FROM alpine:3.16
# WORKDIR /app
#
#
# # 重要：复制二进制文件时添加执行权限
# COPY --from=builder /app/myapp /app/myapp
# RUN chmod +x /app/myapp  # 确保执行权限
#
# # 关键：复制配置文件目录（保持原结构）
# # 确保复制整个config目录及其内容
# COPY --from=builder /build/config /app/config/
#
# # 验证文件结构（调试时使用，正式构建可注释掉）
# # RUN ls -lR /app
#
# # 创建非root用户并赋予配置文件权限
# RUN addgroup -g 1000 appuser && \
#     adduser -u 1000 -G appuser -D -h /app appuser && \
#     chown -R appuser:appuser /app
#
# USER appuser
#
# # 暴露端口
# EXPOSE 8080
#
# # 启动命令（指定配置文件路径）
# CMD ["/app/myapp", "-config", "/app/config/config.ini"]

