package mq

import (
	"context"
	"errors"
)

// Message represents a message in the queue
type Message struct {
	Topic   string
	Key     string
	Value   []byte
	Headers map[string]string
}

// MQError represents message queue related errors
type MQError struct {
	Op  string
	Err error
}

func (e *MQError) Error() string {
	if e.Err != nil {
		return e.Op + ": " + e.Err.Error()
	}
	return e.Op
}

// Producer defines the interface for message publishing operations
type Producer interface {
	// Publish sends a message to the specified topic
	Publish(ctx context.Context, msg *Message) error
	// Close closes the producer connection
	Close() error
}

// Consumer defines the interface for message consuming operations
type Consumer interface {
	// Subscribe 被动订阅模式
	Subscribe(ctx context.Context, topic string, handler func(*Message) error) error
	// Unsubscribe 取消订阅
	Unsubscribe(topic string) error
	// Poll 主动拉取消息，返回消息或超时
	Poll(ctx context.Context, topic string) (*Message, error)
	// Close closes the consumer connection
	Close() error
}

// Config represents configuration for message queue
type Config struct {
	Brokers []string
	Group   string
	// Additional fields can be added based on specific implementation needs
}

var (
	ErrTopicEmpty        = errors.New("topic cannot be empty")
	ErrNilMessage        = errors.New("message cannot be nil")
	ErrNotConnected      = errors.New("message queue not connected")
	ErrUnsubscribeFailed = errors.New("failed to unsubscribe from topic")
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrTimeout           = errors.New("poll timeout")
)

// NewProducer creates a new producer instance
type ProducerFactory interface {
	NewProducer(cfg *Config) (Producer, error)
}

// NewConsumer creates a new consumer instance
type ConsumerFactory interface {
	NewConsumer(cfg *Config) (Consumer, error)
}
