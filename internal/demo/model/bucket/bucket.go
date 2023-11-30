package bucket

import (
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	bucket2 "imsdk/internal/common/dao/bucket"
	"imsdk/internal/common/pkg/config"
	"imsdk/internal/sdkserver/model/bucket"
)

type GetBucketInfoParams struct {
	BucketId string `json:"bucket_id" binding:"required"`
}

func GetBucketInfo(ctx *gin.Context, params GetBucketInfoParams) bucket2.Bucket {
	//(pkg.CurlResponse, error)
	params = GetBucketInfoParams{
		BucketId: "49ba59abbe56e057",
	}
	dataByte, _ := json.Marshal(params)
	//data, err := model.RequestIMServer("bucketInfo", string(dataByte))
	//if err != nil {
	//	return pkg.CurlResponse{}, err
	//}
	//return data, nil
	ak, _ := config.GetConfigAk()
	ctx.Set("ak", ak)
	ctx.Set("data", string(dataByte))
	res, _ := bucket.GetBucketInfo(ctx)
	return res
}
