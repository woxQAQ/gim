package client

import (
	"sync"

	"github.com/gorilla/websocket"

	"github.com/woxQAQ/gim/internal/types"
)

// Client 表示测试用的WebSocket客户端
type Client struct {
	conn      *websocket.Conn
	url       string
	userID    string
	platform  int32
	messages  []types.Message
	msgMutex  sync.Mutex
	closeOnce sync.Once
}

// New 创建一个新的测试客户端
func New(url, userID string, platform int32) *Client {
	return &Client{
		url:      url,
		userID:   userID,
		platform: platform,
	}
}

// Connect 连接到WebSocket服务器
func (c *Client) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return err
	}
	c.conn = conn

	// 启动消息接收协程
	go c.readMessages()
	return nil
}

// Close 关闭连接
func (c *Client) Close() {
	c.closeOnce.Do(func() {
		if c.conn != nil {
			c.conn.Close()
		}
	})
}

// readMessages 持续读取消息
func (c *Client) readMessages() {
	for {
		var msg types.Message
		if err := c.conn.ReadJSON(&msg); err != nil {
			return
		}

		c.msgMutex.Lock()
		c.messages = append(c.messages, msg)
		c.msgMutex.Unlock()
	}
}

// SendMessage 发送消息
func (c *Client) SendMessage(msg *types.Message) error {
	return c.conn.WriteJSON(msg)
}

// GetMessages 获取接收到的所有消息
func (c *Client) GetMessages() []types.Message {
	c.msgMutex.Lock()
	defer c.msgMutex.Unlock()
	return append([]types.Message{}, c.messages...)
}
