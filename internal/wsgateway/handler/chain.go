package handler

import (
	"github.com/woxQAQ/gim/internal/apiserver/stores"
	"github.com/woxQAQ/gim/internal/wsgateway/codec"
	"github.com/woxQAQ/gim/internal/wsgateway/user"
	"github.com/woxQAQ/gim/pkg/mq"
)

// NewMessageChain 创建默认的消息处理链
func NewMessageChain(userManager user.IUserManager,
	ms *stores.MessageStore,
	encoder codec.Encoder,
	producer mq.Producer,
) *Chain {
	chain := NewChain()

	// 添加异步消息转发处理器
	chain.AddHandler(NewForwardHandler(userManager, producer, encoder))

	// 添加消息存储处理器
	chain.AddHandler(NewStoreHandler(ms, encoder))

	return chain
}

// Chain 消息处理责任链
type Chain struct {
	head Handler
	tail Handler
}

// NewChain 创建新的处理链
func NewChain() *Chain {
	return &Chain{}
}

// AddHandler 添加处理器到链尾
func (c *Chain) AddHandler(handler Handler) {
	if c.head == nil {
		c.head = handler
		c.tail = handler
		return
	}

	c.tail.SetNext(handler)
	c.tail = handler
}

// Process 处理消息
func (c *Chain) Process(data []byte) error {
	if c.head == nil {
		return nil
	}

	current := c.head
	for current != nil {
		continue_, err := current.Handle(data)
		if err != nil {
			return err
		}
		if !continue_ {
			break
		}
		current = current.GetNext()
	}

	return nil
}
