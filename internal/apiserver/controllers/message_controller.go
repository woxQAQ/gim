package controllers

import (
	"github.com/go-fuego/fuego"
	"github.com/woxQAQ/gim/internal/apiserver/services"
	"github.com/woxQAQ/gim/internal/apiserver/types/request"
	"github.com/woxQAQ/gim/internal/apiserver/types/response"
)

// MessageController 处理消息相关的HTTP请求
type MessageController struct {
	messageService *services.MessageService
}

// NewMessageController 创建MessageController实例
func NewMessageController(messageService *services.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

func (c *MessageController) Route(sv *fuego.Server) {
	g := fuego.Group(sv, "/messages",
		fuego.OptionDescription("消息相关接口"),
		fuego.OptionTags("message"),
	)

	fuego.Get(g, "/history", c.GetMessageHistory, fuego.OptionDescription("获取消息历史记录"))
}

// GetMessageHistory 处理获取消息历史记录请求
func (c *MessageController) GetMessageHistory(ctx fuego.ContextWithBody[request.GetMessageHistoryRequest]) (*response.MessageHistoryResponse, error) {
	req, err := ctx.Body()
	if err != nil {
		return nil, err
	}

	// 调用service层获取消息历史记录
	return c.messageService.GetMessageHistory(req.UserID, req.PageSize, req.PageToken)
}
