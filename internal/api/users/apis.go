package users

// users godoc
// users 实现了用户相关的api
// 包括登陆，注册，获取用户信息（有限的），注销用户

import (
	"errors"
	"fmt"
	"github.com/woxQAQ/gim/config"
	"github.com/woxQAQ/gim/internal/db"
	"github.com/woxQAQ/gim/internal/middleware/jwt"
	"github.com/woxQAQ/gim/internal/models"
	"github.com/woxQAQ/gim/pkg/util"
	"net/http"
	"strconv"
	"time"

	vad "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type loginMsg struct {
	UserId string `json:"id"`
	// userName string
	UserPwd string `json:"password"`
}

type registerMsg struct {
	UserName   string `json:"name"`
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Gender     string `json:"gender"`
	Birthday   string `json:"birthday"`
}

// LoginById godoc
// @Summary 登陆
// @Description 使用用户名和密码进行登陆，并返回授权token
// @Accept json
// @Produce json
// @Success 200
// @Router /v1/user/login [post]
func LoginById(ctx *gin.Context) {
	var loginMessage loginMsg
	err := ctx.ShouldBindJSON(&loginMessage)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "传输格式错误",
			"error":   err,
		})
		return
	}
	Id, err := strconv.Atoi(loginMessage.UserId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "登录失败",
			"error":   err,
		})
		return
	}
	userId := uint(Id)
	userPwd := loginMessage.UserPwd

	// 确认用户存在
	data, err := db.QueryById(userId)
	if err != nil {
		zap.S().Error("查询用户失败:", err)

		statusCode := http.StatusInternalServerError
		msg := "其他错误"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = "用户不存在"
			statusCode = http.StatusBadRequest
		}
		ctx.AbortWithStatusJSON(statusCode, gin.H{
			"code":    -1,
			"message": msg,
			"error":   err.Error(),
		})
		return
	}

	ok := checkPwd(userPwd, data.Salt, data.Password)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "登陆失败：密码错误",
		})
		return
	}

	data.LoginTime = time.Now()
	data.Online = true
	// todo 此处的用户名密码都是通过明文传输的，不安全，如何进行密文传输？

	token, err := jwt.GenerateToken(data.Name, data.Password, "woxQAQ")
	if err != nil {
		zap.S().With(err).Info("生成token失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "生成token失败",
		})
		return
	}
	// 鉴权成功

	for i := 0; i < config.RecallTimes; i++ {
		err = db.UpdateUser(data)
		if err != nil {
			zap.S().Info("更新数据失败")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": fmt.Sprintf("更新数据失败,正在重传, 重传次数: %x", i+1),
				"error":   err.Error(),
			})
			continue
		}
		// 返回 token
		zap.S().Info("更新数据成功")
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "登陆成功",
			"token":   token,
		})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code":    -1,
		"message": "服务器更新数据失败！请重新登录",
		"error":   err.Error(),
	})
}

// Signup godoc
// @Summary 注册
// @Description 使用用户名和密码进行注册
// @Accept json
// @Produce json
// @Success 200
// @Router /vi/user/signup [post]
func Signup(ctx *gin.Context) {
	registerMessage := registerMsg{}
	err := ctx.ShouldBindJSON(&registerMessage)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "消息格式错误",
			"error":   err.Error(),
		})
		return
	}
	Name := registerMessage.UserName
	Email := registerMessage.Email
	Phone := registerMessage.Phone
	Gender := registerMessage.Gender
	Password := registerMessage.Password
	repassword := registerMessage.RePassword
	birthday := registerMessage.Birthday

	// 以下这部分内容实际上应该交给前端来做，但是为了保证前端和后端的兼容性，这里保留了一些简单的校验
	if vad.IsNotNull(Name) || vad.IsNotNull(Password) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：用户名或密码为空",
		})
		return
	}

	if repassword != Password {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：两次输入的密码不一致",
			"data":    Password,
		})
		return
	}
	if !vad.IsEmail(Email) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：邮箱格式错误",
			"data":    Email,
		})
		return
	}

	if !vad.Matches(Phone, "1^[3~9]{1}\\\\d{9}$") {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：手机号格式错误",
			"data":    Phone,
		})
		return
	}

	// todo 邮箱认证or 手机号认证
	if !vad.IsTime(birthday, config.DateTemp) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：日期格式错误",
			"data":    birthday,
		})
		return
	}
	if !vad.Matches(Gender, "1^(男|女)$") {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：性别格式错误",
			"data":    Gender,
		})
		return
	}
	// 用户名允许重复

	t := time.Now()
	user := models.UserBasic{
		Name:       Name,
		LogOutTime: t,
		LoginTime:  t,
	}

	salt, err := getSalt()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：生成salt失败",
			"err":     err,
		})
		return
	}
	user.Password = encryptPwd(Password, salt)
	user.Salt = salt

	if err := db.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "注册失败，请重试",
			"error":   err,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功，请进行登陆",
		"data":    user,
	})
}

func DelUser(ctx *gin.Context) {
	// todo 注销用户需要邮箱验证，手机号验证，还有密码验证，如何验证？

	user := models.UserBasic{}
	userid, err := util.ConvParamToUINT(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "获取id失败",
		})
		return
	}
	user.ID = userid

	if err := db.DeleteUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "删除用户失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注销账号成功",
	})
}

// InfoUser godoc
// @Summary 获取某个用户的有限信息
// @Description 使用用户名作为参数
// @Param name
// @Accept json
// @Produce json
// @Success 200
// @Router /vi/user/signup [post]
func InfoUser(ctx *gin.Context) {
	// 获取用户名
	id, err := util.ConvParamToUINT(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "获取id失败",
		})
		return
	}
	users, err := db.QueryById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "不存在此用户",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "查询成功",
		// 只返回部分信息
		// todo 权限管理
		"data": models.UserBasic{
			Name:       users.Name,
			Gender:     users.Gender,
			Phone:      users.Phone,
			Email:      users.Email,
			Avatar:     users.Avatar,
			LoginTime:  users.LoginTime,
			LogOutTime: users.LogOutTime,
			Online:     users.Online,
		},
	})
}

// UpdateUser godoc
// @Summary 更新用户信息
// @Description 更新用户信息，包括用户名，密码，电话，邮箱，性别等
// @Accept json
// @Produce json
// @Success 200
// @Router /vi/user/update [post]
// @Param id query int true "用户id"
// @Param name formData string true "用户名"
func UpdateUser(ctx *gin.Context) {

	id, err := util.ConvParamToUINT(ctx, "id")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取id失败",
		})
		return
	}
	user := models.UserBasic{}
	user.ID = id

	Name := ctx.Request.FormValue("name")
	Password := ctx.Request.FormValue("password")
	Email := ctx.Request.FormValue("email")
	Phone := ctx.Request.FormValue("phone")
	Gender := ctx.Request.FormValue("gender")

	if Name != "" {
		user.Name = Name
	}
	if Password != "" {
		salt, err := getSalt()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": "获取盐值失败",
			})
			return
		}
		user.Salt = salt
		user.Password = encryptPwd(Password, salt)
	}
	if Email != "" {
		user.Email = Email
	}
	if Phone != "" {
		user.Phone = Phone
	}
	if Gender != "" {
		user.Gender = Gender
	}

	_, err = vad.ValidateStruct(user)
	if err != nil {
		zap.S().Info("参数不合法", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数不匹配",
		})
		return
	}
	err = db.UpdateUser(user)
	if err != nil {
		zap.S().Info("更新用户失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "修改失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "修改成功",
	})
}
