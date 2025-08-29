# 使用 Node.js 官方镜像作为构建环境
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制 package 文件=
COPY package.json package-lock.json ./

# 安装依赖=
RUN npm ci --only=production && \
    npm cache clean --force

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 使用 Nginx 作为生产环境
FROM nginx:alpine

# 设置语言环境和时区
ENV LANG=C.UTF-8 \
    TZ=Asia/Shanghai

# 复制构建产物到 Nginx 目录
COPY --from=builder /app/build /usr/share/nginx/html

# 复制自定义 Nginx 配置（需确保配置了 UTF-8）
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 3000

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]