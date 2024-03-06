package server

import (
	"demo/handle"
	"demo/tools/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetPushMessage
// @Summary 创建通知通道
// @Tags 消息操作
// @Param x-token query string true "用户令牌"
// @Param id query string true "用户id"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "内部错误"
// @Failure 403 {string} string "拒绝访问"
// @Router /message/getPush [get]
func GetPushMessage(c *gin.Context) {
	id := c.Query("id")
	ts := c.Query("x-token")
	ok := token.CheckToken(ts, id)
	if !ok {
		c.String(http.StatusForbidden, "拒绝访问")
		return
	}
	handle.WsHandler(id, c)
}
