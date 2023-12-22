package jwt

import (
	"github.com/woxQAQ/gim/config"
	"github.com/woxQAQ/gim/internal/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Jwt struct {
	SignKey []byte
	//MaxRefresh time.Duration
}

type Claims struct {
	UserId string `json:"user_id"`
	//Password      string `json:"password"`
	ExpiredAtTime int64 `json:"expired_time"`
	jwt.RegisteredClaims
}

var secret = []byte(config.JwtSecret)

//var maxRefreash = 0

func GenerateToken(userid string, iss string) (string, error) {
	now := time.Now()
	expired := now.Add(3 * time.Hour)

	claims := Claims{
		UserId: userid,
		//Password:      password,
		ExpiredAtTime: expired.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			ExpiresAt: jwt.NewNumericDate(expired),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func Auth(context *gin.Context) {
	token := context.GetHeader("Authorization")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "请在请求头中的 Authorization 中增加登录后返回的token",
		})
	} else {
		// 验证token
		userId := context.Query("userId")
		claims, err := ParseToken(token)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"data": gin.H{
					"message": "解析token时出错，token不合法",
				},
				"err": err.Error(),
			})
			return
		} else if time.Now().Unix() > claims.ExpiredAtTime {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"data": gin.H{
					"message": "授权已过期",
					"now":     time.Now(),
					"claims":  claims,
				},
				"err": errors.ErrAuthenticationFailed.Error(),
			})
			return
		} else if claims.UserId != userId {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"data": gin.H{
					"message": "授权失败",
				},
				"err": errors.ErrAuthenticationFailed.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"message": "授权成功",
			},
			"err": "",
		})
		zap.S().Info("token认证")
	}
}

func JWY() gin.HandlerFunc {
	return func(context *gin.Context) {
		Auth(context)
		context.Next()
	}
}
