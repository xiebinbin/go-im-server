package chat

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/demo/model/chat"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func JoinChat(ctx *gin.Context) {
	var params chat.JoinParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := chat.JoinChat(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
