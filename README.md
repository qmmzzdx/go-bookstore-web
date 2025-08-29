# 📚 MZZDX书城 - 现代化在线图书商城系统

## 🎯 项目概述

MZZDX书城是一个基于Go语言开发的现代化在线图书商城系统，采用前后端分离架构，提供完整的图书浏览、搜索、购物车、收藏、订单管理等核心功能。系统设计注重用户体验，支持响应式布局，适配多种设备访问。

## 🌟 核心特性

-   **📖 图书管理**: 支持图书分类、搜索、详情展示
-   **🛒 购物车系统**: 本地存储，支持数量调整和持久化
-   **💳 订单管理**: 完整的下单、支付、订单历史流程
-   **👤 用户系统**: 注册、登录、个人资料管理
-   **❤️ 收藏功能**: 用户可收藏喜欢的图书
-   **🎨 轮播展示**: 动态轮播图展示推荐内容
-   **📱 响应式设计**: 完美适配桌面端和移动端
-   **👑 管理员系统**: 完整的后台管理功能

## 🏗️ 技术架构

### 后端技术栈

| 技术          | 版本      | 用途           |
| ------------- | --------- | -------------- |
| **Go**        | 1.24+     | 主要开发语言   |
| **Gin**       | 1.9+      | Web框架        |
| **GORM**      | 1.25+     | ORM框架        |
| **MySQL**     | 5.7+      | 主数据库       |
| **Redis**     | 7.2.0+    | 缓存数据库     |
| **Nginx**     | 1.21.0+   | 代理服务器     |
| **JWT**       | -         | 用户认证       |
| **Validator** | -         | 参数验证       |

## 📋 功能模块

### 1. 用户管理模块
-   用户注册与登录
-   JWT Token认证
-   个人资料管理
-   密码修改

### 2. 图书管理模块
-   图书列表展示
-   分类筛选
-   搜索功能（模糊匹配）
-   图书详情页
-   热销图书推荐
-   新书上架展示

### 3. 购物车模块
-   添加/删除商品
-   数量调整
-   本地存储持久化
-   价格计算

### 4. 订单管理模块
-   订单创建
-   订单支付
-   订单历史查询
-   订单状态管理

### 5. 收藏功能模块
-   收藏/取消收藏
-   收藏列表展示
-   收藏数量统计

### 6. 管理员功能模块
-   **仪表盘**: 统计卡片、最近图书、系统进度、快速操作
-   **图书管理**: 图书列表、增删改查、上下架操作
-   **分类管理**: 分类列表、增删改查、状态管理
-   **订单管理**: 订单列表查看、状态管理、详情查看
-   **用户管理**: 用户列表查看、状态管理、详情查看

## 🎨 设计特色

### 圆角设计
-   所有按钮采用圆角设计 (`border-radius: 6px/8px`)
-   输入框、选择器统一圆角风格
-   卡片、模态框采用圆角设计
-   标签和进度条也采用圆角

### 现代化UI
-   渐变按钮效果
-   阴影和过渡动画
-   优雅的加载状态
-   友好的错误提示

### 用户体验
-   直观的操作界面
-   清晰的信息层级
-   友好的错误提示
-   流畅的交互体验

## 🔗 API接口

### 基础信息

-   **Base URL**: `http://localhost:8080` (用户端), `http://localhost:8081` (管理员端)
-   **API Version**: `v1`
-   **Content-Type**: `application/json`
-   **认证方式**: JWT Token (Bearer Token)

### 响应格式

所有API响应都遵循统一的JSON格式：

```json
{
  "code": 0,            // 状态码：0表示成功，-1表示失败
  "message": "success", // 响应消息
  "data": {}            // 响应数据
}
```

### 管理员API (端口: 8081)

#### 图书管理
-   `GET /api/v1/admin/books/list` - 获取图书列表
-   `GET /api/v1/admin/books/:id` - 获取图书详情
-   `POST /api/v1/admin/books/create` - 创建图书
-   `PUT /api/v1/admin/books/:id` - 更新图书
-   `DELETE /api/v1/admin/books/:id` - 删除图书
-   `PUT /api/v1/admin/books/:id/status` - 更新图书状态

