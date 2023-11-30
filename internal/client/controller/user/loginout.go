package user

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/model/user/login"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/response"
)

func LoginOut(ctx *gin.Context) {
	uid, _ := ctx.Get("uid")
	os := ctx.Value(base.HeaderFieldOs).(string)
	deviceName := funcs.GetHeadersFields(ctx, "Device-Name")
	loginOutData := login.OutRequest{
		UId:        uid.(string),
		Os:         os,
		DeviceName: deviceName,
	}
	login.Out(ctx, loginOutData)
	response.RespSuc(ctx)
	return
}

func LoginOutV2(ctx *gin.Context) {
	var params struct {
		UId string `json:"uid"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	if params.UId == "" {
		response.RespSuc(ctx)
		return
	}
	os := ctx.Value(base.HeaderFieldOs).(string)
	deviceId := ctx.Value(base.HeaderFieldOs).(string)
	loginOutData := login.OutRequest{
		UId:      params.UId,
		DeviceId: deviceId,
		Os:       os,
	}
	login.Out(ctx, loginOutData)
	response.RespSuc(ctx)
	return
}
