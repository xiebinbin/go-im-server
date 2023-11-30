package common

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/common"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/response"
)

func GetCountryList(ctx *gin.Context) {
	lang := funcs.GetHeadersFields(ctx, "Lang")
	list, err := common.GetCountryList(lang)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.ResData(ctx, list)
	return
}

func GetCityList(ctx *gin.Context) {
	lang := funcs.GetHeadersFields(ctx, "Lang")
	var params struct {
		Ver int `json:"ver"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	list, ver, err := common.GetCityList(lang, params.Ver)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	data := map[string]interface{}{
		"ver":   ver,
		"items": list,
	}
	response.RespData(ctx, data)
	return
}
