package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"os"
	"wechat/docs"
	"wechat/middleware"
	"wechat/service"
)

func Router() *gin.Engine {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("router.Router:", err)
		}
	}()
	logFile, err := os.Create("gin.log")
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile)
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	tool := r.Group("/tool")
	{
		tool.POST("/getMessageJson", service.GetMessageJson)
	}
	user := r.Group("/user")
	{
		user.POST("/create", service.CreateUser)
		user.POST("/login", service.Login)
		user.GET("/exist", service.IsUserExist)

		user.Use(middleware.MwCheckToken)
		user.GET("/find", service.FindUser)
		user.POST("/delete", service.DeleteUser)
		user.POST("/update", service.UpdateUser)
	}
	r.Use(middleware.MwCheckToken)
	message := r.Group("/message")
	{
		message.GET("/getPush", service.GetPushMessage)
	}
	friend := r.Group("/friend")
	{
		friend.POST("/request", service.FriendReq)
		friend.POST("/agree", service.FriendAgree)
		friend.POST("/delete", service.DeleteFriend)
		friend.GET("/getList", service.GetFriendList)
	}
	group := r.Group("/group")
	{
		group.GET("/getMembers", service.GetGroupMembers)
		group.GET("/getGroup", service.GetGroup)
		group.POST("/enterReq", service.EnterGroupReq)
		group.POST("/enterAgree", service.EnterGroupAgree)
		group.POST("/create", service.CreateGroup)
		group.POST("/delete", service.DeleteGroup)
		group.POST("/removeMember", service.RemoveGroupMember)
		group.POST("/inviteMember", service.InviteToGroup)
	}
	return r
}
