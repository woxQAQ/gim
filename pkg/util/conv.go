package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ConvParamToUINT(ctx *gin.Context, key string) (uint, error) {
	param := ctx.Param(key)
	return StringToUint(param)
}

func ConvFormToUINT(ctx *gin.Context, key string) (uint, error) {
	param := ctx.PostForm(key)
	return StringToUint(param)
}

func StringToUint(param string) (uint, error) {
	target, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return uint(target), nil
}
