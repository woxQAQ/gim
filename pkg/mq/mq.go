package mq

import (
	"context"
	"errors"
	"time"
)

// Message 表示队列中的一条消息
type Message struct {
	Topic       string            // 主题
	Key         string            // 消息键
	Value       []byte            // 消息内容
	Headers     map[string]string // 消息头
	Priority    int               // 消息优先级 0-9，越大优先级越高
	Expiration  time.Duration     // 消息过期时间
	CreateTime  time.Time         // 消息创建时间
	DeliverTime time.Time         // 期望投递时间（延迟消息）
}

// Result 表示消息处理结果
type Result struct {
	MessageID string
	Topic     string
	Error     error
}

// Producer 定义生产者接口
type Producer interface {
	// Publish 发送单条消息
	Publish(ctx context.Context, msg *Message) (Result, error)
	// PublishBatch 批量发送消息
	PublishBatch(ctx context.Context, msgs []*Message) ([]Result, error)
	// Close 关闭生产者连接
	Close() error
}

// DeliveryMode 消息投递模式
type DeliveryMode int

const (
	// AtMostOnce 最多投递一次
	AtMostOnce DeliveryMode = iota
	// AtLeastOnce 至少投递一次
	AtLeastOnce
	// ExactlyOnce 精确投递一次
	ExactlyOnce
)

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	// 消费者组ID
	GroupID string
	// 消息投递模式
	DeliveryMode DeliveryMode
	// 批量消费大小
	BatchSize int
	// 消费超时时间
	Timeout time.Duration
	// 是否自动提交
	AutoCommit bool
}

// Consumer 定义消费者接口
type Consumer interface {
	// Subscribe 订阅主题
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	// SubscribeBatch 批量订阅主题
	SubscribeBatch(ctx context.Context, topic string, handler BatchMessageHandler) error
	// Unsubscribe 取消订阅
	Unsubscribe(topic string) error
	// Commit 手动提交消息确认
	Commit(ctx context.Context, msgID string) error
	// Reject 拒绝消息
	Reject(ctx context.Context, msgID string, requeue bool) error
	// Pause 暂停消费
	Pause(topics ...string) error
	// Resume 恢复消费
	Resume(topics ...string) error
	// Close 关闭消费者
	Close() error
}

// MessageHandler 消息处理函数
type MessageHandler func(ctx context.Context, msg *Message) error

// BatchMessageHandler 批量消息处理函数
type BatchMessageHandler func(ctx context.Context, msgs []*Message) error

// Admin 定义管理接口
type Admin interface {
	// CreateTopic 创建主题
	CreateTopic(ctx context.Context, topic string) error
	// DeleteTopic 删除主题
	DeleteTopic(ctx context.Context, topic string) error
	// ListTopics 列出所有主题
	ListTopics(ctx context.Context) ([]string, error)
	// GetOffset 获取主题偏移量
	GetOffset(ctx context.Context, topic string, groupID string) (int64, error)
}

// MQ 定义消息队列接口
type MQ interface {
	// NewProducer 创建生产者
	NewProducer() (Producer, error)
	// NewConsumer 创建消费者
	NewConsumer(cfg *ConsumerConfig) (Consumer, error)
	// Admin 获取管理接口
	Admin() Admin
	// Close 关闭消息队列连接
	Close() error
}

var (
	ErrTopicEmpty        = errors.New("topic cannot be empty")
	ErrNilMessage        = errors.New("message cannot be nil")
	ErrNotConnected      = errors.New("message queue not connected")
	ErrUnsubscribeFailed = errors.New("failed to unsubscribe from topic")
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrTimeout           = errors.New("poll timeout")
	// 新增错误类型
	ErrDuplicateMessage = errors.New("duplicate message")
	ErrMessageExpired   = errors.New("message expired")
	ErrDeliveryMode     = errors.New("unsupported delivery mode")
	ErrConsumerPaused   = errors.New("consumer paused")
	ErrInvalidBatchSize = errors.New("invalid batch size")
	ErrTopicNotFound    = errors.New("topic not found")
	ErrGroupNotFound    = errors.New("consumer group not found")
)
