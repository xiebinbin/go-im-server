package token

import (
	"context"
	"encoding/base64"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/config"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/redis"
	"strconv"
	"strings"
	"time"
)

const (
	LoginTokenExpire = 86400 * 30
)

type Token struct {
	Token        string
	RefreshToken string
	Expire       int64
}

type GetTokenParams struct {
	Uid          string `json:"uid"`
	AK           string `json:"ak"`
	DeviceId     string `json:"device_id"`
	RefreshToken string `json:"refresh_token"`
}

func getTokenOldCacheKey(uid, deviceId string) string {
	cacheKey := base.TypeUserSDK + deviceId + ":" + uid
	version := redis.Client.Get("user:version:" + uid).Val()
	ak, _ := config.GetConfigAk()
	if ak == base.AKSingapore && funcs.CompareVersion(version, "1.1.0") == -1 {
		if funcs.CompareVersion(version, "1.1.0") == -1 {
			if deviceId == base.AppDeviceId {
				cacheKey = base.TypeUserApp + uid
			}
			if deviceId == base.WebDeviceId {
				cacheKey = base.TypeUserPc + uid
			}
		}
	}
	return cacheKey
}
func getTokenCacheKey(uid, deviceId string) string {
	cacheKey := base.TypeUserSDK + deviceId + ":" + uid
	return cacheKey
}

func GetToken(ctx context.Context, params GetTokenParams) (Token, error) {
	var res Token
	key := getTokenCacheKey(params.Uid, params.DeviceId)
	t := time.Now().Unix()
	uidStr := base64.StdEncoding.EncodeToString([]byte(params.Uid))
	token := strings.ToUpper(funcs.Md5Str(funcs.GetRandString(32)+params.Uid+string(t))) + uidStr

	logCtx := log.WithFields(ctx, map[string]string{"action": "GetToken", "uid": params.Uid, "device_id": params.DeviceId})
	refreshToken := params.RefreshToken
	if params.RefreshToken == "" {
		refreshToken = strings.ToUpper(funcs.Md5Str(funcs.GetRandString(32) + params.Uid + string(t)))
	} else {
		cache, err := redis.Client.Get(key).Result()
		if err != nil && err != redis.NilErr {
			log.Logger().Error(logCtx, "failed to generate token, err: ", err)
			return Token{}, errno.Add("failed to generate token", errno.SysErr)
		}
		data := strings.Split(cache, "|")
		if len(data) < 3 || data[1] != params.RefreshToken {
			return Token{}, errno.Add("wrong-req", errno.WrongReq)
		}
	}

	cacheStr := token + "|" + res.RefreshToken + "|" + strconv.Itoa(0)
	log.Logger().Info(logCtx, "key--", key, cacheStr)
	if _, err := redis.Client.Set(key, cacheStr, 0).Result(); err != nil {
		log.Logger().Error(logCtx, " failed to save token, err: ", err)
		return res, errno.Add("fail", errno.DefErr)
	}
	res = Token{
		Token:        token,
		Expire:       0,
		RefreshToken: refreshToken,
	}
	return res, nil
}

func IsLatestVersion(uid, deviceId string) bool {
	cacheKey := uid + deviceId
	cacheVer, _ := redis.Client.HMGet(base.RedisUserLatestVer, cacheKey).Result()
	latestVer := "2.2.0"
	userVer := "0"
	if len(cacheVer) != 0 && cacheVer[0] != nil {
		userVer = cacheVer[0].(string)
	}
	res := funcs.CompareVersion(latestVer, userVer)
	if res == 1 {
		return false
	}
	return true
}

func RemoveToken(ctx context.Context, uid, deviceId string) bool {
	logCtx := log.WithFields(ctx, map[string]string{"action": "RemoveToken", "uid": uid, "device_id": deviceId})
	key := getTokenCacheKey(uid, deviceId)
	_, err := redis.Client.Del(key).Result()
	if err != nil {
		log.Logger().Info(logCtx, "redis del err:", err, ":key:", key)
	}
	return true
}
