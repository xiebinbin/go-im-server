package contacts

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user/contact"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type UpdateDatesRequest = contact.UpdateBaseRequest

func UpdateDates(ctx *gin.Context) {
	var params UpdateDatesRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid, _ := ctx.Get("uid")
	err := contact.UpdateDates(uid.(string), params.UId, params.Dates)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}
