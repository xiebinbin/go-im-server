package group

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/group"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func AddAdministrators(ctx *gin.Context) {
	var params group.AddAdministratorRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	_, err = group.AddAdministrators(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespData(ctx, data)
	return
}

func RemoveAdministrators(ctx *gin.Context) {
	var params group.RemoveAdministratorRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.RemoveAdministrator(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
