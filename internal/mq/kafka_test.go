package mq

import (
	"github.com/segmentio/kafka-go"
	"testing"
)

func TestKafka(t *testing.T) {
	addr := "192.168.31.172:9092"
	conn, err := InitKafka(addr)
	defer func() {
		err := conn.Close()
		if err != nil {
			t.Error(err)
		}
	}()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("hello1")},
		kafka.Message{Value: []byte("hello2")},
		kafka.Message{Value: []byte("hello3")},
	)
	if err != nil {
		t.Error(err.Error())
	}
}
