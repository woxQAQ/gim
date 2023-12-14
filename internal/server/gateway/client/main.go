package main

import (
	"fmt"
	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/woxQAQ/gim/pkg/requests"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/"}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, err := net.Dial("tcp", u.Host)
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
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	buf := make([]byte, 512)

	// 发送认证消息
	req := requests.AuthenticateReq{
		UserName: "woxQAQ",
		Token:    "123456",
	}

	buf, err = req.Marshal()
	fmt.Println(buf)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Write(buf)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = c.Read(buf)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(buf)
}
