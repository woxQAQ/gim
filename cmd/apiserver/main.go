package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-fuego/fuego"
	"github.com/woxQAQ/gim/internal/apiserver/config"
	"github.com/woxQAQ/gim/pkg/db"
	"github.com/woxQAQ/gim/pkg/logger"
	"github.com/woxQAQ/gim/pkg/middleware"
	"go.uber.org/zap"
)

var (
	addr         string
	logLevel     string
	logFile      string
	databasePath string
)

func init() {
	flag.StringVar(&addr, "addr", ":8081", "API服务监听地址")
	flag.StringVar(&logLevel, "log-level", "info", "日志级别 (debug, info, warn, error)")
	flag.StringVar(&logFile, "log-file", "", "日志文件路径，为空时仅输出到控制台")
	flag.StringVar(&databasePath, "db-path", "gim.db", "SQLite数据库文件路径")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 初始化日志系统
	l, err := logger.NewLogger(logger.DomainAPIServer, &logger.Config{
		Level:    logLevel,
		FilePath: logFile,
	}, zap.AddCallerSkip(1))
	if err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库连接
	if err := db.Init(&db.Config{
		DatabasePath: databasePath,
	}); err != nil {
		l.Error("初始化数据库连接失败", logger.Error(err))
		os.Exit(1)
	}

	// 创建Fuego服务器实例
	server := fuego.NewServer(fuego.WithAddr(addr),
		fuego.WithoutLogger(),
		fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			JsonFilePath:     "../../api/openapi.json",
			PrettyFormatJson: true,
		}),
		fuego.WithoutAutoGroupTags(),
	)

	// 设置全局中间件
	fuego.Use(server, middleware.Recovery(l))

	// 设置日志中间件
	fuego.Use(server, middleware.Logger(l))

	config.Register(server, db.GetDB())

	// 启动服务器
	go func() {
		l.Info("API服务器启动", logger.String("addr", addr))
		if err := server.Run(); err != nil {
			l.Error("服务器异常退出", logger.Error(err))
			os.Exit(1)
		}
	}()

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 优雅关闭
	l.Info("正在关闭服务...")
	if err := server.Shutdown(context.Background()); err != nil {
		l.Error("关闭服务器失败", logger.Error(err))
	}

	// 关闭数据库连接
	if err := db.Close(); err != nil {
		l.Error("关闭数据库连接失败", logger.Error(err))
	}

	l.Info("服务已完全关闭")
}
