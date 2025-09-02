package global

import (
	"context"
	"fmt"
	"log"

	"bookstore/config"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBClient 全局数据库客户端实例
// 用于整个应用程序的数据库操作，通过GetDB()获取
var DBClient *gorm.DB

// RedisClient 全局Redis客户端实例
// 用于整个应用程序的Redis操作
var RedisClient *redis.Client

// GetDB 获取全局数据库连接对象
// 返回值:
//
//	*gorm.DB - 数据库连接实例
//
// 注意: 如果DBClient未初始化会触发log.Fatalln
func GetDB() *gorm.DB {
	if DBClient == nil {
		log.Fatalln("DBClient is not initialized")
	}
	return DBClient
}

// InitMySQL 初始化MySQL数据库连接
// 从全局配置中读取数据库配置，建立连接并设置GORM的日志模式
// 依赖:
//   - config.AppConfig.Database 必须已正确配置
//
// 副作用:
//   - 初始化全局变量DBClient
//   - 连接失败会终止程序
func InitMySQL() {
	// 从全局配置中获取数据库配置
	cfg := config.AppConfig.Database

	// 构建MySQL连接字符串(DSN)
	// 格式: "用户名:密码@协议(地址:端口)/数据库名?参数"
	// 参数说明:
	//   - charset=utf8mb4: 使用完整的UTF-8编码(支持4字节字符，如emoji表情)
	//   - parseTime=True: 将数据库中的DATE/DATETIME/TIMESTAMP类型自动解析为Go的time.Time类型
	//   - loc=Local: 使用时区本地化处理时间，确保时间字段与应用程序所在时区一致
	//   - collation=utf8mb4_unicode_ci: 设置字符集排序规则为Unicode大小写不敏感，支持多语言排序
	//   - tcp: 使用TCP协议进行数据库连接（默认协议）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
		cfg.User,     // 数据库用户名
		cfg.Password, // 数据库密码
		cfg.Host,     // 数据库主机地址
		cfg.Port,     // 数据库端口
		cfg.Name,     // 数据库名称
	)

	// 连接数据库
	var err error
	// 使用GORM打开MySQL连接，并配置:
	// - 使用mysql驱动打开指定DSN
	// - 设置GORM配置: 启用Info级别的SQL日志
	DBClient, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 开发环境显示SQL日志
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("MySQL connected successfully")
}

// InitRedis 初始化Redis连接
// 从全局配置中读取Redis配置，建立连接并测试连通性
// 依赖:
//   - config.AppConfig.Redis 必须已正确配置
//
// 副作用:
//   - 初始化全局变量RedisClient
//   - 连接失败会终止程序
func InitRedis() {
	// 获取Redis配置
	cfg := config.AppConfig.Redis

	// 创建Redis客户端实例
	// 配置选项包括:
	//   - Addr: Redis服务器地址(主机:端口)
	//   - Password: Redis认证密码
	//   - DB: 选择的Redis数据库编号
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // Redis服务器地址
		Password: cfg.Password,                             // Redis密码
		DB:       cfg.DB,                                   // Redis数据库编号
	})

	// 测试连接，使用context.Background()来传递上下文
	ctx := context.Background()
	str, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	log.Printf("Redis connected successfully, ping result: %s", str)
}

// CloseDB 安全关闭全局数据库连接
// 释放数据库连接资源
// 注意:
//   - 应该在程序退出前调用
//   - 无错误返回，仅做资源清理
func CloseDB() {
	if DBClient != nil {
		sqlDB, err := DBClient.DB()
		if err == nil {
			sqlDB.Close() // 关闭底层数据库连接
		}
	}
}

// CloseRedis 安全关闭全局Redis连接
// 释放Redis连接资源
// 注意:
//   - 应该在程序退出前调用
//   - 无错误返回，仅做资源清理
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close() // 关闭Redis客户端连接
	}
}
