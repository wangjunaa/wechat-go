package server

import (
	"demo/handle"
	"demo/tools/token"
	"github.com/gin-gonic/gin"
)

type membersJson struct {
	MembersId []string `json:"members"`
}

// CreateGroup
// @Summary 创建群组
// @Tags 群
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param members body membersJson true "初始群员id"
// @Success 200 {object} Model.GroupBasic "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /group/create [post]
func CreateGroup(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}

	temp := membersJson{}
	err := c.ShouldBindJSON(&temp)
	if err != nil {
		c.String(412, err.Error())
		return
	}

	group, err := handle.CreateGroup(id, temp.MembersId)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, group)
}

// DeleteGroup
// @Summary 删除群聊
// @Tags 群
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param groupId query string true "群id"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /group/delete [post]
func DeleteGroup(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}

	gid := c.Query("groupId")
	err := handle.CheckOwner(id, gid)
	if err != nil {
		c.String(412, err.Error())
		return
	}

	err = handle.DeleteGroup(gid)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.String(200, "OK")
}

// RemoveGroupMember
// @Summary 删除群员
// @Tags 群
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param groupId query string true "群id"
// @Param deletedId query string true "被删除用户id"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /group/removeMember [post]
func RemoveGroupMember(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}

	gid := c.Query("groupId")
	err := handle.CheckOwner(id, gid)
	if err != nil {
		c.String(412, err.Error())
		return
	}

	err = handle.CheckOwner(gid, id)
	if err != nil {
		c.String(412, err.Error())
		return
	}

	deletedId := c.Query("deletedId")
	err = handle.RemoveGroupMember(gid, deletedId)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, "OK")
}

// InviteToGroup
// @Summary 邀请入群
// @Tags 群
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param groupId query string true "群id"
// @Param invitedMembers body membersJson true "被邀请用户id"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /group/inviteMember [post]
func InviteToGroup(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}

	gid := c.Query("groupId")
	err := handle.CheckOwner(id, gid)
	if err != nil {
		c.String(412, err.Error())
		return
	}

	temp := membersJson{}
	err = c.ShouldBindJSON(&temp)
	if err != nil {
		c.String(412, err.Error())
		return
	}

	err = handle.InviteToGroup(gid, temp.MembersId)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.String(200, "OK")
}

// GetGroupMembers
// @Summary 获取群员信息
// @Tags 群
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param groupId query string true "群id"
// @Success 200 {array} Model.ShowUser "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /group/getMembers [get]
func GetGroupMembers(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}
	gid := c.Query("groupId")
	g, err := handle.GetGroupById(gid)
	if err != nil {
		c.String(500, "内部错误")
		return
	}
	c.JSON(200, handle.UserListToShow(g.Members))
}

// GetGroup
// @Summary 获取群信息
// @Tags 群
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param groupId query string true "群id"
// @Success 200 {object} Model.GroupBasic "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /group/getGroup [get]
func GetGroup(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}

	gid := c.Query("groupId")
	group, err := handle.GetGroupById(gid)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, group)
}

// EnterGroupReq
// @Summary 发送加群申请
// @Tags 群
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param groupId query string true "群id"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /group/enterReq [post]
func EnterGroupReq(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}

	gid := c.Query("groupId")
	err := handle.EnterGroupReq(gid, id)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, "成功")
}

// EnterGroupAgree
// @Summary 获取群信息
// @Tags 群
// @Param x-token header string true "用户令牌"
// @Param id query string true "用户id"
// @Param groupId query string true "群id"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /group/enterAgree [post]
func EnterGroupAgree(c *gin.Context) {
	t := c.Request.Header.Get("x-token")
	id := c.Query("id")
	if ok := token.CheckToken(t, id); !ok {
		c.String(403, "拒绝访问")
		return
	}

	gid := c.Query("groupId")
	err := handle.CheckOwner(gid, id)
	if err != nil {
		c.String(412, err.Error())
		return
	}

	err = handle.EnterGroupAgree(gid, id)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, "成功")
}
