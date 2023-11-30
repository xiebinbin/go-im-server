package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func GetAuthInfo(ctx *gin.Context) {
	var params user.GetUserInfoRequest
	uid, _ := ctx.Get(base.HeaderFieldUID)
	params.Ids = []string{uid.(string)}
	list, err := user.GetUserInfo(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.ResData(ctx, list)
	return
}

func GetListInfo(ctx *gin.Context) {
	var params user.GetUserInfoRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	list, err := user.GetUserInfo(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.ResData(ctx, list)
	return
}

func UpdateName(ctx *gin.Context) {
	var params user.UpdateNameRequest
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	err := user.UpdateName(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func UpdateAvatar(ctx *gin.Context) {
	var params user.UpdateAvatarRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := user.UpdateAvatar(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
