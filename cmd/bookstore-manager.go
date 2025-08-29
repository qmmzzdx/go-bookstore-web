package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"bookstore/config"
	"bookstore/global"
	"bookstore/web/router"
)

// init 应用程序初始化函数
// 在main函数执行前被自动调用，用于初始化全局变量
func init() {
	// 初始化全局数据库客户端为nil
	global.DBClient = nil
	// 初始化全局Redis客户端为nil
	global.RedisClient = nil
}

// main 应用程序主入口函数
// 负责初始化配置、数据库连接，启动HTTP服务器并处理优雅关闭
func main() {
	// 初始化应用程序配置，从指定路径加载YAML配置文件
	config.InitConfig("conf/config.yaml")
	// 获取全局配置实例
	cfg := config.AppConfig

	// 初始化MySQL数据库连接
	global.InitMySQL()

	// 初始化Redis连接
	global.InitRedis()

	// 创建等待组，用于等待所有服务器关闭
	var wg sync.WaitGroup

	// 初始化主服务路由和管理员路由
	mainRouter := router.InitRouter()
	adminRouter := router.InitAdminRouter()

	// 创建两个HTTP服务器实例
	mainServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mainRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	adminServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.AdminPort),
		Handler:      adminRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// 创建系统信号通道，用于接收操作系统信号
	quit := make(chan os.Signal, 1)
	// 注册感兴趣的信号类型
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	// 启动主服务
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("🚀 主服务启动成功，端口: %d", cfg.Server.Port)
		log.Printf("📖 访问地址: http://localhost:%d", cfg.Server.Port)
		if err := mainServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("主服务启动失败: %v", err)
		}
	}()

	// 启动管理员服务
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("🚀 管理员系统启动成功，端口: %d", cfg.Server.AdminPort)
		log.Printf("📖 访问地址: http://localhost:%d", cfg.Server.AdminPort)
		if err := adminServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("管理员系统启动失败: %v", err)
		}
	}()

	// 等待接收操作系统信号
	sig := <-quit
	log.Printf("收到信号: %v", sig)

	// 根据接收到的信号类型输出相应的日志信息
	switch sig {
	case syscall.SIGINT:
		log.Println("收到 SIGINT (Ctrl+C)，正在关闭服务器...")
	case syscall.SIGTERM:
		log.Println("收到 SIGTERM，正在关闭服务器...")
	case syscall.SIGHUP:
		log.Println("收到 SIGHUP，正在关闭服务器...")
	case syscall.SIGQUIT:
		log.Println("收到 SIGQUIT，正在关闭服务器...")
	default:
		log.Printf("收到信号 %v，正在关闭服务器...", sig)
	}

	// 创建带有5秒超时的上下文，用于服务器优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	log.Println("正在关闭主服务...")
	if err := mainServer.Shutdown(ctx); err != nil {
		log.Printf("主服务强制关闭: %v", err)
	} else {
		log.Println("主服务优雅停止")
	}

	log.Println("正在关闭管理员服务...")
	if err := adminServer.Shutdown(ctx); err != nil {
		log.Printf("管理员服务强制关闭: %v", err)
	} else {
		log.Println("管理员服务优雅停止")
	}

	// 等待所有服务器关闭完成
	wg.Wait()

	// 清理应用程序资源（数据库连接、Redis连接等）
	log.Println("正在清理资源...")
	cleanupResources()

	// 记录应用程序成功退出日志
	log.Println("所有服务退出成功")

	// 退出程序，确保所有goroutine和资源被释放
	os.Exit(0)
}

// cleanupResources 清理应用程序资源
// 负责关闭数据库连接、Redis连接等资源释放操作
func cleanupResources() {
	// 关闭数据库连接
	if global.DBClient != nil {
		log.Println("正在关闭数据库连接...")
		global.CloseDB()
	}

	// 关闭Redis连接
	if global.RedisClient != nil {
		log.Println("正在关闭Redis连接...")
		global.CloseRedis()
	}

	// 等待一小段时间确保所有资源都被正确释放
	time.Sleep(100 * time.Millisecond)
	log.Println("所有资源已清理完成")
}
