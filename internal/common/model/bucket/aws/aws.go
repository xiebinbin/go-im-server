package aws

import (
	"context"
	"encoding/json"
	"fmt"
	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s32 "github.com/aws/aws-sdk-go-v2/service/s3"
	"imsdk/internal/common/dao/bucket"
	"imsdk/internal/common/pkg/aws"
	"imsdk/internal/common/pkg/ip"
	"imsdk/pkg/app"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"strings"
	"time"
)

type GetAwsHeadObjectRequest struct {
	BucketId string `json:"bucket_id" binding:"required"`
	Key      string `json:"key" binding:"required"`
}

const (
	StsAwsCacheTag = "tmm-im:aws:sts:"
	StsAliCacheTag = "tmm-im:ali:sts:"
	StsExpire      = 3600
)

func GetAwsBucketInfoByBucketId(bucketId string) (bucket.Bucket, error) {
	bucketInfo := bucket.GetBucketInfoById(bucketId)
	sts := GetSts(bucketInfo.BucketName, bucketInfo.Region)
	if sts == nil {
		return bucketInfo, errno.Add("fail", errno.DefErr)
	}
	bucketInfo.Sts = sts
	bucketInfo.Expire = sts["expire"].(int64)
	return bucketInfo, nil
}

func GetAwsHeadObjectExists(request GetAwsHeadObjectRequest) bool {
	stsInfo, _ := GetAwsBucketInfoByBucketId(request.BucketId)
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
		"s3_bucket_id":  bucketInfo["bucket_id"],
		"s3_bucket":     bucketInfo["bucket"],
		"s3_region":     bucketInfo["region"],
		"is_accelerate": isAccelerate,
	}
}

func GetSts(bucketName string, region string) map[string]interface{} {
	cacheTag := StsAwsCacheTag + funcs.Md516(bucketName+region)
	cacheRes, _ := redis.Client.Get(cacheTag).Result()
	res := make(map[string]interface{}, 0)
	fmt.Println("cacheRes---", cacheRes)
	if cacheRes != "" {
		err := json.Unmarshal([]byte(cacheRes), &res)
		if err != nil {
			return nil
		}
	}
	if len(res) == 0 { // failed to get cache or cache is empty
		return cachePreSign(bucketName, region)
	} else {
		expire := res["expire"].(float64)
		if int64(expire)-funcs.GetMillis() <= 10*60*1000 { // expired
			return cachePreSign(bucketName, region)
		}
		res["expire"] = int64(expire)
		return res
	}
}

func cachePreSign(bucketName string, region string) map[string]interface{} {
	cacheTag := StsAwsCacheTag + funcs.Md516(bucketName+region)
	stsInfo := aws.GetSts(region, aws.StsExpire)
	if stsInfo == nil {
		return nil
	}
	// todo check code
	cacheStr, _ := json.Marshal(stsInfo)
	err := redis.Client.Set(cacheTag, string(cacheStr), time.Duration(StsExpire-180)*time.Second).Err()
	if err != nil {
	}
	return stsInfo
}
