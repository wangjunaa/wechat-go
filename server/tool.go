package server

import (
	Model "demo/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

// GetMessageJson
// @Summary 获取消息json格式
// @Tags 工具
// @Param message body Model.Message true "消息"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Failure 412 {string} string "先决条件错误"
// @Router /tool/getMessageJson [post]
func GetMessageJson(c *gin.Context) {
	msg := &Model.Message{}
	err := c.ShouldBindJSON(msg)
	if err != nil {
		c.String(412, err.Error())
		return
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		c.String(412, err.Error())
		return
	}
	c.String(200, string(bytes))
}
