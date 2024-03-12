package server

import (
	"github.com/gin-gonic/gin"
	"wechat/handler"
)

// GetPushMessage
// @Summary 创建通知通道
// @Tags 消息操作
// @Param Authenticate header string true "用户令牌"
// @Success 200 {object} RespJson "成功"
// @Failure 401 {object} RespJson "验证失败"
// @Failure 400 {object} RespJson "参数有误"
// @Failure 500 {object} RespJson "内部错误"
// @Router /message/getPush [get]
func GetPushMessage(c *gin.Context) {
	id := c.GetString("id")
	handler.WsHandler(id, c)
}
