package mq

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRingBuffer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RingBuffer Suite")
}

var _ = Describe("RingBuffer", func() {
	var rb *ringBuffer

	BeforeEach(func() {
		rb = newRingBuffer(4) // 创建大小为4的缓冲区
	})

	Context("创建新的ringBuffer", func() {
		It("应该正确初始化所有字段", func() {
			Expect(rb.size).To(Equal(int64(4)))
			Expect(rb.mask).To(Equal(int64(3))) // 4-1
			Expect(rb.head).To(Equal(int64(0)))
			Expect(rb.tail).To(Equal(int64(0)))
			Expect(rb.items).To(HaveLen(4))
		})

		It("size应该是2的幂", func() {
			rb := newRingBuffer(3)              // 传入3
			Expect(rb.size).To(Equal(int64(4))) // 应该向上取整到4
		})
	})

	Context("写入消息", func() {
		It("应该能正确写入消息", func() {
			msg := &Message{Topic: "test", Value: []byte("hello")}
			err := rb.put(msg)
			Expect(err).NotTo(HaveOccurred())

			item, ok := rb.get(0)
			Expect(ok).To(BeTrue())
			Expect(item.msg).To(Equal(msg))
		})

		It("缓冲区满时应该返回错误", func() {
			// 写满缓冲区(size-1)个消息，因为环形缓冲区要保留一个空位
			for i := 0; i < 3; i++ {
				msg := &Message{Topic: "test", Value: []byte("hello")}
				err := rb.put(msg)
				Expect(err).NotTo(HaveOccurred())
			}

			// 再写入一条应该失败
			msg := &Message{Topic: "test", Value: []byte("hello")}
			err := rb.put(msg)
			Expect(err).To(Equal(ErrBufferFull))
		})
	})

	Context("读取消息", func() {
		It("应该能按顺序读取消息", func() {
			messages := make([]*Message, 3)
			for i := 0; i < 3; i++ {
				messages[i] = &Message{
					Topic: "test",
					Value: []byte("msg"),
				}
				err := rb.put(messages[i])
				Expect(err).NotTo(HaveOccurred())
			}

			// 按顺序读取并验证
			for i := int64(0); i < 3; i++ {
				item, ok := rb.get(i)
				Expect(ok).To(BeTrue())
				Expect(item.msg).To(Equal(messages[i]))
				Expect(item.offset).To(Equal(i))
			}
		})

		It("读取无效偏移量应该返回false", func() {
			_, ok := rb.get(-1)
			Expect(ok).To(BeFalse())

			_, ok = rb.get(100)
			Expect(ok).To(BeFalse())
		})
	})

	Context("并发安全性", func() {
		It("应该在并发写入时保持数据一致性", func() {
			successCount := int32(0)
			receivedCount := int32(0)
			writers := 2
			messagesPerWriter := 2

			// 使用WaitGroup同步写入和读取
			var wg sync.WaitGroup
			wg.Add(writers + 1) // writers + 1个reader

			// 启动读取协程
			go func() {
				defer wg.Done()
				offset := int64(0)
				for {
					if atomic.LoadInt32(&receivedCount) >= atomic.LoadInt32(&successCount) &&
						atomic.LoadInt32(&successCount) > 0 {
						break
					}
					if _, ok := rb.get(offset); ok {
						atomic.AddInt32(&receivedCount, 1)
						offset++
					}
					time.Sleep(time.Millisecond)
				}
			}()

			// 启动写入协程
			for i := 0; i < writers; i++ {
				go func() {
					defer wg.Done()
					for j := 0; j < messagesPerWriter; j++ {
						msg := &Message{Topic: "test", Value: []byte("test")}
						if err := rb.put(msg); err == nil {
							atomic.AddInt32(&successCount, 1)
						}
					}
				}()
			}

			// 等待所有操作完成
			wg.Wait()

			// 验证结果
			Expect(atomic.LoadInt32(&receivedCount)).To(Equal(atomic.LoadInt32(&successCount)))
		})
	})
})
