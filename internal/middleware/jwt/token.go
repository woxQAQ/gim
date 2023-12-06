package jwt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/woxQAQ/gim/internal/global"
	"go.uber.org/zap"
)

type Jwt struct {
	SignKey []byte
	//MaxRefresh time.Duration
}

type Claims struct {
	UserName      string `json:"userName"`
	Password      string `json:"password"`
	ExpiredAtTime int64  `json:"expired_time"`
	jwt.RegisteredClaims
}

var secret = []byte(global.JwtSecret)

//var maxRefreash = 0

func GenerateToken(userName string, password string, iss string) (string, error) {
	now := time.Now()
	expired := now.Add(3 * time.Hour)

	claims := Claims{
		UserName:      userName,
		Password:      password,
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

func JWY() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Query("token")
		if token == "" {
			context.JSON(http.StatusUnauthorized, gin.H{
				"message": "请输入登录时token",
			})
			context.Abort()
			return
		} else {
			claims, err := ParseToken(token)
			if err != nil {
				context.JSON(http.StatusUnauthorized, gin.H{
					"message": "token失效",
					"error":   err,
				})
				context.Abort()
				return
			} else if time.Now().Unix() > claims.ExpiredAtTime {
				context.JSON(http.StatusUnauthorized, gin.H{
					"now":     time.Now(),
					"claims":  claims,
					"message": "授权已过期",
				})
				context.Abort()
				return
			}
		}

		zap.S().Info("token认证")
		context.Next()
	}
}
