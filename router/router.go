package router

import (
	"demo/docs"
	"demo/server"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("router.Router:", err)
		}
	}()
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	tool := r.Group("/tool")
	{
		tool.POST("/getMessageJson", server.GetMessageJson)
	}
	user := r.Group("/user")
	{
		user.POST("/create", server.CreateUser)
		user.POST("/login", server.Login)

		user.Use(server.MwCheckToken)
		user.GET("/find", server.FindUser)
		user.POST("/delete", server.DeleteUser)
		user.POST("/update", server.UpdateUser)
	}
	r.Use(server.MwCheckToken)
	message := r.Group("/message")
	{
		message.GET("/getPush", server.GetPushMessage)
	}
	friend := r.Group("/friend")
	{
		friend.POST("/request", server.FriendReq)
		friend.POST("/agree", server.FriendAgree)
		friend.POST("/delete", server.DeleteFriend)
		friend.GET("/getList", server.GetFriendList)
	}
	group := r.Group("/group")
	{
		group.GET("/getMembers", server.GetGroupMembers)
		group.GET("/getGroup", server.GetGroup)
		group.POST("/enterReq", server.EnterGroupReq)
		group.POST("/enterAgree", server.EnterGroupAgree)
		group.POST("/create", server.CreateGroup)
		group.POST("/delete", server.DeleteGroup)
		group.POST("/removeMember", server.RemoveGroupMember)
		group.POST("/inviteMember", server.InviteToGroup)
	}
	return r
}
