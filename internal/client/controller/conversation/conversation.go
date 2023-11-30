package conversation

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/conversation/notice"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func GetNoticeSettings(ctx *gin.Context) {
	//list := notice.GetNoticeSetting(ctx.Request.Header)
	list := notice.GetSetting(ctx)
	response.ResData(ctx, list)
	return
}

func GetNoticeSettingsV2(ctx *gin.Context) {
	var params notice.GetSeqRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	list := notice.GetSettingV2(ctx, params)
	response.RespListData(ctx, list)
	return
}
