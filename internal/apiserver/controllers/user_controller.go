package controllers

import (
	"time"

	"github.com/go-fuego/fuego"
	"github.com/golang-jwt/jwt/v5"

	"github.com/woxQAQ/gim/internal/apiserver/services"
	"github.com/woxQAQ/gim/internal/apiserver/types/request"
	"github.com/woxQAQ/gim/internal/apiserver/types/response"
	"github.com/woxQAQ/gim/internal/models"
)

// UserController 处理用户相关的HTTP请求
type UserController struct {
	userService *services.UserService
}

func (c *UserController) Route(sv *fuego.Server) {
	g := fuego.Group(sv, "/users",
		fuego.OptionDescription("用户相关接口"),
		fuego.OptionTags("user"),
	)
	fuego.Post(g, "/register", c.Register, fuego.OptionDescription("注册用户"))
	fuego.Post(g, "/login", c.Login, fuego.OptionDescription("用户登录"))
}

// NewUserController 创建UserController实例
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register 处理用户注册请求
func (uc *UserController) Register(c fuego.ContextWithBody[request.RegisterRequest]) (*response.UserResponse, error) {
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
	return user.ToResponse(), nil
}

// Login 处理用户登录请求
func (uc *UserController) Login(c fuego.ContextWithBody[request.LoginRequest]) (*response.LoginResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	token, err := uc.generateToken(req.UserId)
	if err != nil {
		return nil, fuego.InternalServerError{}
	}

	// 调用service层处理登录逻辑
	user, err := uc.userService.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (uc *UserController) generateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(""))
}
