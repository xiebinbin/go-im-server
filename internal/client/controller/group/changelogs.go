package group

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/group"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func Changelogs(ctx *gin.Context) {
	var params struct {
		Group []group.ChangeLogsListParams `json:"group"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}

	data, err := group.ChangeLogsList(params.Group)
	if err == nil {
		response.RespListData(ctx, data)
		return
	}
	response.RespErr(ctx, err)
	return
}
