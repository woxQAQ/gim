package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/woxQAQ/gim/internal/wsgateway"
	"github.com/woxQAQ/gim/pkg/db"
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

	// 初始化数据库
	// 设置测试数据库为内存模式
	Err := db.Init(&db.Config{DatabasePath: ":memory:"})
	Expect(Err).NotTo(HaveOccurred())

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

	// 清理测试数据库
	_ = os.Remove(filepath.Join(os.TempDir(), "gim_test.db"))
})
