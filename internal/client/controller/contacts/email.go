package contacts

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user/contact"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type UpdateEmailsRequest = contact.UpdateBaseRequest

func UpdateEmails(ctx *gin.Context) {
	var params UpdateEmailsRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid, _ := ctx.Get("uid")
	err := contact.UpdateEmails(uid.(string), params.UId, params.Emails)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}
