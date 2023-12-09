package friends

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/woxQAQ/gim/internal/db"
	"go.uber.org/zap"
	"net/http"
	"sort"
	"strconv"
)

// SendFriendRequest
// @Router /friends/{friend_id} [post]
// @Summary 发送等待对方处理的好友请求
// @Description 在Relation表中创建一对关系，分别对应请求发起方
// 和请求接收方的。userId为请求发起方，friendId为请求接收方。
// @Accept form-data
// @Produce json
// @Param id query number true "User ID"
// @Param id path number true "Friend ID"
// @Success 200 {object} gin.H
// @Failure 400,500 {object} gin.H
func SendFriendRequest(ctx *gin.Context) {
	uid := ctx.Param("id")
	if !govalidator.StringMatches(uid, `^[0-9]+$`) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "User ID is invalid",
		})
		return
	}
	fid := ctx.PostForm("id")
	if !govalidator.StringMatches(uid, `^[0-9]+$`) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "Friend ID is invalid",
		})
		return
	}
	if uid == fid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "You Cannot add yourself as friend",
			"data":    fid,
		})
		return
	}
	friends := ctx.PostFormArray("friends")

	if sort.Search(len(friends), func(i int) bool { return friends[i] >= fid }) < len(friends) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "好友已经存在",
		})
		return
	}

	userId, err := strconv.Atoi(uid)
	if err != nil {
		zap.S().Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "Internal Server Error:" + err.Error(),
		})
		return
	}
	friendId, err := strconv.Atoi(fid)
	if err != nil {
		zap.S().Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "Internal Server Error:" + err.Error(),
		})
		return
	}

	relation, err := db.CreateRelation(uint(userId), uint(friendId))
	if err != nil {
		zap.S().Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "create relation error: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    relation,
	})
	return
}

func GetFriendList(ctx *gin.Context) {

}

func AllowFriend(ctx *gin.Context) {

}
