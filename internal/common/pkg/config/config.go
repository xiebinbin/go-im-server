package config

import (
	"fmt"
	"imsdk/pkg/app"
)

var (
	sk string
	ak,
	hollowUId,
	offlineCallBackUrl,
	staticUrl,
	defBucketId string
)

func GetHollowUId() (string, error) {
	if hollowUId != "" {
		return hollowUId, nil
	}
	hollowUIdConf, err := app.Config().GetChildConf("global", "system_user", "im_hollow_uid")
	if err != nil {
		fmt.Println("GetHollowUId Err:", err)
		return "", nil
	}
	hollowUId = hollowUIdConf.(string)
	return hollowUId, nil
}

func GetOfflineCallBackUrl() (string, error) {
	if offlineCallBackUrl != "" {
		return offlineCallBackUrl, nil
	}
	url, err := app.Config().GetChildConf("global", "hosts", "offline_callback")
	if err != nil {
		fmt.Println("GetOfflineCallBackUrl Err:", err)
		return "", nil
	}
	offlineCallBackUrl = url.(string)
	return offlineCallBackUrl, nil
}

func GetDefBucketId() (string, error) {
	if defBucketId != "" {
		return defBucketId, nil
	}
	defBucketIdConf, _ := app.Config().GetChildConf("global", "system", "oss_def_bucket_id")
	defBucketId = defBucketIdConf.(string)
	return defBucketId, nil
}

func GetStaticUrl() (string, error) {
	if staticUrl != "" {
		return staticUrl, nil
	}
	staticUrlConf, err := app.Config().GetChildConf("global", "system", "static_url")
	if err != nil {
		return "", nil
	}
	return staticUrlConf.(string), nil
}

func GetDefRegion() (string, error) {
	defRegion, err := app.Config().GetChildConf("global", "system", "oss_def_region")
	if err != nil {
		fmt.Println("GetDefRegion Err:", err)
		return "singapore", nil
	}
	return defRegion.(string), nil
}

func GetConfigAk() (string, error) {
	if ak != "" {
		return ak, nil
	}
	Ak, _ := app.Config().GetChildConf("global", "system", "pk")
	if Ak == nil {
		fmt.Println("not read ak config:", ak)
		return "", nil
	}
	ak = Ak.(string)
	return ak, nil
}

func GetConfigSk() (string, error) {
	if sk != "" {
		return sk, nil
	}
	Sk, _ := app.Config().GetChildConf("global", "system", "sk")
	sk = Sk.(string)
	return sk, nil
}

func GetConfigTopic(key string) (string, error) {
	topic, err := app.Config().GetChildConf("mq", "topictag", key)
	if err != nil {
		return "", nil
	}
	return topic.(string), nil
}

type IMCredentials struct {
	AK string `json:"ak"`
	SK string `json:"sk"`
}

func GetIMSdkKey() (IMCredentials, error) {
	// todo config file get
	ak, _ = GetConfigAk()
	sk, _ = GetConfigSk()
	return IMCredentials{
		AK: ak,
		SK: sk,
	}, nil
}