#### 分类管理
-   `GET /api/v1/admin/categories/list` - 获取分类列表
-   `POST /api/v1/admin/categories/create` - 创建分类
-   `PUT /api/v1/admin/categories/:id` - 更新分类
-   `DELETE /api/v1/admin/categories/:id` - 删除分类

#### 订单管理
-   `GET /api/v1/admin/orders/list` - 获取订单列表
-   `GET /api/v1/admin/orders/:id` - 获取订单详情
-   `PUT /api/v1/admin/orders/:id/status` - 更新订单状态

#### 用户管理
-   `GET /api/v1/admin/users/list` - 获取用户列表
-   `GET /api/v1/admin/users/:id` - 获取用户详情
-   `PUT /api/v1/admin/users/:id/status` - 更新用户状态

### 用户端API (端口: 8080)

#### 认证相关
-   `POST /api/v1/user/register` - 用户注册
-   `POST /api/v1/user/login` - 用户登录
-   `GET /api/v1/user/profile` - 获取用户信息

#### 图书相关
-   `GET /api/v1/book/list` - 获取图书列表
-   `GET /api/v1/book/search` - 搜索图书
-   `GET /api/v1/book/detail/{id}` - 获取图书详情
-   `GET /api/v1/book/category/{category}` - 获取分类图书
-   `GET /api/v1/book/hot` - 获取热销图书
-   `GET /api/v1/book/new` - 获取新书

#### 订单相关
-   `POST /api/v1/order/create` - 创建订单
-   `GET /api/v1/order/list` - 获取订单列表
-   `POST /api/v1/order/{id}/pay` - 支付订单

#### 收藏相关
-   `POST /api/v1/favorite/add` - 添加收藏
-   `DELETE /api/v1/favorite/remove` - 取消收藏
-   `GET /api/v1/favorite/list` - 获取收藏列表
-   `GET /api/v1/favorite/count` - 获取收藏数量

#### 其他接口
-   `GET /api/v1/carousel/list` - 获取轮播图列表
-   `GET /api/v1/category/list` - 获取分类列表
-   `GET /api/v1/captcha/generate` - 生成验证码

## 📚 MZZDX书城 API 文档

## 📋 概述

本文档详细描述了MZZDX书城系统的所有API接口，包括请求方法、参数、响应格式和示例。

### 基础信息

- **Base URL**: `http://localhost:8080`
- **API Version**: `v1`
- **Content-Type**: `application/json`
- **认证方式**: JWT Token (Bearer Token)

### 响应格式

所有API响应都遵循统一的JSON格式：

```json
{
  "code": 0,            // 状态码：0表示成功，-1表示失败
  "message": "success", // 响应消息
  "data": {}            // 响应数据
}
```

## 🔐 认证相关

### 用户注册

**接口**: `POST /api/v1/user/register`

**请求参数**:
```json
{
  "username": "testuser",
  "password": "123456",
  "confirm_password": "123456",
  "email": "test@example.com",
  "phone": "12345678901",
  "captcha_id": "captcha_id",
  "captcha_value": "1234"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "注册成功"
}
```

### 用户登录

**接口**: `POST /api/v1/user/login`

**请求参数**:
```json
{
  "username": "testuser",
  "password": "123456",
  "captcha_id": "captcha_id",
  "captcha_value": "1234"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "phone": "12345678901",
      "avatar": "https://example.com/avatar.jpg",
      "is_admin": false
    }
  }
}
```

### 获取用户信息

**接口**: `GET /api/v1/user/profile`

**请求头**: `Authorization: Bearer {token}`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取用户信息成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "phone": "12345678901",
    "avatar": "https://example.com/avatar.jpg",
    "is_admin": false,
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

## 📖 图书相关

### 获取图书列表

**接口**: `GET /api/v1/book/list`

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 12)

