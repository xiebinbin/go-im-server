package group

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/group"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

// AddApp group/addApp
func AddApp(ctx *gin.Context) {
	var params group.AddAppRequest
	uid, _ := ctx.Get("uid")
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := group.AddApp(ctx, uid.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func AppList(ctx *gin.Context) {
	var params group.IdsRequest
	uid, _ := ctx.Get("uid")
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	res, err := group.AppList(ctx, uid.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, res)
	return
}

func AppInfo(ctx *gin.Context) {
	var params group.IdsRequest
	uid, _ := ctx.Get("uid")
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	res, err := group.AppInfo(ctx, uid.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, res)
	return
}

func AppUpdate(ctx *gin.Context) {
	var params group.UpdateAppRequest
	uid, _ := ctx.Get("uid")
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := group.UpdateApp(ctx, uid.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func AppDelete(ctx *gin.Context) {
	var params group.IdsRequest
	uid, _ := ctx.Get("uid")
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := group.DeleteByIds(ctx, uid.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func DeleteAppByGIds(ctx *gin.Context) {
	var params group.DeleteByGIdsRequest
	uid, _ := ctx.Get("uid")
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := group.DeleteByGIds(ctx, uid.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
