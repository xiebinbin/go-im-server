package bucket

import (
	"encoding/json"
	"imsdk/internal/client/model/bucket"
	"imsdk/internal/common/pkg/aws"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetAwsHeadObjectRequest = bucket.GetAwsHeadObjectRequest
type GetPreSignURLRequest struct {
	Key string `json:"key" binding:"required"`
}

func GetPreSignURL(ctx *gin.Context) {
	var params GetPreSignURLRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	if len(params.Key) == 0 || !strings.HasPrefix(params.Key, "/") {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	url := aws.GetPreSignURL(aws.GetR2Client(), "bobobo-test", params.Key)
	response.RespData(ctx, map[string]interface{}{
		"url": url,
	})
}

func StartBucket(ctx *gin.Context) {
	//bucket.GetR2STS(ctx)
	//data := bucket.StartBucketV2()
	//response.RespData(ctx, data)
	return
}

func GetBucketInfo(ctx *gin.Context) {
	var params struct {
		BucketId string `json:"bucket_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
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
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	res := bucket.GetAwsHeadObjectExists(params)
	response.RespData(ctx, map[string]interface{}{
		"res": res,
	})
	return
}
