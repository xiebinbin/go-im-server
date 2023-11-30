package chat

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/model/chat"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type IdRequest = chat.GetMembersParams
type GetChatRequest = chat.GetChatParams

type IdsRequest struct {
	ChatIDs []string `json:"ids" binding:"required"`
}

func GetIdList(ctx *gin.Context) {
	userId, _ := ctx.Get("uid")
	data, _ := chat.GetIdListByUid(userId.(string))
	data = append(data, base.OfficialPaymentChatId, base.OfficialNoticeChatId)
	response.RespListData(ctx, data)
	return
}

func GetMemberIds(ctx *gin.Context) {
	var params IdsRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, _ := chat.GetChatsMemberIds(params.ChatIDs)
	response.RespListData(ctx, data)
	return
}

func GetMembersInfo(ctx *gin.Context) {
	var params chat.MemberInfoParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, _ := chat.GetMembersInfo(params)
	response.RespListData(ctx, data)
	return
}

func DeleteChat(ctx *gin.Context) {
	var params DeleteChatRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	params.UId = ctx.Value("uid").(string)
	if err := chat.DeleteChat(ctx, params); err != nil {
		response.RespErr(ctx, errno.Add("fail", errno.DefErr))
		return
	}
	response.RespSuc(ctx)
	return
}
