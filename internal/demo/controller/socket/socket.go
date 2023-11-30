package socket

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/demo/model/socket"
	"imsdk/pkg/response"
)

func SendCmd(ctx *gin.Context) {
	 err := socket.SendCmd(ctx)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
