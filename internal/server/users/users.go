package users

import (
	"gIM/internal/db"
	"gIM/internal/middleware/jwt"
	"gIM/internal/models"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Login godoc
// @Summary 登陆
// @Description 使用用户名和密码进行登陆，并返回授权token
// @Accept json
// @Produce json
// @Success 200
// @Router /user/login [post]
func Login(ctx *gin.Context) {
	userName := ctx.PostForm("name")
	userPwd := ctx.PostForm("password")

	// 确认用户存在
	_, err := db.QueryByUserName(userName)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    -1,
			"message": "登录失败：用户名错误",
		})
		return
	}

	// 检查密码
	data, err := db.QueryByNameAndPwd(userName, userPwd)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    -1,
			"message": "登陆失败：密码错误",
		})
		return
	}

	//if !data.IsLogOut {
	//	ctx.JSON(http.StatusForbidden, gin.H{
	//		"code":    -1,
	//		"message": "登陆失败：用户已登陆",
	//	})
	//	return
	//}

	data.LoginTime = time.Now()
	data.IsLogOut = false
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
		"userID":  data.ID,
		"token":   token,
	})

	_, err = db.UpdateUser(*data)
	if err != nil {
		zap.S().Info("更新数据失败")
		// todo 请求重传
		return
	}

}

func Signup(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.PostForm("name")
	user.Password = ctx.PostForm("password")
	repassword := ctx.PostForm("Identify")

	if user.Name == "" || user.Password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：用户名或密码为空",
		})
		return
	}

	if repassword != user.Password {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：两次输入的密码不一致",
			"data":    user.Password,
		})
		return
	}

	isExist := db.UserExist(user.Name)
	if isExist {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "注册失败：用户已存在，请进行登陆",
			"data":    user.Name,
		})
		return
	}

	t := time.Now()
	user.LogOutTime = t
	user.LoginTime = t
	user.HeartBeatTime = t
	// todo 加密

	_, err := db.CreateUser(user)
	if err != nil {
		log.Fatalln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功，请进行登陆",
		"data":    user,
	})
}

func DelUser(ctx *gin.Context) {
	user := models.UserBasic{}

	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info("获取ID失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "删除用户失败",
		})
		return
	}

	user.ID = uint(id)
	err = db.DeleteUser(user)
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

func InfoUser(ctx *gin.Context) {
	// 获取用户名
	name := ctx.Param("name")

	// todo 鉴权

	users, err := db.QueryByUserName(name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "不存在此用户",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.UserBasic{
		Name:       users.Name,
		Gender:     users.Gender,
		Phone:      users.Phone,
		Email:      users.Email,
		Avatar:     users.Avatar,
		LoginTime:  users.LoginTime,
		LogOutTime: users.LogOutTime,
		IsLogOut:   users.IsLogOut,
		DeviceInfo: users.DeviceInfo,
	})
}

// UpdateUser 用来更新重要的，用来标识用户的内容
func UpdateUser(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.PostForm("id"))
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

	// TOdo 鉴权
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		zap.S().Info("参数不合法", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数不匹配",
		})
		return
	}
	_, err = db.UpdateUser(user)
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
