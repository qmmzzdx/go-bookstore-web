package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// ServerConfig 定义服务器相关配置
// 包含常规服务端口和管理端口配置
type ServerConfig struct {
	Port      int `yaml:"port"`       // 主服务监听端口，默认8080
	AdminPort int `yaml:"admin_port"` // 管理接口监听端口，默认8081
}

// DatabaseConfig 定义数据库连接配置
// 包含MySQL数据库连接所需的所有参数
type DatabaseConfig struct {
	Host     string `yaml:"host"`     // 数据库服务器地址
	Port     int    `yaml:"port"`     // 数据库端口，MySQL默认3306
	User     string `yaml:"user"`     // 数据库用户名
	Password string `yaml:"password"` // 数据库密码
	Name     string `yaml:"name"`     // 要连接的数据库名称
}

// Validate 验证数据库配置完整性
// 返回:
//
//	error - 如果任何必填字段为空或无效则返回错误
func (dbc *DatabaseConfig) Validate() error {
	if dbc.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if dbc.User == "" {
		return fmt.Errorf("database user is required")
	}
	if dbc.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if dbc.Port <= 0 {
		return fmt.Errorf("database port must be positive")
	}
	return nil
}

// RedisConfig 定义Redis连接配置
// 包含连接Redis所需的所有参数
type RedisConfig struct {
	Host     string `yaml:"host"`     // Redis服务器地址
	Port     int    `yaml:"port"`     // Redis端口，默认6379
	Password string `yaml:"password"` // Redis认证密码（可选）
	DB       int    `yaml:"db"`       // 选择的Redis数据库编号，默认0
}

// Validate 验证Redis配置完整性
// 返回:
//
//	error - 如果任何必填字段为空或无效则返回错误
func (rc *RedisConfig) Validate() error {
	if rc.Host == "" {
		return fmt.Errorf("redis host is required")
	}
	if rc.Port <= 0 {
		return fmt.Errorf("redis port must be positive")
	}
	return nil
}

// Config 应用程序主配置结构
// 包含所有子系统的配置信息
type Config struct {
	Server   ServerConfig   `yaml:"server"`   // HTTP服务器配置
	Database DatabaseConfig `yaml:"database"` // 数据库配置
	Redis    RedisConfig    `yaml:"redis"`    // Redis缓存配置
}

// Validate 验证整个应用程序配置
// 会对所有子配置进行验证
// 返回:
//
//	error - 如果任何子配置验证失败则返回错误
func (c *Config) Validate() error {
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}
	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("redis config validation failed: %w", err)
	}
	return nil
}

// AppConfig 全局配置实例
// 在InitConfig初始化后，可以通过该变量访问配置
var AppConfig *Config

// InitConfig 初始化应用程序配置
// 从指定路径加载YAML配置文件并解析验证
// 参数:
//
//	path - 配置文件的路径
//
// 返回:
//
//	error - 加载或解析过程中出现的任何错误
func InitConfig(path string) error {
	// 检查配置文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", path)
	}

	// 读取配置文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析YAML内容到Config结构
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置完整性
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// 设置全局配置实例
	AppConfig = &cfg
	log.Printf("Configuration loaded successfully from: %s", path)
	return nil
}
