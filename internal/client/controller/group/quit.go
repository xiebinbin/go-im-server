package group

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/group"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func Quit(ctx *gin.Context) {
	var params group.QuitRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.QuitGroup(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func QuitAll(ctx *gin.Context) {
	var params group.QuitAllRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	err = group.QuitAllGroup(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
