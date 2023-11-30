package bucket

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/demo/model/bucket"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

func GetBucketInfo(ctx *gin.Context) {
	var params bucket.GetBucketInfoParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	response.RespData(ctx, bucket.GetBucketInfo(ctx, params))
	return
}
