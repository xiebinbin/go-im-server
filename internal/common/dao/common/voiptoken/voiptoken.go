package voiptoken

import (
	"encoding/json"
	"imsdk/pkg/redis"
)

type VoIpToken struct {
	Token string `json:"voip_token"`
	OS    string `json:"os"`
	Time  int64  `json:"time"`
}

func deviceTokenRedisCacheTag() string {
	return "chat:voip_token"
}

func deviceTokenUniqueCacheTag() string {
	return "chat:voip_token:unique_tag"
}

func AddVoIpTokenToRedis(uid string, params VoIpToken) error {
	data := map[string]interface{}{
		"token": params.Token,
		"time":  params.Time,
		"os":    params.OS,
	}
	dataSave, _ := json.Marshal(data)
	_, err := redis.Client.HSet(deviceTokenRedisCacheTag(), uid, string(dataSave)).Result()
	if err == nil {
		return nil
	}
	return err
}

func GetVoIpTokenUniqueInfo(token string) string {
	data := redis.Client.HGet(deviceTokenUniqueCacheTag(), token).Val()
	return data
}

func RemoveVoIpTokenUniqueInfo(token string) error {
	err := redis.Client.Del(deviceTokenUniqueCacheTag(), token).Err()
	return err
}

func AddVoIpTokenUniqueTag(token, uid string) error {
	_, err := redis.Client.HSet(deviceTokenUniqueCacheTag(), token, uid).Result()
	if err == nil {
		return nil
	}
	return err
}

func GetVoIpTokenDetailByUid(uid string) (map[string]interface{}, error) {
	data := redis.Client.HGet(deviceTokenRedisCacheTag(), uid).Val()
	resMap := make(map[string]interface{}, 0)
	if data != "" {
		err := json.Unmarshal([]byte(data), &resMap)
		if err != nil {
			return resMap, err
		}
	}
	return resMap, nil
}
func GetDeviceTokenDetailByUIds(uIds []string) ([]map[string]interface{}, error) {
	resMap := make([]map[string]interface{}, 0)
	for _, v := range uIds {
		tokenInfo, _ := GetVoIpTokenDetailByUid(v)
		if len(tokenInfo) == 0 {
			continue
		}
		resMap = append(resMap, tokenInfo)
	}
	return resMap, nil
}

func RemoveDeviceTokenByUid(uid string) error {
	_, err := redis.Client.HDel(deviceTokenRedisCacheTag(), uid).Result()
	if err == nil {
		return nil
	}
	return err
}

func EmptyTokenInfo(uid string) bool {
	info, _ := GetVoIpTokenDetailByUid(uid)
	if info["token"] == nil {
		return true
	}
	err := RemoveDeviceTokenByUid(uid)
	if err != nil {
	}
	//err = RemoveDeviceTokenUniqueInfo(info["token"].(string))
	//if err != nil {
	//}
	return true
}
