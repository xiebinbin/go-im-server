package contacts

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user/contact"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type UpdateAliasRequest struct {
	UId   string `json:"uid" binding:"required"`
	Alias string `json:"alias"`
}
type UpdatePhoneRequest struct {
	UId   string `json:"uid" binding:"required"`
	Phone string `json:"phone"`
}

func UpdateAlias(ctx *gin.Context) {
	var params UpdateAliasRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid, _ := ctx.Get("uid")
	err := contact.UpdateAlias(uid.(string), params.UId, params.Alias)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}

func UpdatePhone(ctx *gin.Context) {
	var params UpdatePhoneRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid, _ := ctx.Get("uid")
	err := contact.UpdatePhone(uid.(string), params.UId, params.Phone)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}
