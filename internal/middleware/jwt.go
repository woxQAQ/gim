package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrTokenExpired           = errors.New("令牌过期")
	ErrTokenExpiredMaxRefresh = errors.New("令牌已过最大刷新时间")
	ErrHeaderEmpty            = errors.New("需要认证")
	ErrHeaderFormat           = errors.New("认证头格式错误")
)

type Jwt struct {
	SignKey []byte
	//MaxRefresh time.Duration
}

type Claims struct {
	UserID        uint  `json:"userID"`
	ExpiredAtTime int64 `json:"expired_time"`
	jwt.RegisteredClaims
}

var key = []byte("qaq")

//var maxRefreash = 0

func NewJWT() *Jwt {
	return &Jwt{
		SignKey: key,
	}
}

func (j *Jwt) createToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SignKey)
}

func (j *Jwt) IssueToken(userID uint, iss string) (string, error) {
	now := time.Now()
	expireTime := now.Add(7 * 24 * time.Hour)
	claims := Claims{
		UserID:        userID,
		ExpiredAtTime: expireTime.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			ExpiresAt: jwt.NewNumericDate(expireTime),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token, err := j.createToken(claims)
	if err != nil {
		zap.S().Info(err)
		return "", err
	}
	return token, nil
}

func getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("auth")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrHeaderFormat
	}

	return parts[1], nil
}

func parseToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
}

func ParseToken(ctx *gin.Context) (*Claims, error) {
	tokenString, err := getTokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}
	token, err := parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func JWY() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		user := c.Query("userID")
		userId, err := strconv.Atoi(user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "userID不合法",
			})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "请登录",
			})
			c.Abort()
			return
		} else {
			claims, err := ParseToken(c)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    -1,
					"message": err,
				})
			}

			if claims.UserID != uint(userId) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    -1,
					"message": "登录不合法",
				})
			}

			zap.S().Info("token认证成功")
			c.Next()
		}
	}
}
