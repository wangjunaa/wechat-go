package server

import (
	"github.com/gin-gonic/gin"
	"wechat/utils/token"
)

func MwCheckToken(c *gin.Context) {
	t := c.Request.Header.Get("Authenticate")
	id := token.CheckToken(t)
	if id == "" {
		RespFailure(c, 401, "身份验证失败")
		c.Abort()
		return
	}
	c.Set("id", id)
}
