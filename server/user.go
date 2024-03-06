package server

import (
	"demo/handle"
	Model "demo/models"
	"demo/tools/token"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// FindUser
// @Summary 获取用户
// @Tags 用户操作
// @Param x-token query string true "用户令牌"
// @Param id path string true "用户id"
// @Success 200 {object} Model.ShowUser "成功"
// @Failure 412 {string} string "先决条件错误"
// @Router /user/find/{id} [get]
func FindUser(c *gin.Context) {
	id := c.Param("id")
	ts := c.Query("x-token")
	ok := token.CheckToken(ts, id)
	if !ok {
		c.String(http.StatusForbidden, "拒绝访问")
		return
	}
	u, err := handle.GetUser(id)
	if err != nil {
		c.String(http.StatusPreconditionFailed, err.Error())
		return
	}
	log.Println("serer.GetUser:", u)
	c.JSON(http.StatusOK, u.ToShowUser())
}

// CreateUser
// @Summary 创建用户
// @Tags 用户操作
// @Param json body Model.RegisterJson true "用户json"
// @Success 200 {string} string "添加成功"
// @Failure 412 {string} string "先决条件错误"
// @Router /user/create [post]
func CreateUser(c *gin.Context) {
	var json Model.RegisterJson
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.String(http.StatusPreconditionFailed, err.Error())
		return
	}

	if json.Password1 != json.Password2 {
		c.String(http.StatusPreconditionFailed, "两次密码不同")
		return
	}
	err = handle.CreateUser(json)
	if err != nil {
		c.String(http.StatusPreconditionFailed, err.Error())
		return
	}
	//fmt.Println("server.CreateUser:", id, password)
	c.String(http.StatusOK, "添加成功")
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户操作
// @Param x-token query string true "用户令牌"
// @Param id query string true "用户id"
// @Param password query string true "密码"
// @Success 200 {string} string "成功"
// @Failure 412 {string} string "先决条件错误"
// @Router /user/delete [post]
func DeleteUser(c *gin.Context) {
	id := c.Query("id")
	password := c.Query("password")
	ts := c.Query("x-token")
	ok := token.CheckToken(ts, id)
	if !ok {
		c.String(http.StatusForbidden, "拒绝访问")
		return
	}
	err := handle.DeleteUser(id, password)
	if err != nil {
		c.String(http.StatusPreconditionFailed, err.Error())
		return
	}
	c.String(http.StatusOK, "删除成功")
}

// UpdateUser
// @Summary 更新用户
// @Tags 用户操作
// @Param x-token header string true "用户令牌"
// @Param user body Model.UserBasic true "用户json信息"
// @Success 200 {string} string "成功"
// @Failure 412 {string} string "先决条件错误"
// @Failure 500 {string} string "内部错误"
// @Router /user/update [post]
func UpdateUser(c *gin.Context) {
	var userBasic Model.UserBasic
	err := c.ShouldBindJSON(&userBasic)
	if err != nil {
		log.Println("server.UpdateUser:", "绑定数据错误")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	ts := c.Request.Header.Get("x-token")
	log.Println(ts)
	ok := token.CheckToken(ts, userBasic.ID)
	if !ok {
		c.String(http.StatusForbidden, "拒绝访问")
		return
	}
	err = handle.UpdateUser(userBasic)
	if err != nil {
		log.Println("server.UpdateUser:", "更新数据错误")
		c.String(http.StatusPreconditionFailed, err.Error())
		return
	}
	c.String(http.StatusOK, "更新成功")
}

// Login
// @Summary 用户登录
// @Tags 用户操作
// @Param id header string true "用户id"
// @Param password header string true "密码"
// @Success 200 {string} string "成功"
// @Failure 412 {string} string "先决条件错误"
// @Router /user/login [post]
func Login(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("server.User.Login:", r)
		}
	}()
	id := c.Request.Header.Get("id")
	password := c.Request.Header.Get("password")
	tk, err := handle.Login(id, password)
	if err != nil {
		c.String(http.StatusPreconditionFailed, err.Error())
		return
	}
	c.String(http.StatusOK, tk)
}
