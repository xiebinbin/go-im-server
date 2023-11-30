package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"imsdk/pkg/app"
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

func GetSts(region string, expire int32) map[string]interface{} {
	role, err := app.Config().GetChildConf("aws", "roles", "s3sts")
	var serverUser map[string]string
	app.Config().Bind("aws", "user_serveruser", &serverUser)
	if err != nil {
		fmt.Println("action:get sts:desc:field to get role conf", err)
		return nil
	}
	fmt.Println(role, serverUser)
	sessionName := "tmm_s3_sts"
	//credentials := aws.Credentials{AccessKeyID: serverUser["key"], SecretAccessKey: serverUser["secret"]}
	staticProvider := credentials.NewStaticCredentialsProvider(
		serverUser["key"],
		serverUser["secret"],
		"",
	)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(staticProvider))
	if err != nil {
		fmt.Println("action:get sts:desc:field to load role conf", err)
		return nil
	}

	client := sts.NewFromConfig(cfg)
	roleArn := role.(string)
	input := &sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: &sessionName,
		DurationSeconds: &expire,
	}

	res, err := client.AssumeRole(context.TODO(), input)
	if err != nil {
		fmt.Println("action:get sts:desc:failed to request sts api:", err)
		return nil
	}

	return map[string]interface{}{
		"expire":            (res.Credentials.Expiration.UnixNano() / 1e6) - 5*180,
		"access_key_id":     *res.Credentials.AccessKeyId,
		"access_key_secret": *res.Credentials.SecretAccessKey,
		"session_token":     *res.Credentials.SessionToken,
	}
}
