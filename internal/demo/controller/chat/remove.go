package chat

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/demo/model/chat"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func RemoveMember(ctx *gin.Context) {
	var params chat.RemoveParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := chat.RemoveMember(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
