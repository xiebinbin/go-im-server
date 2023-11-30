package chat

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/model/chat"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func ChangeLogs(ctx *gin.Context) {
	var params chat.ChangeLogsReq
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}

	data, err := chat.ChangeLogsList(ctx, params)
	if err == nil {
		response.RespListData(ctx, data)
		return
	}
	response.RespErr(ctx, err)
	return
}
