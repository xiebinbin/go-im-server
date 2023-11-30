package bucket

import (
	"imsdk/internal/common/dao/bucket"
	"imsdk/internal/common/model/bucket/ali"
	"imsdk/internal/common/model/bucket/aws"
)

func GetBucketInfoByBucketId(bucketId string) (bucket.Bucket, error) {
	//bucketInfo := bucket.GetBucketInfoById(bucketId)
	bucketInfo := bucket.GetBucketInfoByIdV2(bucketId)
	sts := make(map[string]interface{}, 0)
	if bucketInfo.Provider == "ali" {
		stsInfo, _ := ali.GetAliSts(bucketInfo.Region)
		sts = map[string]interface{}{
			"expire":            stsInfo.Expire,
			"access_key_id":     stsInfo.AccessKeyId,
			"access_key_secret": stsInfo.AccessKeySecret,
			"session_token":     stsInfo.SessionToken,
		}
		bucketInfo.Expire = stsInfo.Expire
	} else {
		//return bucket.Bucket{},errno.Add("ocean-test-err", errno.DataNotExist)
		sts = aws.GetSts(bucketInfo.BucketName, bucketInfo.Region)
		bucketInfo.Expire = sts["expire"].(int64)
	}
	bucketInfo.Sts = sts
	return bucketInfo, nil
}
