package start

import (
	"github.com/gin-gonic/gin"
	"imsdk/pkg/app"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type Resp struct {
	InviteUrl string `toml:"invite_url" json:"invite_url"`
}

func Start(ctx *gin.Context) {
	var data Resp
	err := app.Config().Bind("global", "start", &data)
	if err != nil {
		response.RespErr(ctx, errno.Add("failed", errno.SysErr))
		return
	}

	response.RespData(ctx, data)
	return
}
