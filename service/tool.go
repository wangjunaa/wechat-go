package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"wechat/models"
)

// GetMessageJson
// @Summary 获取消息json格式
// @Tags 工具
// @Param message body models.Message true "消息"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /tool/getMessageJson [post]
func GetMessageJson(c *gin.Context) {
	msg := &models.Message{}
	err := c.ShouldBindJSON(msg)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		RespFailure(c, 400, err.Error())
		return
	}
	RespSuccess(c, 200, "成功", string(bytes), 1)
}
