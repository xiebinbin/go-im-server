package bucket

import (
	"context"
	"fmt"
	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s32 "github.com/aws/aws-sdk-go-v2/service/s3"
	"imsdk/internal/common/dao/bucket"
	bucketModel "imsdk/internal/common/model/bucket"
	"imsdk/internal/common/pkg/config"
	"imsdk/internal/common/pkg/ip"
	"imsdk/pkg/app"
	"imsdk/pkg/funcs"
	"strconv"
	"strings"
)

const (
	StsCacheTag = "chat-im:aws:sts:"
	StsExpire   = 3600
)

type GetAwsHeadObjectRequest struct {
	BucketId string `json:"bucket_id" binding:"required"`
	Key      string `json:"key" binding:"required"`
}

func GetBucketInfoByBucketId(bucketId string) (bucket.Bucket, error) {
	return bucketModel.GetBucketInfoByBucketId(bucketId)
}

func GetAwsHeadObjectExists(request GetAwsHeadObjectRequest) bool {
	t1 := funcs.GetMillis()
	stsInfo, _ := GetBucketInfoByBucketId(request.BucketId)
	accessKey := ""
	secretKey := ""
	sts := ""
	if stsInfo.Sts != nil {
		accessKey = stsInfo.Sts["access_key_id"].(string)
		secretKey = stsInfo.Sts["access_key_secret"].(string)
		sts = stsInfo.Sts["session_token"].(string)
	}
	if sts == "" {
		return false
	}
	region := stsInfo.Region
	client1 := s32.New(s32.Options{
		Region:        region,
		UseAccelerate: true,
		Credentials:   aws2.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, sts)),
	})
	bucketName := stsInfo.BucketName
	key := request.Key
	_, err := client1.HeadObject(context.TODO(), &s32.HeadObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	t2 := funcs.GetMillis()
	dur := t2 - t1
	fmt.Println("*********** HeadObject DurTime ***********", strconv.Itoa(int(dur)))
	return err == nil
}

func GetBucketByRegion(info ip.Info) map[string]interface{} {
	timezone := info.Timezone
	region := "default"
	if timezone != "" {
		sep := "/"
		arr := strings.Split(timezone, sep)
		tmpRegion := strings.ToLower(arr[0])
		agreeRegion := []string{"asia", "europe"}
		if funcs.In(tmpRegion, agreeRegion) {
			region = tmpRegion
		}
	}
	var bucketInfo map[string]string
	app.Config().Bind("aws", "s3_"+strings.ToLower(region), &bucketInfo)
	isAccelerate := true
	if strings.ToLower(info.Country) == "china" {
		isAccelerate = false
	}
	isAccelerate = false
	return map[string]interface{}{
		"bucket_id":     bucketInfo["bucket_id"],
		"bucket":        bucketInfo["bucket"],
		"region":        bucketInfo["region"],
		"is_accelerate": isAccelerate,
	}
}

func StartBucketV2() map[string]interface{} {
	defBucketId, _ := config.GetDefBucketId()
	var bucketInfo bucket.Bucket
	err := app.Config().Bind("buckets", defBucketId, &bucketInfo)
	if err != nil {
		fmt.Println(bucketInfo, err)
	}
	return map[string]interface{}{
		"bucket_id":     defBucketId,
		"is_accelerate": false,
	}
}
