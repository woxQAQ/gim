package users

// users godoc
// users 实现了用户相关的api
// 包括登陆，注册，获取用户信息（有限的），注销用户

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	vad "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/woxQAQ/gim/internal/db"
	"github.com/woxQAQ/gim/internal/global"
	"github.com/woxQAQ/gim/internal/middleware/jwt"
	"github.com/woxQAQ/gim/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Login godoc
// @Summary 登陆
// @Description 使用用户名和密码进行登陆，并返回授权token
// @Accept json
// @Produce json
// @Success 200
// @Router /v1/user/login [post]
func Login(ctx *gin.Context) {
	Id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "登录失败",
			"error":   err,
		})
		return
	}
	userId := uint(Id)
	// userName := ctx.PostForm("name")
	userPwd := ctx.PostForm("password")

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

	for i := 0; i < 3; i++ {
		err = db.UpdateUser(data)
		if err != nil {
			zap.S().Info("更新数据失败")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    -1,
				"message": fmt.Sprintf("更新数据失败,正在重传, 重传次数: %x", i+1),
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
}

// Signup godoc
// @Summary 注册
// @Description 使用用户名和密码进行注册
// @Accept json
// @Produce json
// @Success 200
// @Router /vi/user/signup [post]
func Signup(ctx *gin.Context) {
	Name := ctx.PostForm("name")
	Password := ctx.PostForm("password")
	repassword := ctx.PostForm("repassword")
	birthday := ctx.PostForm("birthday")

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
	// todo 邮箱认证or 手机号认证
	if !vad.IsTime(birthday, global.DateTemp) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：日期格式错误",
			"data":    birthday,
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
	userid, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "获取id失败",
		})
		return
	}
	user.ID = uint(userid)

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
	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "获取id失败",
		})
		return
	}
	uid := uint(id)
	users, err := db.QueryById(uid)
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

	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取id失败",
		})
		return
	}
	user := models.UserBasic{}
	user.ID = uint(id)

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
