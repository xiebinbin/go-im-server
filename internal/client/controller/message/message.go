package message

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/message"
	"imsdk/pkg/errno"
	"imsdk/pkg/log"
	"imsdk/pkg/response"
)

func GetMessageStatus(ctx *gin.Context) {
	var params struct {
		Ids []string `json:"ids" binding:"required"`
	}
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := message.GetMessageStatus(uid, params.Ids)
	logCtx := log.WithFields(ctx, map[string]string{"action": "createChat"})
	if err != nil {
		log.Logger().Error(logCtx, "GetMessageStatus error -err : ", params, err)
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, data)
	return
}
