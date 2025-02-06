package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/woxQAQ/gim/internal/wsgateway"
	"github.com/woxQAQ/gim/pkg/logger"
)

var (
	gateway *wsgateway.WSGateway
	server  *httptest.Server
	baseURL string
	testCtx context.Context
	cancel  context.CancelFunc
)

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gateway Suite")
}

var _ = BeforeSuite(func() {
	testCtx, cancel = context.WithCancel(context.Background())

	// 创建网关实例
	l, _ := logger.NewLogger(logger.DomainWSGateway, &logger.Config{Level: "error"})
	l.Disable()
	gateway, _ = wsgateway.NewWSGateway(wsgateway.WithLogger(l))
	_ = gateway.Start(testCtx)

	// 创建测试服务器
	server = httptest.NewServer(http.HandlerFunc(gateway.HandleNewConnection))
	baseURL = "ws" + server.URL[4:]
})

var _ = AfterSuite(func() {
	// 停止服务
	cancel()
	_ = gateway.Stop()
	server.Close()
})
