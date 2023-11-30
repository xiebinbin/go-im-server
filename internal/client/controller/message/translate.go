package message

import (
	"github.com/gin-gonic/gin"
	"html"
	"imsdk/internal/client/model/message"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type TranslateRequest = message.TranslateParams

func Translate(ctx *gin.Context) {
	var params TranslateRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	target, _ := ctx.Get("lang")
	data, err := message.Translate(params, target.(string))
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	res := map[string]interface{}{
		"res": html.UnescapeString(data),
		"mid": params.MId,
	}
	response.ResData(ctx, res)
	return
}
