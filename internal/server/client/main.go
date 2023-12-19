package main

import (
	"fmt"
	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/panjf2000/gnet/v2"
	"github.com/woxQAQ/gim/internal/server/message"
	"log"
	"net/url"
	"os"
	"os/signal"
)

func authenticate(c gnet.Conn) (bool, error) {
	req := message.NewAuthReq(message.ReqAuthenticate, &message.AuthenticateData{
		Token:    "token",
		UserName: "username",
		Password: "password",
	})

	jsonData, err := req.Marshal()
	if err != nil {
		return false, err
	}

	_, err = c.Write(jsonData)
	if err != nil {
		return false, err
	}
	err := c.Flush()
	if err != nil {
		return false, err
	}

	logging.Infof("waiting for response...\n")

	_, err = c.Read(jsonData)
	if err != nil {
		return false, err
	}
}

func main() {
	imClient, err := gnet.NewClient(&gnet.BuiltinEventEngine{})
	if err != nil {
		panic(err)
	}

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/"}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//c, err := net.Dial("tcp", u.Host)
	c, err := imClient.Dial("tcp", u.Host)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func() {
		logging.Infof("connection stop: %s\n", c.LocalAddr().String())
		c.Close()
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		buf := make([]byte, 512)
		for {
			_, err := c.Read(buf)
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Printf("recv: %s\n", buf)
		}
	}()
	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()

	// 发送认证消息

}
