package user

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/model/user"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type getAuthRequest = user.GetAuthParams

func GetPubKey(ctx *gin.Context) {
	data := map[string]string{
		"pub_key": user.GetPubKey(ctx),
	}
	response.RespDataWithNoEnc(ctx, data)
	return
}

func Register(ctx *gin.Context) {
	data, err := user.RegUser(ctx)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespData(ctx, data)
	return
}

func IsRegister(ctx *gin.Context) {
	data, err := user.IsRegister(ctx)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	isReg := false
	if data.ID != "" {
		isReg = true
	}
	res := map[string]bool{
		"is_register": isReg,
	}
	response.RespData(ctx, res)
	return
}

func Login(ctx *gin.Context) {
	var params getAuthRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	response.RespSuc(ctx)
	return
}
