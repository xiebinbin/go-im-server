package user

import (
	"imsdk/internal/common/model/user"
	"imsdk/internal/common/pkg/config"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"

	"github.com/gin-gonic/gin"
)

type getAuthRequest = user.GetAuthParams

func GetSysInfo(ctx *gin.Context) {
	staticUrl, err := config.GetStaticUrl()
	if err != nil {
		response.RespErr(ctx, errno.Add("get-static-url-err", errno.ParamsErr))
		return
	}
	data := map[string]string{
		"static_url": staticUrl,
		"pub_key":    user.GetPubKey(ctx),
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
