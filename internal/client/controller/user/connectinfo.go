package user

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user"
	user2 "imsdk/internal/common/model/user"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func GetConnectInfo(ctx *gin.Context) {
	var params user2.GetConnectInfoRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}

	list, err := user.GetConnectInfo(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, list)
	return
}
