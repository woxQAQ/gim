package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/fasthttp"
	"github.com/woxQAQ/gim/internal/api/users"
	"github.com/woxQAQ/gim/internal/protobuf/proto_pb"
	"github.com/woxQAQ/gim/internal/server/message"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"time"
)

func login() (string, error) {
	req := users.LoginMsg{
		UserId:  "1234",
		UserPwd: "password",
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	request := fasthttp.AcquireRequest()
	request.SetRequestURI("http://127.0.0.1:11111/v1/auth/login")
	request.SetBody(jsonData)
	request.Header.SetMethod("POST")
	request.Header.Set("Accept", "application/json")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(request)

	if err := fasthttp.Do(request, resp); err != nil {
		return "", err
	}
	respData := &message.Response{}
	err = json.Unmarshal(resp.Body(), respData)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("status code: %d, error: %v", resp.StatusCode(), respData.Err)
	}

	token, ok := respData.Data_["token"].(string)
	if token == "" {
		return token, fmt.Errorf("there is not token field in the response body")
	}
	if !ok {
		return "", fmt.Errorf("type conv error")
	}
	return token, nil
}

type imClient struct {
	*gnet.BuiltinEventEngine
	wg sync.WaitGroup
}

func (ic *imClient) OnTraffic(c gnet.Conn) (action gnet.Action) {
	size := c.InboundBuffered()
	buf := make([]byte, size)
	_, err := c.Read(buf)
	if err != nil {
		return gnet.Close
	}
	fmt.Printf("receive data, %v\n", string(buf))
	ic.wg.Done()
	return
}

func main() {
	// 首先进行登陆
	token, err := login()
	if err != nil {
		logging.Errorf("loginging error: %v", err.Error())
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	if err != nil {
		panic(err)
	}

	//header := &http.Header{}
	// todo 发起websocket连接
	websocket.
		u := url.URL{Scheme: "ws", Host: "127.0.0.1:8088", Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func() {
		logging.Infof("connection stop: %s\n", conn.LocalAddr().String())
		c.Close()
	}()

	done := make(chan struct{})
	//req := message.NewRequest(message.ReqTestGatewayConn, &message.RequestData{
	//	"message": "hello, im",
	//}, token, "1234")
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				logging.Errorf("goroutine error: %v", err.Error())
				return
			}
			logging.Infof("recv: %v", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			var msgContent string
			_, err = fmt.Scan(msgContent)
			if err != nil {
				logging.Infof("panic: %v", err.Error())
				panic(err)
			}
			msg := &proto_pb.SingleMessage{
				Content:    msgContent,
				SenderId:   1,
				ReceiverId: 2,
				Timestamp:  time.Now().Unix(),
			}
			buf, err := proto.Marshal(msg)
			if err != nil {
				logging.Infof("panic: %v", err.Error())
				panic(err)
			}
			err = conn.WriteMessage(websocket.TextMessage, buf)
			if err != nil {
				logging.Infof("panic: %v", err.Error())
				panic(err)
			}
		case <-interrupt:
			logging.Infof("Exiting")
			err = conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logging.Infof("panic: %v", err.Error())
				panic(err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
