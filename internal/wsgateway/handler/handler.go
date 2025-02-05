package handler

import (
	"github.com/woxQAQ/gim/internal/types"
)

// Handler 定义消息处理器接口
type Handler interface {
	// Handle 处理消息，返回是否继续处理链
	Handle(msg types.Message) (bool, error)
	// SetNext 设置下一个处理器
	SetNext(handler Handler)
	// GetNext 获取下一个处理器
	GetNext() Handler
}

// BaseHandler 处理器基础实现
type BaseHandler struct {
	next Handler
}

// SetNext 设置下一个处理器
func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

// GetNext 获取下一个处理器
func (h *BaseHandler) GetNext() Handler {
	return h.next
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
func (c *Chain) Process(msg types.Message) error {
	if c.head == nil {
		return nil
	}

	current := c.head
	for current != nil {
		continue_, err := current.Handle(msg)
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
