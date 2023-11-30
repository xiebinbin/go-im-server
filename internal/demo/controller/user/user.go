package user

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/demo/model/user"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)
type GetConnectInfoRequest = user.GetConnectInfoRequest

func GetConnectInfo(ctx *gin.Context) {
	var params GetConnectInfoRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ResErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := user.GetConnectInfo(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, data)
	return
}