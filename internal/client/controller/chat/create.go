package chat

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/model/chat"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func CreateChat(ctx *gin.Context) {
	var params chat.CreateParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	_, err := chat.Create(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func GetMyChat(ctx *gin.Context) {
	var params chat.GetMyChatParams
	params.UId = ctx.Value("uid").(string)
	data, err := chat.GetMyChat(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, data)
	return
}

func DeleteMyChat(ctx *gin.Context) {
	var params chat.DeleteChatRequest
	params.UId = ctx.Value("uid").(string)
	err := chat.DeleteChat(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
