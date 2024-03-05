package server

import (
	"demo/handle"
	Model "demo/models"
	"demo/tools/token"
	"github.com/gin-gonic/gin"
)

type friendListResp struct {
	Friends []Model.ShowUser `json:"friends"`
}

// FriendReq
// @Summary 发送好友请求
// @Tags 好友
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param id2 query string true "被请求的用户id"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /friend/request [post]
func FriendReq(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	id2 := c.Query("id2")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}
	if err := handle.FriendReq(id, id2); err != nil {
		c.String(500, err.Error())
		return
	}
	c.String(200, "OK")
}

// FriendAgree
// @Summary 同意好友请求
// @Tags 好友
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param id2 query string true "被同意的用户id"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /friend/agree [post]
func FriendAgree(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	id2 := c.Query("id2")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}
	if err := handle.FriendAgree(id, id2); err != nil {
		c.String(500, err.Error())
		return
	}
	c.String(200, "OK")
}

// GetFriendList
// @Summary 获取好友列表
// @Tags 好友
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Success 200 {object} friendListResp "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /friend/getList [get]
func GetFriendList(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}
	friendList, err := handle.GetFriendList(id)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, friendListResp{
		Friends: friendList,
	})
}

// DeleteFriend
// @Summary 删除好友
// @Tags 好友
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param deletedId query string true "被删除的用户id"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /friend/delete [post]
func DeleteFriend(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}

	deletedId := c.Query("deletedId")
	if err := handle.DeleteFriend(id, deletedId); err != nil {
		c.String(500, err.Error())
		return
	}
	c.String(200, "OK")
}
