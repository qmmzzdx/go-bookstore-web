# 使用官方Go 1.24镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖
RUN apk add --no-cache git ca-certificates tzdata make gcc vim

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bookstore-manager ./cmd/bookstore-manager/

# 使用轻量级的alpine镜像作为运行环境
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata wget

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S bookstore && \
  adduser -u 1001 -S bookstore -G bookstore

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bookstore-manager .

# 复制配置文件
COPY conf/ ./conf/

# 复制SQL文件
COPY sql/ ./sql/

# 设置权限
RUN chown -R bookstore:bookstore /app

# 切换到非root用户
USER bookstore

# 暴露端口（用户端和管理员端）
EXPOSE 8080 8081

# 健康检查 - 同时检查用户端和管理员端
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health && \
  wget --no-verbose --tries=1 --spider http://localhost:8081/api/v1/admin/health || exit 1

# 启动服务（同时启动用户端和管理员端）
CMD ["./bookstore-manager"]