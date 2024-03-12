package server

import (
	"github.com/gin-gonic/gin"
	"wechat/handler"
)

// FriendReq
// @Summary 发送好友请求
// @Tags 好友
// @Param Authenticate header string true "用户令牌"
// @Param requestedId formData string true "被请求的用户id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /friend/request [post]
func FriendReq(c *gin.Context) {
	id := c.GetString("id")
	requestedId := c.DefaultPostForm("requestedId", "")
	if err := handler.FriendReq(id, requestedId); err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", nil, 1)
}

// FriendAgree
// @Summary 同意好友请求
// @Tags 好友
// @Param Authenticate header string true "用户令牌"
// @Param agreedId formData string true "被同意的用户id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /friend/agree [post]
func FriendAgree(c *gin.Context) {
	id := c.GetString("id")
	agreedId := c.DefaultPostForm("agreedId", "")
	if err := handler.FriendAgree(id, agreedId); err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", nil, 1)
}

// GetFriendList
// @Summary 获取好友列表
// @Tags 好友
// @Param Authenticate header string true "用户令牌"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /friend/getList [get]
func GetFriendList(c *gin.Context) {
	id := c.GetString("id")

	friendList, err := handler.GetFriendList(id)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", friendList, 1)
}

// DeleteFriend
// @Summary 删除好友
// @Tags 好友
// @Param Authenticate header string true "用户令牌"
// @Param deletedId formData string true "被删除的用户id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /friend/delete [post]
func DeleteFriend(c *gin.Context) {
	id := c.GetString("id")
	deletedId := c.DefaultPostForm("deletedId", "")
	if err := handler.DeleteFriend(id, deletedId); err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", nil, 1)
}
