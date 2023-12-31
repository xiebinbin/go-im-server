package group

import (
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"imsdk/internal/client/model/group"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type GetNoticeRequest = group.GetNoticeRequest

func UpdateNotice(ctx *gin.Context) {
	var params group.UpdateNoticeRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.UpdateNotice(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func UpdateDesc(ctx *gin.Context) {
	var params group.UpdateDescRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.UpdateDesc(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func GetNotice(ctx *gin.Context) {
	var params GetNoticeRequest
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := group.GetNotice(uid, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespData(ctx, data)
	return
}

func GetDesc(ctx *gin.Context) {
	var params group.GetDescRequest
	uid := ctx.Value(base.HeaderFieldUID).(string)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := group.GetDesc(ctx, uid, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespData(ctx, data)
	return
}
