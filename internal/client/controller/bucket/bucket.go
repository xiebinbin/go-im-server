package bucket

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/bucket"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type GetAwsHeadObjectRequest = bucket.GetAwsHeadObjectRequest

func StartBucket(ctx *gin.Context) {
	data := bucket.StartBucketV2()
	//if ctx.GetHeader("Os") == "ios" {
	//	data = map[string]interface{}{
	//		"bucket_id":     "2lqxqjg9lnrt1",
	//		"is_accelerate": false,
	//	}
	//}
	response.RespData(ctx, data)
	return
}

func GetBucketInfo(ctx *gin.Context) {
	var params struct {
		BucketId string `json:"bucket_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ResErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := bucket.GetBucketInfoByBucketId(params.BucketId)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespData(ctx, data)
	return
}

func GetAwsHeadObjectExists(ctx *gin.Context) {
	var params GetAwsHeadObjectRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ResErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	res := bucket.GetAwsHeadObjectExists(params)
	response.RespData(ctx, map[string]interface{}{
		"res": res,
	})
	return
}
