package gateway

import (
	"github.com/valyala/fasthttp"
	"github.com/woxQAQ/gim/internal/server/message"
	"time"
)

// authHandler 当数据到来时，对连接进行鉴权。当然，不是所有连接都需要鉴权
// 鉴权成功后，我们会将session保存在redis中，其中包括
// - sessionID
// - userId
// - token
// - conn
// - token销毁时间
func authHandler(token string, userId string) (time.Time, error) {
	// 编码数据
	//token := buf.GetToken()
	//userId := buf.GetUserId()

	client := &fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://127.0.0.1:11111/v1/auth" + "?" + "userId=" + userId)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	response := message.Response{}
	err := client.Do(req, resp)
	if err != nil {
		return time.Time{}, err
	}
	err = response.UnMarshal(resp.Body())
	if resp.StatusCode() != 200 {
		return time.Time{}, err
	}

	return response.Data_["expired_time"].(time.Time), nil
}
