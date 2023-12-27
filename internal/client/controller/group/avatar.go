package group

import (
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"imsdk/internal/client/model/group"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type UpdateAvatarRequest = group.UpdateAvatarRequest

func UpdateAvatar(ctx *gin.Context) {
	var params group.UpdateAvatarRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.UpdateAvatar(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func UpdateCover(ctx *gin.Context) {
	var params group.UpdateCoverRequest
	userId, _ := ctx.Get("uid")
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := group.UpdateCover(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
