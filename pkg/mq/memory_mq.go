package mq

import (
	"context"
	"sync"
)

type memoryMQ struct {
	msgs     map[string]chan *Message
	handlers map[string][]func(*Message) error
	mu       sync.RWMutex
}

type memoryProducer struct {
	mq *memoryMQ
}

type memoryConsumer struct {
	mq    *memoryMQ
	group string
}

type MemoryMQFactory struct {
	mq *memoryMQ
}

func NewMemoryMQFactory() *MemoryMQFactory {
	return &MemoryMQFactory{
		mq: &memoryMQ{
			msgs:     make(map[string]chan *Message),
			handlers: make(map[string][]func(*Message) error),
		},
	}
}

func (f *MemoryMQFactory) NewProducer(_ *Config) (Producer, error) {
	return &memoryProducer{mq: f.mq}, nil
}

func (f *MemoryMQFactory) NewConsumer(_ *Config) (Consumer, error) {
	return &memoryConsumer{mq: f.mq}, nil
}

// Producer implementation
func (p *memoryProducer) Publish(ctx context.Context, msg *Message) error {
	if msg == nil {
		return ErrNilMessage
	}
	if msg.Topic == "" {
		return ErrTopicEmpty
	}

	p.mq.mu.Lock()
	if _, exists := p.mq.msgs[msg.Topic]; !exists {
		p.mq.msgs[msg.Topic] = make(chan *Message, 100) // 缓冲区大小可配置
	}
	p.mq.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.mq.msgs[msg.Topic] <- msg:
		// 触发所有该主题的处理器
		p.mq.mu.RLock()
		handlers := p.mq.handlers[msg.Topic]
		p.mq.mu.RUnlock()

		for _, handler := range handlers {
			go handler(msg)
		}
		return nil
	}
}

func (p *memoryProducer) Close() error {
	return nil
}

// Consumer implementation
func (c *memoryConsumer) Subscribe(ctx context.Context, topic string, handler func(*Message) error) error {
	if topic == "" {
		return ErrTopicEmpty
	}

	c.mq.mu.Lock()
	if _, exists := c.mq.msgs[topic]; !exists {
		c.mq.msgs[topic] = make(chan *Message, 100)
	}
	c.mq.handlers[topic] = append(c.mq.handlers[topic], handler)
	c.mq.mu.Unlock()

	return nil
}

func (c *memoryConsumer) Unsubscribe(topic string) error {
	c.mq.mu.Lock()
	delete(c.mq.handlers, topic)
	c.mq.mu.Unlock()
	return nil
}

func (c *memoryConsumer) Close() error {
	return nil
}

// Poll 实现主动消费模式
func (c *memoryConsumer) Poll(ctx context.Context, topic string) (*Message, error) {
	if topic == "" {
		return nil, ErrTopicEmpty
	}

	c.mq.mu.RLock()
	msgChan, exists := c.mq.msgs[topic]
	c.mq.mu.RUnlock()

	if !exists {
		c.mq.mu.Lock()
		c.mq.msgs[topic] = make(chan *Message, 100)
		msgChan = c.mq.msgs[topic]
		c.mq.mu.Unlock()
	}

	select {
	case msg := <-msgChan:
		return msg, nil
	case <-ctx.Done():
		return nil, ErrTimeout
	}
}
