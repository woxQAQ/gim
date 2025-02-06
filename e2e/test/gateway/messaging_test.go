package gateway

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"time"

	"github.com/woxQAQ/gim/e2e/pkg/client"
	"github.com/woxQAQ/gim/internal/types"
)

var _ = Describe("WebSocket Gateway Messaging Tests", func() {
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

	Context("消息发送测试", func() {
		It("应该能成功发送和接收消息", func() {
			// 创建两个测试客户端
			client1 := client.New(baseURL+"?user_id=test1", "test1", 1)
			client2 := client.New(baseURL+"?user_id=test2", "test2", 1)

			Err := client1.Connect()
			Expect(Err).NotTo(HaveOccurred())
			Err = client2.Connect()
			Expect(Err).NotTo(HaveOccurred())

			clients = append(clients, client1, client2)

			// 等待连接建立
			time.Sleep(100 * time.Millisecond)

			// 发送测试消息
			testMsg := types.Message{
				Header: types.MessageHeader{
					Type: types.MessageTypeText,
					From: "test1",
					To:   "test2",
				},
				Payload: []byte("hello"),
			}

			Err = client1.SendMessage(testMsg)
			Expect(Err).NotTo(HaveOccurred())

			// 验证消息接收
			Eventually(func() []types.Message {
				return client2.GetMessages()
			}).Should(ContainElement(testMsg))
		})
	})

	Context("心跳测试", func() {
		It("应该正确处理心跳消息", func() {
			c := client.New(baseURL+"?user_id=test1", "test1", 1)
			Err := c.Connect()
			Expect(Err).NotTo(HaveOccurred())
			clients = append(clients, c)

			// 发送心跳消息
			heartbeat := types.Message{
				Header: types.MessageHeader{
					Type: types.MessageTypeHeartbeat,
					From: "test1",
				},
			}

			Err = c.SendMessage(heartbeat)
			Expect(Err).NotTo(HaveOccurred())

			// 验证最后心跳时间更新
			Eventually(func() error {
				_, err := gateway.GetUserHeartbeatStatus("test1", 1)
				return err
			}).ShouldNot(HaveOccurred())
		})
	})
})
