package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"github.com/woxQAQ/gim/internal/apiserver/config"
	"github.com/woxQAQ/gim/pkg/constants"
	"github.com/woxQAQ/gim/pkg/db"
	"github.com/woxQAQ/gim/pkg/logger"
)

func init() {
	viper.SetDefault(constants.OpenapiFilePath, "../../api/openapi.json")
	viper.SetDefault(constants.EnableOpenapiSpec, true)
	viper.SetDefault(constants.OpenapiRoute, "/swagger")
	viper.SetDefault(constants.LogLevel, "info")
	viper.SetDefault(constants.LogFilePath, "")
	viper.SetDefault(constants.ApiPath, ":8081")
	viper.SetDefault(constants.DBPath, "gim.db")
}

func main() {
	// 初始化日志系统
	l := config.SetupLogger()

	// 初始化数据库连接
	config.SetupDatabase(l)
	server := config.SetupApiServer(l)

	config.Register(server, db.GetDB())

	// 启动服务器
	go func() {
		l.Info("API服务器启动", logger.String("addr", server.Addr))
		if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
