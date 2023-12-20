package middlewares

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
	"imsdk/pkg/sdk"
	"imsdk/pkg/unique"
	"net/http"
)

type Headers struct {
	Os      int
	Version string
	Over    string
	Lang    string
}

func CheckHeaders(ctx *gin.Context) {
	pubKey := ctx.Request.Header.Get(base.HeaderFieldPubKey)
	uid := ctx.Request.Header.Get(base.HeaderFieldUID)
	// 计算发送者的地址与address对比
	sign := ctx.Request.Header.Get(base.HeaderFieldSign)
	time := ctx.Request.Header.Get(base.HeaderFieldTime)
	if pubKey == "" || uid == "" || sign == "" || time == "" {
		ctx.Abort()
		response.ResPubErr(ctx, errno.Add("header-params-err", http.StatusBadRequest))
		return
	}
	reqId := ctx.Request.Header.Get(base.HeaderFieldReqId)
	if reqId == "" {
		reqId = unique.Id12()
	}

	version := ctx.Request.Header.Get(base.HeaderFieldVersion)
	if version == sdk.EmptyString {
		version = "0.0.0.0"
	}

	userAgent := ctx.Request.Header.Get(base.HeaderFieldUserAgent)
	if userAgent == sdk.EmptyString {
		userAgent = "UNKNOWN"
	}

	deviceId := ctx.Request.Header.Get(base.HeaderFieldDeviceId)
	if deviceId == sdk.EmptyString {
		deviceId = "DEFAULT"
	}

	lang := ctx.Request.Header.Get(base.HeaderFieldLang)
	if lang == sdk.EmptyString {
		lang = "en"
	}

	isEnc := ctx.Request.Header.Get(base.HeaderIsEnc)
	if isEnc == sdk.EmptyString {
		isEnc = ""
	}

	ctx.Set(base.HeaderFieldTime, time)
	ctx.Set(base.HeaderFieldSign, sign)
	ctx.Set(base.HeaderFieldPubKey, pubKey)
	ctx.Set(base.HeaderFieldUID, uid)
	ctx.Set(base.HeaderFieldReqId, reqId)
	ctx.Set(base.HeaderFieldVersion, version)
	ctx.Set(base.HeaderFieldUserAgent, userAgent)
	ctx.Set(base.HeaderFieldDeviceId, deviceId)
	ctx.Set(base.HeaderFieldLang, lang)
	ctx.Set(base.HeaderIsEnc, isEnc)

	//os := strings.ToLower(ctx.Request.Header.Get(base.HeaderFieldOs))
}