**响应示例**:
```json
{
  "code": 0,
  "message": "获取图书列表成功",
  "data": {
    "books": [
      {
        "id": 1,
        "title": "三体",
        "author": "刘慈欣",
        "price": 5900,
        "discount": 80,
        "type": "科幻",
        "stock": 100,
        "cover_url": "https://example.com/cover.jpg",
        "description": "地球文明与三体文明的星际战争...",
        "isbn": "9787536692930",
        "publisher": "重庆出版社",
        "publish_date": "2008-01-01",
        "pages": 302,
        "language": "中文",
        "format": "平装"
      }
    ],
    "total": 20,
    "total_page": 2,
    "current_page": 1
  }
}
```

### 搜索图书

**接口**: `GET /api/v1/book/search`

**查询参数**:
- `q`: 搜索关键词
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 12)

**响应示例**:
```json
{
  "code": 0,
  "message": "搜索成功",
  "data": {
    "books": [...],
    "total": 5,
    "total_page": 1,
    "current_page": 1
  }
}
```

### 获取图书详情

**接口**: `GET /api/v1/book/detail/{id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取图书详情成功",
  "data": {
    "id": 1,
    "title": "三体",
    "author": "刘慈欣",
    "price": 5900,
    "discount": 80,
    "type": "科幻",
    "stock": 100,
    "cover_url": "https://example.com/cover.jpg",
    "description": "地球文明与三体文明的星际战争...",
    "isbn": "9787536692930",
    "publisher": "重庆出版社",
    "publish_date": "2008-01-01",
    "pages": 302,
    "language": "中文",
    "format": "平装"
  }
}
```

### 获取分类图书

**接口**: `GET /api/v1/book/category/{category}`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取分类图书成功",
  "data": [
    {
      "id": 1,
      "title": "三体",
      "author": "刘慈欣",
      "price": 5900,
      "discount": 80,
      "type": "科幻",
      "stock": 100,
      "cover_url": "https://example.com/cover.jpg"
    }
  ]
}
```

### 获取热销图书

**接口**: `GET /api/v1/book/hot`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取热销图书成功",
  "data": [
    {
      "id": 1,
      "title": "三体",
      "author": "刘慈欣",
      "price": 5900,
      "discount": 80,
      "type": "科幻",
      "stock": 100,
      "cover_url": "https://example.com/cover.jpg"
    }
  ]
}
```

### 获取新书

