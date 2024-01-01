package mq

import (
	"github.com/segmentio/kafka-go"
	"github.com/woxQAQ/gim/config"
	"gopkg.in/yaml.v3"
	"os"
)

type kafkaManager struct {
	conn *kafka.Conn
}

type kafkaConfig struct {
	Address string   `yaml:"kafka_address"`
	Topics  []string `yaml:"topics"`
}

func InitKafka() (*kafka.Conn, error) {

	data, err := os.ReadFile(config.KafkaFilePath)
	if err != nil {
		return nil, err
	}
	var configs kafkaConfig
	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	conn, err := kafka.Dial("tcp", configs.Address)
	if err != nil {
		return nil, err
	}
	for _, topic := range configs.Topics {
		err := conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}
