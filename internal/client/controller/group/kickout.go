package group

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/group"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func KickOut(ctx *gin.Context) {
	var params group.KickOutRequest
	uid := ctx.Value(base.HeaderFieldUID).(string)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := group.KickOutGroup(ctx, uid, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
