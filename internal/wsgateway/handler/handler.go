package handler

// Handler 定义消息处理器接口
type Handler interface {
	// Handle 处理消息，返回是否继续处理链
	Handle(msg []byte) (bool, error)
	// SetNext 设置下一个处理器
	SetNext(handler Handler)
	// GetNext 获取下一个处理器
	GetNext() Handler
}

// BaseHandler 处理器基础实现
type BaseHandler struct {
	next Handler
}

// SetNext 设置下一个处理器
func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

// GetNext 获取下一个处理器
func (h *BaseHandler) GetNext() Handler {
	return h.next
}
