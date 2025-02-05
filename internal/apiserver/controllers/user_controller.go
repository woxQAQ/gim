package controllers

// UserController 处理用户相关的HTTP请求
type UserController struct {
	// 在这里注入所需的service
}

// NewUserController 创建UserController实例
func NewUserController() *UserController {
	return &UserController{}
}
