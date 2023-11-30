package message

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/demo/model/message"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func GetMessageInfo(ctx *gin.Context) {
	var params message.GetMessageInfoParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := message.GetMessageInfo(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, data.Data.Items)
	return
}


func GetMessageList(ctx *gin.Context) {
	var params message.GetMessageListParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := message.GetMessageList(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, data.Data.Items)
	return
}
