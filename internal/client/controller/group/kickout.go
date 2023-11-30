package group

import (
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"imsdk/internal/client/model/group"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func KickOut(ctx *gin.Context) {
	var params group.KickOutRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.KickOutGroup(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
