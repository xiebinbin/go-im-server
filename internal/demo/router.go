package demo

import (
	"github.com/gin-gonic/gin"
	middlewares2 "imsdk/internal/common/middlewares"
	"imsdk/internal/demo/controller/bucket"
	"imsdk/internal/demo/controller/chat"
	"imsdk/internal/demo/controller/health"
	"imsdk/internal/demo/controller/login"
	"imsdk/internal/demo/controller/message"
	"imsdk/internal/demo/controller/socket"
	"imsdk/internal/demo/controller/user"
)

func GetEngine(engine *gin.Engine) {
	engine.Use(middlewares2.Cors)
	engine.GET("/healthCheck", health.AwsHealthCheck)
	engine.POST("/login", login.Login)
	engine.POST("/getAuth", login.GetAuth)
	//engine.Use(middlewares.CheckHeaders)
	//engine.Use(middlewares.Auth)
	engine.POST("/bucketInfo", bucket.GetBucketInfo)
	engine.POST("/createChat", chat.CreateChat)
	engine.POST("/createNoticeChat", chat.CreateNoticeChat)
	engine.POST("/joinChat", chat.JoinChat)
	engine.POST("/removeMember", chat.RemoveMember)
	engine.POST("/getChatList", chat.GetChatList)
	engine.POST("/getMemberByAChatIds", chat.GetMemberByAChatIds)

	//engine.POST("/sendMessage", message.SendTextMessage)
	engine.POST("/sendTextMessage", message.SendTextMessage)
	engine.POST("/sendImageMessage", message.SendImageMessage)
	engine.POST("/sendAttachmentMessage", message.SendAttachmentMessage)
	engine.POST("/sendCardAndTempMessage", message.SendCardAndTempMessage)

	engine.POST("/sendCardMessage", message.SendCardMessage)
	engine.POST("/sendVerticalCardMessage", message.SendVerticalCardMessage)
	engine.POST("/sendCustomizeMessage", message.SendCustomizeMessage)
	engine.POST("/sendMiddleMessage", message.SendMiddleMessage)
	engine.POST("/sendNotificationMessage", message.SendNotificationMessage)
	engine.POST("/setMessageDisable", message.SetMessageDisable)
	engine.POST("/getMessageInfo", message.GetMessageInfo)
	engine.POST("/getMessageList", message.GetMessageList)

	engine.POST("/getUserConnectInfo", user.GetConnectInfo)

	engine.POST("/sendCmd", socket.SendCmd)
}
