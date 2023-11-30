package bucket

import (
	"imsdk/pkg/app"
)

type Bucket struct {
	AccelerateHost string                 `toml:"AccelerateHost" json:"accelerate_host"`
	BaseHost       string                 `toml:"BaseHost" json:"base_host"`
	CropHost       string                 `toml:"CropHost" json:"crop_host"`
	Provider       string                 `toml:"Provider" json:"provider"`
	BucketId       string                 `toml:"BucketId" json:"bucket_id"`
	BucketName     string                 `toml:"BucketName" json:"bucket_name"`
	Region         string                 `toml:"Region" json:"region"`
	Expire         int64                  `json:"expire"`
	Sts            map[string]interface{} `json:"sts"`
}

type FileBaseField struct {
	BucketId string `json:"bucket_id"`
	FileType string `json:"file_type"`
	Text     string `json:"text"`
	Width    int8   `json:"width"`
	Height   int8   `json:"height"`
}

var buckets = map[string]Bucket{
	// production  (new account)
	"38b10711b195d0cf": {
		CropHost:       "https://cropv2.tdim.com.tr/crop?src=",
		BaseHost:       "https://tdpub.s3.eu-central-1.amazonaws.com/",
		AccelerateHost: "https://tdpub.s3-accelerate.amazonaws.com/",
		BucketId:       "38b10711b195d0cf",
		BucketName:     "tdpub",
		Region:         "eu-central-1",
	},
	"7a2e9ace81e810be": {
		CropHost:       "https://cropv2.tdim.com.tr/crop?src=",
		BaseHost:       "https://tdfrankurt.s3.eu-central-1.amazonaws.com/",
		AccelerateHost: "https://tdfrankurt.s3-accelerate.amazonaws.com/",
		BucketId:       "7a2e9ace81e810be",
		BucketName:     "tdfrankurt",
		Region:         "eu-central-1",
	},
	"a9fd5cf1ccc56003": {
		CropHost:       "https://crop-asia.tdim.com.tr/crop?src=",
		BaseHost:       "https://tdtokyo.s3.ap-northeast-1.amazonaws.com/",
		AccelerateHost: "https://tdtokyo.s3-accelerate.amazonaws.com/",
		BucketId:       "a9fd5cf1ccc56003",
		BucketName:     "tdtokyo",
		Region:         "ap-northeast-1",
	},
	// test (new account)
	"858251218d9da00a": {
		CropHost:       "https://crop-asia.tdim.com.tr/crop?src=",
		BaseHost:       "https://tdpubttokyo.s3.ap-northeast-1.amazonaws.com/",
		AccelerateHost: "https://tdpubttokyo.s3-accelerate.amazonaws.com/",
		BucketId:       "858251218d9da00a",
		BucketName:     "tdpubttokyo",
		Region:         "ap-northeast-1",
	},
	"040ea79af8d5f001": {
		CropHost:       "https://cropv2.tdim.com.tr/crop?src=",
		BaseHost:       "https://tdpub-test.s3.eu-central-1.amazonaws.com/",
		AccelerateHost: "",
		BucketId:       "040ea79af8d5f001",
		BucketName:     "tdpub-test",
		Region:         "eu-central-1",
	},
	"2ok55satkskw1": {
		CropHost:       "https://crop-asia.tdim.com.tr/crop?src=",
		BaseHost:       "https://ntok2.s3.ap-northeast-1.amazonaws.com/",
		AccelerateHost: "https://ntok2.s3-accelerate.amazonaws.com/",
		BucketId:       "2ok55satkskw1",
		BucketName:     "ntok2",
		Region:         "ap-northeast-1",
	},
	"2lqxqjg9lnrt1": {
		CropHost:       "https://cropv2.tdim.com.tr/crop?src=",
		BaseHost:       "https://ntd2.s3.eu-central-1.amazonaws.com/",
		AccelerateHost: "https://ntd2.s3-accelerate.amazonaws.com/",
		BucketId:       "2lqxqjg9lnrt1",
		BucketName:     "ntd2",
		Region:         "eu-central-1",
	},
	"44335b0272c2588f": {
		CropHost:       "https://crop-asia.tdim.com.tr/crop?src=",
		BaseHost:       "https://tmtok2.s3.ap-northeast-1.amazonaws.com/",
		AccelerateHost: "",
		BucketId:       "44335b0272c2588f",
		BucketName:     "client-logs-tok",
		Region:         "ap-northeast-1",
	},

	// old bucket (new account)
	"2ok55satkskw2": {
		CropHost:       "https://crop-asia.tdim.com.tr/crop?src=",
		BaseHost:       "https://tdtok2.s3.ap-northeast-1.amazonaws.com/",
		AccelerateHost: "https://tdtok2.s3-accelerate.amazonaws.com/",
		BucketId:       "2ok55satkskw2",
		BucketName:     "tdtok2",
		Region:         "ap-northeast-1",
	},

	"2lqxqjg9lnrt2": {
		CropHost:       "https://cropv2.tdim.com.tr/crop?src=",
		BaseHost:       "https://tdim.s3.eu-central-1.amazonaws.com/",
		AccelerateHost: "https://tdim.s3-accelerate.amazonaws.com/",
		BucketId:       "2lqxqjg9lnrt2",
		BucketName:     "tdim",
		Region:         "eu-central-1",
	},
	"local": {
		CropHost:       "https://cropv2.tdim.com.tr/crop?src=",
		BaseHost:       "https://ntd2.s3.eu-central-1.amazonaws.com/",
		AccelerateHost: "https://ntd2.s3-accelerate.amazonaws.com/",
		BucketId:       "local",
		BucketName:     "ntd2",
		Region:         "eu-central-1",
	},
}

func GetBucketInfoById(bucketId string) Bucket {
	if v, ok := buckets[bucketId]; ok {
		return v
	}
	return Bucket{}
}

func GetBucketInfoByIdV2(bucketId string) Bucket {
	var bucketInfo Bucket
	err := app.Config().Bind("buckets", bucketId, &bucketInfo)
	if err != nil {
		return bucketInfo
	}
	return bucketInfo
}
