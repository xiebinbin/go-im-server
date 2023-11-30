package chat

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/middlewares"
	"imsdk/internal/demo/model/chat"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func CreateChat(ctx *gin.Context) {
	var params chat.CreateParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := chat.CreateChat(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func CreateNoticeChat(ctx *gin.Context) {
	var params chat.CreateParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	params.Creator = middlewares.GetUId(ctx)
	err := chat.CreateNoticeChat(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func GetChatList(ctx *gin.Context) {
	data, _ := chat.GetChatList(ctx)
	response.RespListData(ctx, data)
	return
}

func GetMemberByAChatIds(ctx *gin.Context) {
	var params struct {
		AChatIDs []string `json:"achat_ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, _ := chat.GetMemberByAChatIds(ctx, params.AChatIDs)
	response.RespListData(ctx, data.Data.Items)
	return
}
