package controllers

import (
	"github.com/go-fuego/fuego"
	"github.com/woxQAQ/gim/internal/apiserver/models"
	"github.com/woxQAQ/gim/internal/apiserver/services"
	"github.com/woxQAQ/gim/internal/apiserver/types/request"
)

// UserController 处理用户相关的HTTP请求
type UserController struct {
	userService *services.UserService
}

func (c *UserController) Route(sv *fuego.Server) {
	g := fuego.Group(sv, "/users", fuego.OptionDescription("用户相关接口"))
	fuego.Post(g, "/register", c.Register, fuego.OptionDescription("注册用户"))
	fuego.Post(g, "/login", c.Login, fuego.OptionDescription("用户登录"))
	fuego.Get(g, "/{userId}", c.GetUserInfo, fuego.OptionDescription("获取用户信息"))
}

// NewUserController 创建UserController实例
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register 处理用户注册请求
func (uc *UserController) Register(c fuego.ContextWithBody[request.RegisterRequest]) (*models.User, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	// 调用service层处理注册逻辑
	err = uc.userService.Register(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 处理用户登录请求
func (uc *UserController) Login(c fuego.ContextWithBody[request.LoginRequest]) (*models.User, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// 调用service层处理登录逻辑
	user, err := uc.userService.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserInfo 获取用户信息
func (uc *UserController) GetUserInfo(c fuego.ContextNoBody) (*models.User, error) {
	// 从上下文中获取用户ID
	userID := c.PathParam("userId")
	if userID == "" {
		return nil, fuego.BadRequestError{
			Title:  "Missing user ID",
			Detail: "User ID is required",
		}
	}

	// 调用service层获取用户信息
	user, err := uc.userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
