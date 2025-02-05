package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/woxQAQ/gim/internal/wsgateway"
	"github.com/woxQAQ/gim/pkg/db"
	"github.com/woxQAQ/gim/pkg/logger"
	"github.com/woxQAQ/gim/pkg/snowflake"
	"go.uber.org/zap"
)

var (
	addr         string
	logLevel     string
	logFile      string
	databasePath string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "WebSocket服务监听地址")
	flag.StringVar(&logLevel, "log-level", "info", "日志级别 (debug, info, warn, error)")
	flag.StringVar(&logFile, "log-file", "", "日志文件路径，为空时仅输出到控制台")
	flag.StringVar(&databasePath, "db-path", "gim.db", "SQLite数据库文件路径")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 初始化日志系统
	l, err := logger.NewLogger(logger.DomainWSGateway, &logger.Config{
		Level:    logLevel,
		FilePath: logFile,
	}, zap.AddCallerSkip(1))
	if err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化消息ID生成器
	if err := snowflake.InitGenerator(1); err != nil {
		l.Error("初始化消息ID生成器失败", logger.Error(err))
		os.Exit(1)
	}

	// 初始化数据库连接
	if err := db.Init(&db.Config{
		DatabasePath: databasePath,
	}); err != nil {
		l.Error("初始化数据库连接失败", logger.Error(err))
		os.Exit(1)
	}

	// 创建网关实例
	gateway, err := wsgateway.NewWSGateway(
		wsgateway.WithLogger(l),
	)
	if err != nil {
		l.Error("创建WebSocket网关失败", logger.Error(err))
		os.Exit(1)
	}

	// 创建HTTP服务器
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gateway.HandleNewConnection(w, r)
		}),
		ReadHeaderTimeout: 10 * time.Second, // 防止 Slowloris 攻击
	}

	// 启动网关服务
	ctx := context.Background()
	if err := gateway.Start(ctx); err != nil {
		l.Error("启动网关服务失败", logger.Error(err))
		os.Exit(1)
	}

	// 启动HTTP服务器
	go func() {
		l.Info("WebSocket服务器启动", logger.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("HTTP服务器异常退出", logger.Error(err))
			os.Exit(1)
		}
	}()

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 优雅关闭
	l.Info("正在关闭服务...")

	// 关闭HTTP服务器
	if err := server.Shutdown(context.Background()); err != nil {
		l.Error("关闭HTTP服务器失败", logger.Error(err))
	}

	// 关闭网关服务
	if err := gateway.Stop(); err != nil {
		l.Error("关闭网关服务失败", logger.Error(err))
	}

	// 关闭数据库连接
	if err := db.Close(); err != nil {
		l.Error("关闭数据库连接失败", logger.Error(err))
	}

	l.Info("服务已完全关闭")
}
