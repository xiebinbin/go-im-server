package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	json "github.com/json-iterator/go"
	"imsdk/pkg/app"
	"imsdk/pkg/funcs"
	"time"
)

type R2Resolver struct {
	AccountId string
}
type R2Account struct {
	AccountId       string `toml:"account_id"`
	AccessKeyId     string `toml:"access_key"`
	AccessKeySecret string `toml:"access_secret"`
}

func (r *R2Resolver) ResolveEndpoint(service, region string) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r.AccountId),
	}, nil
}

func GetR2Client() *s3.Client {
	//var bucketName = "bobobo-test"
	var r2Account R2Account
	app.Config().Bind("global", "r2_account", &r2Account)
	r2Resolver := &R2Resolver{
		AccountId: r2Account.AccountId,
	}
	fmt.Println("r2Account:", r2Account)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolver(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r2Account.AccessKeyId, r2Account.AccessKeySecret, "")),
	)
	if err != nil {
		//log.Fatal(err)
	}
	//client1 := sts.NewFromConfig(cfg)
	client := s3.NewFromConfig(cfg)
	return client
}

func GetPreSignURL(client *s3.Client, bucketName string) string {
	presignClient := s3.NewPresignClient(client)
	t := time.Unix(funcs.GetMillis()+604800000, 0)
	presignResult, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:  aws.String(bucketName),
		Key:     aws.String("miyaya.txt"),
		Expires: aws.Time(t),
	})

	if err != nil {
		panic("Couldn't get presigned URL for PutObject")
	}

	fmt.Printf("Presigned URL For object: %s\n", presignResult.URL)
	return presignResult.URL
}

func ListObjectsV2(client *s3.Client, bucketName string) {
	listObjectsOutput, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		//log.Fatal(err)
		fmt.Println("err:1", err)
	}
	fmt.Println("listObjectsOutput:", listObjectsOutput.Contents)
	//for _, object := range listObjectsOutput.Contents {
	//	obj, _ := json.MarshalIndent(object, "", "\t")
	//	fmt.Println("listObjectsOutput:", string(obj))
	//}

	//  {
	//  	"ChecksumAlgorithm": null,
	//  	"ETag": "\"eb2b891dc67b81755d2b726d9110af16\"",
	//  	"Key": "ferriswasm.png",
	//  	"LastModified": "2022-05-18T17:20:21.67Z",
	//  	"Owner": null,
	//  	"Size": 87671,
	//  	"StorageClass": "STANDARD"
	//  }
}

func ListBuckets(client *s3.Client, bucketName string) {
	listBucketsOutput, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		//log.Fatal(err)
	}

	for _, object := range listBucketsOutput.Buckets {
		obj, _ := json.MarshalIndent(object, "", "\t")
		fmt.Println(string(obj))
	}
}
