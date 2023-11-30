package user

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/user/setting"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type UpdateLangRequest struct {
	Language string `json:"language" binding:"required"`
}

type UpdateSettingParams struct {
	Val uint8  `json:"val"`
	Ts  uint64 `json:"ts"`
}

func UpdateLanguage(ctx *gin.Context) {
	var params UpdateLangRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid, _ := ctx.Get("uid")
	_, err := setting.UpdateLanguage(uid.(string), params.Language)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}
