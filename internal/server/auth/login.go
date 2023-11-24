package auth

import (
	"gIM/internal/db"
	"gIM/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Login(ctx *gin.Context) {
	userName := ctx.PostForm("name")
	userPwd := ctx.PostForm("password")

	// 确认用户存在
	_, err := db.QueryByUserName(userName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "登录失败：用户名错误",
		})
		return
	}

	// 检查密码
	data, err := db.QueryByNameAndPwd(userName, userPwd)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "登陆失败：密码错误",
		})
		return
	}

	// todo 此处的用户名密码都是通过明文传输的，不安全，如何进行密文传输？
	// 鉴权成功
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登陆成功",
		"userID":  data.ID,
	})
}

func Signup(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.PostForm("name")
	user.Password = ctx.PostForm("password")
	//repassword := ctx.PostForm("Identify")

	if user.Name == "" || user.Password == "" {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "注册失败：用户名或密码为空",
			"data":    user,
		})
		return
	}

	_, err := db.CreateUser(user)
	if err != nil {
		log.Fatalln(err)
	}

	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "注册成功，请进行登陆",
		"data":    user,
	})
}
