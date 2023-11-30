package middlewares

import (
	"github.com/gin-gonic/gin"
	"imsdk/pkg/encrypt"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func DecryptParams(ctx *gin.Context) {
	method := ctx.Request.Method
	var oriData string
	if method == "GET" {
		oriData = ctx.Query("data")
	} else {
		oriData = ctx.Request.FormValue("data")
	}
	if oriData == "" || len(oriData) <= 32 {
		ctx.Abort()
		response.ResErr(ctx, errno.Add("wrong-req", errno.WrongReq))
		return
	}
	encKey, _ := ctx.Get("enc_key")
	params := encrypt.AesCbcDecrypt(oriData, encKey.(string))
	ctx.Set("params", params)
}
