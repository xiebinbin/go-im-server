package contacts

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user/contact"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type UpdatePhonesRequest = contact.UpdateBaseRequest

func UpdatePhones(ctx *gin.Context) {
	var params UpdatePhonesRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid, _ := ctx.Get("uid")
	err := contact.UpdatePhones(uid.(string), params.UId, params.Phones)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}
