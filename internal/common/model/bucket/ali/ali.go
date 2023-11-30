package ali

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"imsdk/internal/common/dao/bucket"
	"imsdk/internal/common/pkg/config"
	"imsdk/pkg/app"
	_ "imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/oss/aliyun"
	"imsdk/pkg/redis"
	"time"
)

const (
	StsAliCacheTag = "tmm-im:ali:sts:"
	StsExpire      = 3600
)

type Role struct {
	Arn         string `json:"arn"`
	SessionName string `json:"session_name"`
}
type Account struct {
	User            string `json:"user"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

type StsParams struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SessionToken    string `json:"session_token"`
	Expire          int64  `json:"expire"`
}

func StartBucket() map[string]interface{} {
	var bucketInfo bucket.Bucket
	err := app.Config().Bind("buckets", "49ba59abbe56e057", &bucketInfo)
	if err != nil {

	}
	return map[string]interface{}{
		"bucket_id":     bucketInfo.BucketId,
		"is_accelerate": false,
	}
}

func GetAliSts(region string) (StsParams, error) {
	regionDef, _ := config.GetDefRegion() //hongkong
	if region == "" {
		region = GetRegion(regionDef)
	}
	acc, _ := GetAccount(regionDef)
	role, _ := GetRole(regionDef)
	cacheTag := StsAliCacheTag + funcs.Md516(region+acc.AccessKeyId+acc.AccessKeySecret)
	fmt.Println("GetAliSts cacheTag:", cacheTag)
	cacheRes, _ := redis.Client.Get(cacheTag).Result()
	//cacheRes = ""
	if cacheRes == "" {
		return cacheAliPreSign(region, acc, role)
	}
	var res StsParams
	err := json.Unmarshal([]byte(cacheRes), &res)
	if err != nil {
		return StsParams{}, err
	}
	expire := res.Expire
	if expire-funcs.GetMillis() <= 10*60*1000 { // expired
		return cacheAliPreSign(region, acc, role)
	}
	return res, nil
}

func cacheAliPreSign(regionAddr string, acc Account, role Role) (StsParams, error) {
	logCtx := log.WithFields(context.Background(), map[string]string{"action": "cacheAliPreSign"})
	cacheTag := StsAliCacheTag + funcs.Md516(regionAddr+acc.AccessKeyId+acc.AccessKeySecret)
	stsResponse, err := aliyun.NewSTSClient(regionAddr, acc.AccessKeyId, acc.AccessKeySecret).
		SetRole(role.Arn, role.SessionName).GetSts()
	if err != nil {
		// todo
		return StsParams{}, err
	}
	log.Logger().Info(logCtx, "cacheAliPreSign stsResponse :  ", stsResponse)
	expire := time.Now().UnixNano()/1e6 + 3600*1000
	stsInfo := StsParams{
		AccessKeyId:     stsResponse.AccessKeyId,
		AccessKeySecret: stsResponse.AccessKeySecret,
		Expire:          expire,
		SessionToken:    stsResponse.SessionToken,
	}
	// todo check code
	cacheStr, _ := json.Marshal(stsInfo)
	err = redis.Client.Set(cacheTag, string(cacheStr), time.Duration(StsExpire)*time.Second).Err()
	log.Logger().Info(logCtx, "cacheAliPreSign redis.Client.Set :  ", cacheTag, string(cacheStr))
	if err != nil {
		log.Logger().Error(logCtx, "cacheAliPreSign redis.Client.Set err:  ", err)
	}
	return stsInfo, nil
}

func GetRegion(region string) string {
	res, err := app.Config().GetChildConf("aliyun", "region", region)
	if err != nil {
		return "ap-southeast-1"
	}
	return res.(string)
}

func GetAccount(region string) (Account, error) {
	var account Account
	err := app.Config().Bind("aliyun", region+"_sts_user", &account)
	if err != nil {
		// todo
	}
	return account, nil
}

func GetRole(region string) (Role, error) {
	var role Role
	err := app.Config().Bind("aliyun", region+"_sts_roles", &role)
	if err != nil {
		// todo
	}
	return role, nil
}
