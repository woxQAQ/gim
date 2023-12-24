package main

import (
	"encoding/json"
	"fmt"
	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/fasthttp"
	"github.com/woxQAQ/gim/internal/api/users"
	"github.com/woxQAQ/gim/internal/server/message"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
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
		panic(err)
	}
	ic := &imClient{wg: sync.WaitGroup{}}
	client, err := gnet.NewClient(ic)
	if err != nil {
		panic(err)
	}

	// todo 发起websocket连接
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8088", Path: "/"}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//c, err := net.Dial("tcp", u.Host)
	c, err := client.Dial("tcp", u.Host)
	err = client.Start()
	if err != nil {
		return
	}
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func() {
		logging.Infof("connection stop: %s\n", c.LocalAddr().String())
		c.Close()
	}()

	req := message.NewRequest(message.ReqTestGatewayConn, &message.RequestData{
		"message": "hello, im",
	}, token, "1234")
	jsonData, err := req.MarshalJSON()
	if err != nil {
		panic(err)
	}
	_, err = c.Write(jsonData)
	if err != nil {
		panic(err)
	}
	ic.wg.Add(2)
	ic.wg.Wait()
}
