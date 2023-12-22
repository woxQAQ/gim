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

func main() {
	// 首先进行登陆
	token, err := login()
	if err != nil {
		panic(err)
	}
	imClient, err := gnet.NewClient(&gnet.BuiltinEventEngine{})
	if err != nil {
		panic(err)
	}

	// todo 发起websocket连接
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8088", Path: "/"}
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

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	req := message.NewRequest(message.ReqTemp, &message.RequestData{
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

	buf := make([]byte, 512)
	_, err = c.Read(buf)
	if err != nil {
		panic(err)
	}
	logging.Infof("message arrived: %s\n", string(buf))
}
