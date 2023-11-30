package message

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/middlewares"
	"imsdk/internal/common/model/message"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type SetDisableParams = message.SetDisableParams

func SetDisable(ctx *gin.Context) {
	var params SetDisableParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	params.UId = middlewares.GetUId(ctx)
	err := message.SetDisable(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
	}
	response.RespSuc(ctx)
	return
}
