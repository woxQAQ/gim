package mq

import (
	"errors"
	"sync"
	"time"
)

// messageItem 存储消息及其元数据
type messageItem struct {
	msg       *Message
	offset    int64
	timestamp time.Time
}

// ringBuffer 环形缓冲区实现
type ringBuffer struct {
	items []*messageItem
	mask  int64
	head  int64 // 生产者位置
	tail  int64 // 最老消息位置
	size  int64
	mutex sync.RWMutex
}

func newRingBuffer(size int64) *ringBuffer {
	return &ringBuffer{
		items: make([]*messageItem, size),
		mask:  size - 1,
		head:  0,
		tail:  0,
		size:  size,
	}
}

var ErrBufferFull = errors.New("buffer full")

func (rb *ringBuffer) put(msg *Message) error {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	// 使用位运算计算下一个位置
	nextHead := (rb.head + 1) & rb.mask
	if nextHead == rb.tail {
		return ErrBufferFull
	}

	rb.items[rb.head] = &messageItem{
		msg:       msg,
		offset:    rb.head,
		timestamp: time.Now(),
	}
	rb.head = nextHead
	return nil
}

func (rb *ringBuffer) get(offset int64) (*messageItem, bool) {
	rb.mutex.RLock()
	defer rb.mutex.RUnlock()

	// 简化边界检查，只检查是否在有效范围内
	if offset >= rb.head || offset < rb.tail {
		return nil, false
	}

	item := rb.items[offset&rb.mask]
	if item == nil {
		return nil, false
	}
	return item, true
}
