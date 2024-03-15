package service

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	idNotExist = errors.New("用户id不存在")
	paramError = errors.New("参数错误")
)

type RespJson struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Data  interface{} `json:"data"`
	Count int         `json:"count"`
}

func RespSuccess(c *gin.Context, code int, msg interface{}, data interface{}, count int) {
	c.JSON(200, RespJson{
		Code:  code,
		Msg:   msg,
		Data:  data,
		Count: count,
	})
}

func RespFailure(c *gin.Context, code int, errMsg string) {
	c.JSON(code, RespJson{
		Code:  code,
		Msg:   errMsg,
		Data:  nil,
		Count: 1,
	})
}
