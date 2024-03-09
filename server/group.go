package server

import (
	"demo/handler"
	"github.com/gin-gonic/gin"
)

// CreateGroup
// @Summary 创建群组
// @Tags 群
// @Param Authenticate header string true "用户令牌"
// @Param members formData []string true "初始群员id,应包括群主"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /group/create [post]
func CreateGroup(c *gin.Context) {
	var members []string
	members, ok := c.GetPostFormArray("members")
	if !ok {
		RespFailure(c, 400, "参数有误")
		return
	}

	id := c.GetString("id")
	group, err := handler.CreateGroup(id, members)
	if err != nil {
		RespFailure(c, 500, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", group, 1)
}

// DeleteGroup
// @Summary 删除群聊
// @Tags 群
// @Param Authenticate header string true "用户令牌"
// @Param groupId formData string true "群id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /group/delete [post]
func DeleteGroup(c *gin.Context) {
	id := c.GetString("id")
	gid := c.PostForm("groupId")
	if gid == "" {
		RespFailure(c, 400, paramError.Error())
		return
	}
	err := handler.CheckOwner(id, gid)
	if err != nil {
		RespFailure(c, 500, err.Error())
		return
	}

	err = handler.DeleteGroup(gid)
	if err != nil {
		RespFailure(c, 500, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", nil, 1)
}

// RemoveGroupMember
// @Summary 删除群员
// @Tags 群
// @Param Authenticate header string true "用户令牌"
// @Param deletedId formData string true "被删除用户id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /group/removeMember [post]
func RemoveGroupMember(c *gin.Context) {
	gid := c.PostForm("groupId")
	deletedId := c.PostForm("deletedId")
	id := c.GetString("id")
	err := handler.CheckOwner(id, gid)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	group, err := handler.RemoveGroupMember(gid, deletedId)
	if err != nil {
		RespFailure(c, 500, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", group, 1)
}

// InviteToGroup
// @Summary 邀请入群
// @Tags 群
// @Param Authenticate header string true "用户令牌"
// @Param groupId formData string true "群id"
// @Param invitedMembers formData []string true "被邀请用户id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /group/inviteMember [post]
func InviteToGroup(c *gin.Context) {
	id := c.GetString("id")
	gid := c.PostForm("groupId")
	err := handler.CheckOwner(id, gid)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}

	members, ok := c.GetPostFormArray("invitedMembers")
	if !ok {
		RespFailure(c, 400, "参数错误")
		return
	}

	group, err := handler.AddToGroup(gid, members)
	if err != nil {
		RespFailure(c, 500, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", group, 1)
}

// GetGroupMembers
// @Summary 获取群员信息
// @Tags 群
// @Param Authenticate header string true "用户令牌"
// @Param groupId formData string true "群id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /group/getMembers [get]
func GetGroupMembers(c *gin.Context) {
	gid := c.DefaultPostForm("groupId", "")
	g, err := handler.GetGroupById(gid)
	if err != nil {
		RespFailure(c, 500, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", handler.UserListToShow(g.Members), 1)
}

// GetGroup
// @Summary 获取群信息
// @Tags 群
// @Param Authenticate header string true "用户令牌"
// @Param groupId formData string true "群id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /group/getGroup [get]
func GetGroup(c *gin.Context) {
	gid := c.DefaultPostForm("groupId", "")
	group, err := handler.GetGroupById(gid)
	if err != nil {
		RespFailure(c, 500, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", group, 1)
}

// EnterGroupReq
// @Summary 发送加群申请
// @Tags 群
// @Param Authenticate header string true "用户令牌"
// @Param groupId formData string true "群id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /group/enterReq [post]
func EnterGroupReq(c *gin.Context) {
	id := c.GetString("id")
	gid := c.DefaultPostForm("groupId", "")
	err := handler.EnterGroupReq(gid, id)
	if err != nil {
		RespFailure(c, 500, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", nil, 1)
}

// EnterGroupAgree
// @Summary 同意群申请
// @Tags 群
// @Param Authenticate header string true "用户令牌"
// @Param groupId formData string true "群id"
// @Param agreedId formData string true "被同意id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /group/enterAgree [post]
func EnterGroupAgree(c *gin.Context) {
	gid := c.DefaultPostForm("groupId", "")
	agreedId := c.PostForm("agreedId")
	err := handler.EnterGroupAgree(gid, agreedId)
	if err != nil {
		RespFailure(c, 500, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", nil, 1)
}
