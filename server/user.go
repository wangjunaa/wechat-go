package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"wechat/handler"
	Model "wechat/models"
)

// FindUser
// @Summary 获取用户
// @Tags 用户操作
// @Param Authenticate header string true "用户令牌"
// @Param queriedId query string true "被查询用户Id"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /user/find [get]
func FindUser(c *gin.Context) {
	queriedId := c.Query("queriedId")
	u, err := handler.FindUser(queriedId)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	//log.Println("serer.FindUser:", u)
	RespSuccess(c, 200, "成功", u.ToShowUser(), 1)
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户操作
// @Param Authenticate header string true "用户令牌"
// @Param password formData string true "密码"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /user/delete [post]
func DeleteUser(c *gin.Context) {
	id := c.GetString("id")
	password := c.PostForm("password")
	err := handler.DeleteUser(id, password)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", nil, 1)
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户操作
// @Param Authenticate header string true "用户令牌"
// @Param user body Model.UserBasic true "用户json信息"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /user/update [post]
func UpdateUser(c *gin.Context) {
	var userBasic Model.UserBasic
	err := c.ShouldBindJSON(&userBasic)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	id := c.GetString("id")
	if id != userBasic.ID {
		RespFailure(c, 401, "无法更新非本人信息")
		return
	}
	err = handler.UpdateUser(userBasic)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", nil, 1)
}

// Login
// @Summary 用户登录
// @Tags 用户操作
// @Param id formData string true "用户id"
// @Param password formData string true "密码"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /user/login [post]
func Login(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("server.User.Login:", r)
		}
	}()
	id := c.DefaultPostForm("id", "")
	password := c.DefaultPostForm("password", "")
	log.Println(id, password)
	tk, err := handler.Login(id, password)
	if err != nil {
		RespFailure(c, 400, "登录失败")
		return
	}
	RespSuccess(c, 200, "成功", tk, 1)
}

// CreateUser
// @Summary 创建用户
// @Tags 用户操作
// @Param phone formData string true "用户手机"
// @Param password formData string true "密码"
// @Param userName formData string true "用户名"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /user/create [post]
func CreateUser(c *gin.Context) {
	phone := c.DefaultPostForm("phone", "")
	password := c.DefaultPostForm("password", "")
	userName := c.DefaultPostForm("userName", "")
	if phone == "" || password == "" || userName == "" {
		RespFailure(c, 400, "参数错误")
	}
	t, err := handler.CreateUser(phone, userName, password)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	//fmt.Println("server.CreateUser:", id, password)
	RespSuccess(c, 200, "成功", t, 1)
}

// IsUserExist
// @Summary 判断是否存在用户
// @Tags 用户操作
// @Param phone query string true "用户手机"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /user/exist [get]
func IsUserExist(c *gin.Context) {
	phone := c.Query("phone")
	exist := handler.IsUserExist(phone)
	RespSuccess(c, 200, "成功", exist, 1)
}
