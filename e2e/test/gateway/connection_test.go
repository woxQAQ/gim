package gateway

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/woxQAQ/gim/e2e/pkg/client"
)

var _ = Describe("WebSocket Gateway Connection Tests", func() {
	var clients []*client.Client

	BeforeEach(func() {
		clients = make([]*client.Client, 0)
	})

	AfterEach(func() {
		// 清理客户端连接
		for _, c := range clients {
			c.Close()
		}
		clients = nil
	})

	Context("基本连接测试", func() {
		It("应该能成功建立连接", func() {
			c := client.New(baseURL+"?user_id=test1", "test1", 1)
			Err := c.Connect()
			Expect(Err).NotTo(HaveOccurred())
			clients = append(clients, c)

			// 验证用户在线状态
			Eventually(func() bool {
				return gateway.IsUserOnline("test1")
			}).Should(BeTrue())
		})
	})
})
