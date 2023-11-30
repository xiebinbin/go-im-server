package middlewares

import (
	"github.com/gin-gonic/gin"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type Headers struct {
	Os      int
	Version string
	Over    string
	Lang    string
}

func CheckHeaders(ctx *gin.Context) {
	token := ctx.Request.Header.Get("token")
	if len(token) < 32 {
		ctx.Abort()
		response.ResPubErr(ctx, errno.Add("header-err-token", errno.HeaderErr))
		return
	}
	//ctx.Set("lang", lang)
	//ctx.Set("is_enc", true)
	//ctx.Set("enc_key", getEncryptKey(os))
}
