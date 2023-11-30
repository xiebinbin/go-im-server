package common

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/model/forward"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type SendCmdParams = forward.SendCmdParams

func SendCmd(ctx *gin.Context) {
	var params SendCmdParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := forward.SendCmdMsg(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
