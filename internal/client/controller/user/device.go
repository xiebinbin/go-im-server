package user

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user/device"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/response"
)

type DeviceStatusRequest = device.StatusRequest

func UpdateDeviceStatus(ctx *gin.Context) {
	var params DeviceStatusRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	//params.Os = base.GetTerType(os)
	ip := ctx.ClientIP()
	if ip == "" || funcs.HasLocalIPAddr(ip) {
		ip = funcs.RemoteIp(ctx.Request)
	}
	params.Ip = ip
	err := device.UpdateDeviceStatus(ctx, params)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}
