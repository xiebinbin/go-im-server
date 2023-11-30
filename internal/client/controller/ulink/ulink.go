package ulink

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/ulink"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type Params = ulink.Params

func GetLinkInfo(ctx *gin.Context) {
	var params Params
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ResErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := ulink.GetULinkInfo(params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespData(ctx, data)
	return
}
