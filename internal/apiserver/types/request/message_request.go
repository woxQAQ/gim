package request

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	Type       int32  `json:"type"`
}

// GetMessageHistoryRequest 获取消息历史记录请求
type GetMessageHistoryRequest struct {
	UserID    string `json:"user_id"`
	PageSize  int    `json:"page_size"`
	PageToken string `json:"page_token,omitempty"`
}

// GetUnreadCountRequest 获取未读消息数量请求
type GetUnreadCountRequest struct {
	UserID string `json:"user_id"`
}
