package client

import (
	"imsdk/internal/client/controller/bucket"
	"imsdk/internal/client/controller/chat"
	"imsdk/internal/client/controller/contacts"
	"imsdk/internal/client/controller/copywriting"
	"imsdk/internal/client/controller/group"
	"imsdk/internal/client/controller/health"
	"imsdk/internal/client/controller/message"
	"imsdk/internal/client/controller/test"
	"imsdk/internal/client/controller/user"
	"imsdk/internal/client/controller/ws"
	"imsdk/internal/client/middlewares"
	middlewares2 "imsdk/internal/common/middlewares"

	"github.com/gin-gonic/gin"
)

func GetEngine(engine *gin.Engine) {
	engine.POST("/pushMessageToClient", ws.PushMessageToClient)
	engine.GET("/", ws.Connect)
	engine.GET("/wsConnect", ws.Connect)
	engine.POST("/tt", test.Redis)
	engine.POST("/test/:action", test.GetUrlParams)
	engine.Use(middlewares2.Cors)
	engine.GET("/healthCheck", health.AwsHealthCheck)
	engine.POST("/sys/info", user.GetSysInfo)
	engine.Use(middlewares.CheckHeaders)
	engine.POST("/auth/register", user.Register)
	engine.POST("/user/isRegister", user.IsRegister)
	engine.POST("/auth/getQRCode", user.GenerateLoginCode)
	engine.POST("/auth/getQRCodeScanRes", user.ScanQrCodeRes)
	engine.Use(middlewares.Auth)
	engine.POST("/sys/preSignURL", bucket.GetPreSignURL)
	engine.POST("/auth/info", user.GetAuthInfo)
	engine.POST("/user/updateName", user.UpdateName)
	engine.POST("/user/updateAvatar", user.UpdateAvatar)
	engine.POST("/user/updateGender", user.UpdateGender)
	engine.POST("/user/updateSign", user.UpdateSign)
	engine.POST("/user/getBatchInfo", user.GetListInfo)
	engine.POST("/user/unsubscribe", user.Unsubscribe)

	engine.POST("/auth/appScanLoginQrCode", user.AppScanLoginQrCode)
	engine.POST("/auth/appConfirmLogin", user.AppConfirmLogin)

	engine.POST("/chat/create", chat.CreateChat)
	engine.POST("/chat/list", chat.GetMyChat)
	engine.POST("/chat/delete", chat.DeleteMyChat)

	engine.POST("/message/send", message.SendMessage)
	engine.POST("/message/list", message.GetMessageByChatId)
	engine.POST("/message/deleteBatch", message.DeleteBatch)
	engine.POST("/message/deleteSelfAll", message.DeleteAllByUID)
	engine.POST("/message/deleteBatchByChatIds", message.DeleteBatchByChatIds)
	engine.POST("/message/deleteSelfByChatIds", message.DeleteSelfByChatIds)
	engine.POST("/message/revokeBatch", message.RevokeBatch)
	engine.POST("/message/revokeByChatIds", message.RevokeByChatIds)

	engine.POST("/friend/inviteApply", contacts.Apply)
	engine.POST("/friend/inviteList", contacts.GetList)
	engine.POST("/friend/inviteAgree", contacts.Agree)
	engine.POST("/friend/inviteReject", contacts.Refuse)
	engine.POST("/friend/inviteRead", contacts.ReadApplyFriends)
	engine.POST("/friend/list", contacts.GetFriendList)
	engine.POST("/friend/updateRemark", contacts.UpdateRemark)
	engine.POST("/friend/relationList", contacts.GetRelationList)
	engine.POST("/friend/deleteUnilateral", contacts.DelFriendsUnilateral) //单向删除好友
	engine.POST("/friend/deleteBilateral", contacts.DelFriendsBilateral)   //双向删除

	engine.POST("/member/info", user.GetListInfo)

	engine.POST("/group/create", group.CreateGroup)
	engine.POST("/group/join", group.Join)
	engine.POST("/group/inviteJoin", group.InviteJoin)
	engine.POST("/group/adminAgree", group.AdminAgreeJoin)
	engine.POST("/group/agreeInvite", group.AgreeInvite)
	// 拒绝邀请
	engine.POST("/group/rejectJoin", group.RejectJoin)
	engine.POST("/group/quit", group.Quit)
	engine.POST("/group/quitByIds", group.QuitByIds)
	engine.POST("/group/quitAll", group.QuitAll)
	engine.POST("/group/kickOut", group.KickOut)
	engine.POST("/group/members", group.GetMembersInfo)
	engine.POST("/group/membersByIds", group.GetMembersByIds)
	engine.POST("/group/encInfoByIds", group.GetEncInfoByIds)
	engine.POST("/group/list", group.GetList)
	engine.POST("/group/updateName", group.UpdateName)
	engine.POST("/group/updateAlias", group.UpdateAlias)
	engine.POST("/group/updateAvatar", group.UpdateAvatar)
	engine.POST("/group/updateCover", group.UpdateCover)
	engine.POST("/group/updateNotice", group.UpdateNotice)
	engine.POST("/group/getNotice", group.GetNotice)
	engine.POST("/group/updateDesc", group.UpdateNotice)
	engine.POST("/group/getDesc", group.GetDesc)
	engine.POST("/group/dismiss", group.Disband)
	engine.POST("/group/transfer", group.Transfer)
	engine.POST("/group/addAdmin", group.AddAdministrators)
	engine.POST("/group/removeAdmin", group.RemoveAdministrators)
	engine.POST("/group/clearMessage", group.ClearMessage)
	engine.POST("/group/applyList", group.ApplyList)
	engine.POST("/group/myApplyList", group.MyApplyList)
	engine.POST("/group/detailByIds", group.DetailByIds)
	engine.POST("/group/addApp", group.AddApp)
	engine.POST("/group/appInfo", group.AppInfo)
	engine.POST("/group/appList", group.AppList)
	engine.POST("/group/appUpdate", group.AppUpdate)
	engine.POST("/group/appDelete", group.AppDelete)
	engine.POST("/group/appDeleteByGIds", group.DeleteAppByGIds)

	//engine.POST("/sendCmd", common.SendCmd)
	engine.POST("/getChatMemberIds", chat.GetMemberIds)
	engine.POST("/getChatsMemberIds", chat.GetMemberIds)
	engine.POST("/getChatMemberInfoList", chat.GetMembersInfo)
	engine.POST("/updateChatIsTop", chat.UpdateChatIsTop)
	engine.POST("/updateChatIsTopV2", chat.UpdateChatIsTopV2)
	engine.POST("/updateChatIsMuteNotify", chat.UpdateChatIsMuteNotify)
	engine.POST("/getChatSetting", chat.GetChatSetting)
	engine.POST("/getChatSettings", chat.GetChatSettings)
	engine.POST("/hideChat", chat.HideChat)
	engine.POST("/deleteChat", chat.DeleteChat)
	engine.POST("/chatActive", chat.ReportChatActive)
	engine.POST("/chatChangeLogs", chat.ChangeLogs)

	//engine.POST("/sendMessageAsync", message.SendMessageAsync)
	engine.POST("/getMessageIds", message.GetMessageIds)
	engine.POST("/getMessageIdsV2", message.GetMessageIdsV2)
	engine.POST("/delMessage", message.DelSelf)
	engine.POST("/delMessageByChatIds", message.DelSelfByChatIds)
	engine.POST("/revokeMessage", message.Revoke)
	engine.POST("/revokeMessageByChatIds", message.RevokeByChatIds)
	engine.POST("/clearMessageByChatIds", message.ClearByChatIds)
	engine.POST("/hollowManSendMessage", message.HollowManSendMessage)
	engine.POST("/updateUnreadStock", message.UpdateUnreadStock)
	engine.POST("/translate", message.Translate)
	engine.POST("/userMaxSequence", message.GetUserMaxSequence)
	engine.POST("/getMessageStatus", message.GetMessageStatus)
	engine.POST("/setMessageDisable", message.SetDisable)

	engine.POST("/copywritingList", copywriting.List)
	engine.POST("/updateDeviceStatus", user.UpdateDeviceStatus)

	engine.POST("/getUserConnectInfo", user.GetConnectInfo)
	// pc login

}
