package contacts

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user/contact"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type UpdateRelationsRequest = contact.UpdateBaseRequest

func UpdateRelations(ctx *gin.Context) {
	var params UpdateRelationsRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid, _ := ctx.Get("uid")
	err := contact.UpdateRelations(uid.(string), params.UId, params.Relations)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}
