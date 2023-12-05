package users

// users godoc
// users 实现了用户相关的api
// 包括登陆，注册，获取用户信息（有限的），注销用户

import (
	"gIM/internal/db"
	"gIM/internal/global"
	"gIM/internal/middleware/jwt"
	"gIM/internal/models"
	"net/http"
	"strconv"
	"time"

	vad "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	userId := uint(Id)
	// userName := ctx.PostForm("name")
	userPwd := ctx.PostForm("password")

	// 确认用户存在
	data, err := db.QueryById(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "登录失败",
			"error":   err,
		})
		return
	}

	if data.Name == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户不存在",
		})
		return
	}

	ok := checkPwd(userPwd, data.Salt, data.Password)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    -1,
			"message": "登陆失败：密码错误",
		})
		return
	}

	data.LoginTime = time.Now()
	data.Online = true
	// todo 此处的用户名密码都是通过明文传输的，不安全，如何进行密文传输？

	//zap.S().Info("鉴权")

	token, err := jwt.GenerateToken(data.Name, data.Password, "woxQAQ")
	if err != nil {
		zap.S().Info("生成token失败", err)
		return
	}
	// 鉴权成功

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登陆成功",
		"token":   token,
	})

	err = db.UpdateUser(data)
	if err != nil {
		zap.S().Info("更新数据失败")
		// todo 请求重传
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
	isExist := db.UserExist(Name)
	if isExist {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：用户已存在，请进行登陆",
			"data":    Name,
		})
		return
	}

	t := time.Now()
	user := models.UserBasic{
		Name: Name,
	}
	user.LogOutTime = t
	user.LoginTime = t
	user.HeartBeatTime = t

	salt := getSalt()
	user.Password = encryptPwd(Password, salt)
	user.Salt = salt

	err := db.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "注册成功，请进行登陆",
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
	user := models.UserBasic{}

	if err := ctx.Bind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "获取ID失败",
		})
		return
	}

	err := db.DeleteUser(user)
	if err != nil {
		zap.S().Info("删除用户失败", err)
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
	name := ctx.Param("name")

	users, err := db.QueryByUserName(name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "不存在此用户",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "查询成功",
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

// UpdateUser 用来更新重要的，用来标识用户的内容
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

	Name := ctx.PostForm("name")
	Password := ctx.PostForm("password")
	Email := ctx.PostForm("email")
	Phone := ctx.PostForm("phone")
	Gender := ctx.PostForm("gender")

	if Name != "" {
		user.Name = Name
	}
	if Password != "" {
		user.Password = Password
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
