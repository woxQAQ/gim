package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/server/gateway"
	"go.uber.org/zap"
)

func main() {
	wsMgr := gateway.NewWsMgr()
	wsMgr.Run()
	gsServer := gateway.NewGatewayServer("tcp", true, wsMgr)
	gsServer.WsEngine.Run(gsServer.WebsocketAddress)
	websocketServer := &http.Server{
		Addr:    gsServer.WebsocketAddress,
		Handler: gsServer.WsEngine,
	}

	err := gnet.Run(gsServer, gsServer.Network+"://"+gsServer.TcpAddress)
	if err != nil {
		logging.Errorf("gatewayServer start error: %v", err)
		return
	}

	go func() {
		if err := websocketServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := websocketServer.Shutdown(ctx); err != nil {
		zap.S().Fatalf("Server shutdown", err)
	}
	zap.S().Infoln("server exiting")
}