**接口**: `GET /api/v1/book/new`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取新书成功",
  "data": [
    {
      "id": 1,
      "title": "三体",
      "author": "刘慈欣",
      "price": 5900,
      "discount": 80,
      "type": "科幻",
      "stock": 100,
      "cover_url": "https://example.com/cover.jpg"
    }
  ]
}
```

## 🛒 订单相关

### 创建订单

**接口**: `POST /api/v1/order/create`

**请求头**: `Authorization: Bearer {token}`

**请求参数**:
```json
{
  "items": [
    {
      "book_id": 1,
      "quantity": 2,
      "price": 4720
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建订单成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "order_no": "ORD202401010001",
    "total_amount": 9440,
    "status": 0,
    "is_paid": false,
    "created_at": "2024-01-01T10:00:00Z",
    "order_items": [
      {
        "id": 1,
        "book_id": 1,
        "quantity": 2,
        "price": 4720,
        "subtotal": 9440,
        "book": {
          "title": "三体",
          "author": "刘慈欣",
          "cover_url": "https://example.com/cover.jpg"
        }
      }
    ]
  }
}
```

### 获取订单列表

**接口**: `GET /api/v1/order/list`

**请求头**: `Authorization: Bearer {token}`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取订单列表成功",
  "data": {
    "orders": [
      {
        "id": 1,
        "order_no": "ORD202401010001",
        "total_amount": 9440,
        "status": 1,
        "is_paid": true,
        "created_at": "2024-01-01T10:00:00Z",
        "payment_time": "2024-01-01T10:05:00Z",
        "order_items": [
          {
            "id": 1,
            "book_id": 1,
            "quantity": 2,
            "price": 4720,
            "subtotal": 9440,
            "book": {
              "title": "三体",
              "author": "刘慈欣",
              "cover_url": "https://example.com/cover.jpg"
            }
          }
        ]
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

### 支付订单

**接口**: `POST /api/v1/order/{id}/pay`

**请求头**: `Authorization: Bearer {token}`

**响应示例**:
```json
{
  "code": 0,
  "message": "支付成功"
}
```

## ❤️ 收藏相关

### 添加收藏

**接口**: `POST /api/v1/favorite/add`

**请求头**: `Authorization: Bearer {token}`

**请求参数**:
```json
{
  "book_id": 1
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "添加收藏成功"
}
```

### 取消收藏

**接口**: `DELETE /api/v1/favorite/remove`

**请求头**: `Authorization: Bearer {token}`

**请求参数**:
```json
{
  "book_id": 1
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "取消收藏成功"
}
```

### 获取收藏列表

**接口**: `GET /api/v1/favorite/list`

**请求头**: `Authorization: Bearer {token}`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取收藏列表成功",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "book_id": 1,
      "created_at": "2024-01-01T10:00:00Z",
      "book": {
        "id": 1,
        "title": "三体",
        "author": "刘慈欣",
        "price": 5900,
        "discount": 80,
        "type": "科幻",
        "stock": 100,
        "cover_url": "https://example.com/cover.jpg"
      }
    }
  ]
}
```

### 获取收藏数量

**接口**: `GET /api/v1/favorite/count`

**请求头**: `Authorization: Bearer {token}`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取收藏数量成功",
  "data": {
    "count": 5
  }
}
```

## 🎨 轮播图相关

### 获取轮播图列表

**接口**: `GET /api/v1/carousel/list`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取轮播图成功",
  "data": [
    {
      "id": 1,
      "title": "精选好书推荐",
      "description": "发现更多精彩好书",
      "image_url": "https://example.com/carousel1.jpg",
      "link_url": "/category/文学",
      "sort_order": 1,
      "is_active": true
    }
  ]
}
```

## 📂 分类相关

### 获取分类列表

**接口**: `GET /api/v1/category/list`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取分类列表成功",
  "data": [
    {
      "id": 1,
      "name": "科幻",
      "description": "科幻小说和未来科技",
      "icon": "🚀",
      "color": "#FF6B6B",
      "gradient": "linear-gradient(135deg, #FF6B6B, #FF8E8E)",
      "sort": 1,
      "is_active": true,
      "book_count": 4
    }
  ]
}
```

## 🔐 验证码相关

### 生成验证码

**接口**: `GET /api/v1/captcha/generate`

**响应示例**:
```json
{
  "code": 0,
  "message": "验证码生成成功",
  "data": {
    "captcha_id": "3JU9cHTsexsCQymMZvCb",
    "captcha_base64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAPAAAABQCAMAAAAQlwhOAAAA81BMVEUAAAABfHdm4..."
  }
}
```

## ⚠️ 错误码说明

| 错误码 | 说明                 | HTTP状态码对应 |
| ------ | -------------------- | -------------- |
| 0      | 成功                 | 200 OK         |
| -1     | 失败                 | 400 Bad Request|
| 400    | 请求参数错误         | 400 Bad Request|
| 401    | 未授权               | 401 Unauthorized|
| 403    | 禁止访问             | 403 Forbidden  |
| 404    | 资源不存在           | 404 Not Found  |
| 500    | 服务器内部错误       | 500 Internal Server Error |

## 📝 注意事项

1.  **认证**: 需要认证的接口必须在请求头中包含有效的JWT Token
2.  **价格**: 所有价格字段都以分为单位存储，前端显示时需要除以100
3.  **分页**: 分页接口支持page和page_size参数
4.  **时间格式**: 所有时间字段使用ISO 8601格式
5.  **图片URL**: 图片URL需要完整的HTTPS地址
6.  **参数验证**: 所有输入参数都需要进行验证，防止SQL注入和XSS攻击
7.  **安全性**: 使用HTTPS加密传输数据，对敏感数据进行加密

## 🔧 开发环境

-   **后端服务**: http://localhost:8080
-   **管理员后端服务**: http://localhost:8081
-   **前端服务**: http://localhost:3000
-   **数据库**: MySQL 8.0+
-   **缓存数据库**: Redis 6.0+
-   **API测试工具**: Postman, curl

## 🚀 快速开始

### 环境要求

-   Go 1.21+
-   MySQL 5.7+
-   Redis 6.0+

### 1. 克隆项目

```bash
git clone https://gitee.com/llinfan/bookstore-go
cd bookstore-go
```

### 2. 后端配置

```bash
# 安装依赖
go mod tidy

