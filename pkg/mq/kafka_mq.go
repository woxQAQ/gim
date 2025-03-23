package mq

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

type KafkaConfig struct {
	Config  *Config
	Version string        // Kafka 版本
	Timeout time.Duration // 操作超时时间
}

type kafkaProducer struct {
	producer sarama.SyncProducer
}

type kafkaConsumer struct {
	consumer sarama.Consumer
	group    sarama.ConsumerGroup
	groupID  string
	topics   map[string]bool
}

type KafkaMQFactory struct {
	config *KafkaConfig
}

func NewKafkaMQFactory(config *KafkaConfig) *KafkaMQFactory {
	return &KafkaMQFactory{config: config}
}

func (f *KafkaMQFactory) NewProducer(cfg *Config) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	if f.config.Version != "" {
		version, err := sarama.ParseKafkaVersion(f.config.Version)
		if err != nil {
			return nil, err
		}
		config.Version = version
	}

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &kafkaProducer{producer: producer}, nil
}

func (f *KafkaMQFactory) NewConsumer(cfg *Config) (Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	if f.config.Version != "" {
		version, err := sarama.ParseKafkaVersion(f.config.Version)
		if err != nil {
			return nil, err
		}
		config.Version = version
	}

	consumer, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	group, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.Group, config)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		consumer: consumer,
		group:    group,
		groupID:  cfg.Group,
		topics:   make(map[string]bool),
	}, nil
}

func (p *kafkaProducer) Publish(ctx context.Context, msg *Message) error {
	if msg == nil {
		return ErrNilMessage
	}
	if msg.Topic == "" {
		return ErrTopicEmpty
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: sarama.ByteEncoder(msg.Value),
	}

	if msg.Key != "" {
		kafkaMsg.Key = sarama.StringEncoder(msg.Key)
	}

	_, _, err := p.producer.SendMessage(kafkaMsg)
	return err
}

func (p *kafkaProducer) Close() error {
	return p.producer.Close()
}

func (c *kafkaConsumer) Subscribe(ctx context.Context, topic string, handler func(*Message) error) error {
	if topic == "" {
		return ErrTopicEmpty
	}

	c.topics[topic] = true

	// 启动消费者组
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				partition, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
				if err != nil {
					continue
				}

				for msg := range partition.Messages() {
					message := &Message{
						Topic: msg.Topic,
						Value: msg.Value,
					}
					if msg.Key != nil {
						message.Key = string(msg.Key)
					}

					if err := handler(message); err != nil {
						// 处理错误，可以实现重试逻辑
					}
				}
			}
		}
	}()

	return nil
}

func (c *kafkaConsumer) Poll(ctx context.Context, topic string) (*Message, error) {
	if topic == "" {
		return nil, ErrTopicEmpty
	}

	partition, err := c.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}
	defer partition.Close()

	select {
	case msg := <-partition.Messages():
		return &Message{
			Topic: msg.Topic,
			Key:   string(msg.Key),
			Value: msg.Value,
		}, nil
	case <-ctx.Done():
		return nil, ErrTimeout
	}
}

func (c *kafkaConsumer) Unsubscribe(topic string) error {
	delete(c.topics, topic)
	return nil
}

func (c *kafkaConsumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		return err
	}
	return c.group.Close()
}
