package server

import (
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type Client struct {
	UserId string `json:"userid"`
	w      sync.Mutex
	token  string
	conn   gnet.Conn
}

func NewClient(UserId string, token string, conn gnet.Conn) *Client {
	return &Client{
		UserId: UserId,
		token:  token,
		conn:   conn,
		w:      sync.Mutex{},
	}
}

func (c *Client) SetToken(token string) {
	c.w.Lock()
	defer c.w.Unlock()
	c.token = token
}

func (c *Client) GetToken() string {
	c.w.Lock()
	defer c.w.Unlock()
	return c.token
}

func (c *Client) SetConn(conn gnet.Conn) {
	c.w.Lock()
	defer c.w.Unlock()
	c.conn = conn
}

func (c *Client) GetConn() gnet.Conn {
	c.w.Lock()
	defer c.w.Unlock()
	return c.conn
}

func (c *Client) GetUserId() string {
	c.w.Lock()
	defer c.w.Unlock()
	return c.UserId
}

func (c *Client) SetUserId(userId string) {
	c.w.Lock()
	defer c.w.Unlock()
	c.UserId = userId
}

func (c *Client) Delete() {
	c.w.Lock()
	defer c.w.Unlock()
	c.UserId = ""
	c.token = ""
	c.conn = nil
}
