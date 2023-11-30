package login

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/demo/model/user"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type loginRequest = user.LoginRequest
type getAuthRequest = user.GetAuthRequest

func Login(ctx *gin.Context) {
	var params loginRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ResErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := user.Login(params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespData(ctx, data)
	return
}

func GetAuth(ctx *gin.Context) {
	var params getAuthRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ResErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := user.GetAuth(params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespData(ctx, data)
	return
}
