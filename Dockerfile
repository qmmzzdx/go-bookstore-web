# 使用官方Go 1.24镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置Go代理加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

# 使用国内镜像源安装系统依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk add --no-cache git ca-certificates tzdata make gcc

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bookstore-manager ./cmd/bookstore-manager.go

# 使用官方的alpine镜像作为运行环境
FROM alpine:latest

# 使用国内镜像源安装运行时依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
  apk --no-cache add ca-certificates tzdata curl

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S bookstore && \
  adduser -u 1001 -S bookstore -G bookstore

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder --chown=bookstore:bookstore /app/bookstore-manager .

# 复制配置文件和SQL文件
COPY --chown=bookstore:bookstore conf/ ./conf/
COPY --chown=bookstore:bookstore sql/ ./sql/

# 切换到非root用户
USER bookstore

# 暴露端口
EXPOSE 8080 8081

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/health && \
  curl -f http://localhost:8081/api/v1/admin/health || exit 1

# 启动服务
CMD ["./bookstore-manager"]