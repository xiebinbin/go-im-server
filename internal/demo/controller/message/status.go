package message

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/demo/model/message"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func SetMessageDisable(ctx *gin.Context) {
	var params message.SetDisableParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	_, err := message.SetMessageDisable(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
