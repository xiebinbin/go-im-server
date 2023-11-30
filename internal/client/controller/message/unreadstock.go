package message

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/message"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type UnreadStockRequest = message.UnreadStockParams

func UpdateUnreadStock(ctx *gin.Context) {
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	var params UnreadStockRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}

	if _, err := message.UpdateUnreadStockInfo(uid, params); err != nil {
		response.RespErr(ctx, err)
	} else {
		response.RespSuc(ctx)
	}
	return
}
