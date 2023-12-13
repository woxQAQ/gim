package main

import (
	"bufio"
	"fmt"
	"github.com/panjf2000/gnet/pkg/logging"
	"log"
	"net"
	"net/url"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/ws"}

	c, err := net.Dial("tcp", u.Host)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func() {
		logging.Infof("connection stop: %s\n", c.LocalAddr().String())
		c.Close()
	}()
	rd := bufio.NewReader(c)
	msg, err := rd.ReadBytes('\n')

	if err != nil {
		log.Fatal("read:", err)
	}
	expectedMsg := "connection establishing...\n"
	if string(msg) != expectedMsg {
		logging.Fatalf("the first response packet mismatches, expect: %s, but got: %s", expectedMsg, msg)
	}
	fmt.Printf("the connect has get response: %s\n", msg)
	//count := 5
	//batch := 10
	//packetSize := 1024
	//for i := 0; i < count; i++ {
	//	batchSendAndRecv(c, rd, packetSize, batch)
	//}
}
