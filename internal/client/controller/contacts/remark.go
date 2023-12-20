package contacts

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user/contact"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type UpdateRemarkRequest = contact.UpdateBaseRequest

func UpdateRemarkText(ctx *gin.Context) {
	var params UpdateRemarkRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid, _ := ctx.Get(base.HeaderFieldUID)
	err := contact.UpdateRemarkText(uid.(string), params.UId, params.RemarkText)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}