# 编辑配置文件，设置数据库连接信息
vim conf/config.yaml

# 配置文件示例
server:
  port: 8080        # 用户端端口
  admin_port: 8081  # 管理员端口

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: your_password
  name: bookstore

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

# 初始化数据库
mysql -u root -p < sql/bookstore.sql
mysql -u root -p < sql/mock.sql

# 编译服务
make bookstore-manager

# 启动后端服务
./bin/bookstore-manager
```

### 3. 访问应用

-   **用户端后端API**: http://localhost:8080
-   **管理员后端API**: http://localhost:8081

## 📁 项目结构

```
bookstore-web/
├── cmd/                    # 程序入口
│   ├── bookstore-manager/  # 用户端服务入口
│   └── admin-manager/      # 管理员服务入口
├── model/                  # 数据模型
├── service/                # 业务逻辑层
├── repository/             # 数据访问层
├── web/                    # Web层
│   ├── controller/         # 控制器
│   │   ├── book.go         # 用户端图书控制器
│   │   ├── admin_book.go   # 管理员图书控制器
│   │   ├── admin_order.go  # 管理员订单控制器
│   │   └── admin_user.go   # 管理员用户控制器
│   ├── router/             # 路由配置
│   │   ├── router.go       # 用户端路由
│   │   └── admin_router.go # 管理员路由
│   └── middleware/         # 中间件
├── conf/                   # 配置文件
└── sql/                    # 数据库脚本
```

## 🔒 安全特性

-   **密码加密**: 使用bcrypt进行密码哈希
-   **JWT认证**: 基于Token的无状态认证
-   **SQL注入防护**: 使用GORM参数化查询，避免SQL注入风险
-   **CORS配置**: 使用`github.com/rs/cors`中间件配置跨域请求，支持精确的域名控制
-   **输入验证**: 使用`github.com/go-playground/validator`进行请求参数严格验证
-   **HTTPS支持**: 可通过Nginx配置HTTPS反向代理，保障数据传输安全

## 📖 开发指南

### 添加新的管理员功能

1.  **添加模型定义** (model/)
2.  **添加业务逻辑** (service/)
3.  **添加控制器** (web/controller/admin_*.go)
4.  **添加路由** (web/router/admin_router.go)
5.  **添加前端页面** (bookstore-admin-frontend/src/pages/)

### 安全编程实践

-   使用参数化查询防止SQL注入
-   对用户输入进行严格验证和过滤
-   使用白名单验证列名和输入值
-   避免直接拼接SQL语句

## 🚀 部署指南

### 生产环境部署

1.  **配置生产环境数据库**: 设置正确的数据库连接参数
2.  **设置环境变量**: 区分开发和生产环境配置
3.  **进程守护**: 使用PM2或systemd管理进程，确保服务持续运行
4.  **反向代理**: 配置Nginx反向代理和HTTPS加密
5.  **监控与审计**: 记录所有数据库操作，设置异常查询告警

### 使用Systemd进行进程守护

创建service文件`/etc/systemd/system/bookstore.service`:

```ini
[Unit]
Description=MZZDX Bookstore Service
After=network.target

[Service]
Type=simple
ExecStart=/path/to/bookstore-manager
Restart=always
RestartSec=5s
User=www-data
Group=www-data

[Install]
WantedBy=multi-user.target
```

### Nginx HTTPS配置示例

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate /path/to/your/certificate.crt;
    ssl_certificate_key /path/to/your/private.key;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /admin/ {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```
