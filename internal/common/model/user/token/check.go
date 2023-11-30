package token

import (
	"context"
	"encoding/base64"
	"imsdk/pkg/errno"
	"imsdk/pkg/log"
	"imsdk/pkg/redis"
	"strings"
)

type Info struct {
	UId string `json:"uid"`
}

type CheckTokenParams struct {
	Token    string `json:"token"`
	AK       string `json:"ak"`
	DeviceId string `json:"device_id"`
}

func CheckToken(ctx context.Context, params CheckTokenParams) (Info, error) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "CheckToken", "device": params.DeviceId})
	var res Info
	token := params.Token
	uidByte, err := base64.StdEncoding.DecodeString(token[32:])
	if err != nil {
		log.Logger().Error(logCtx, "token illegal, err: ", err)
		return res, errno.Add("token-illegal", errno.TokenErr)
	}

	uid := string(uidByte)
	log.Logger().Info(logCtx, map[string]string{"uid": uid, "token": token}, "step: 1 ")
	key := getTokenCacheKey(uid, params.DeviceId)
	oldKey := getTokenOldCacheKey(uid, params.DeviceId)
	cache, err := redis.Client.Get(key).Result()
	cacheOld, err := redis.Client.Get(oldKey).Result()
	log.Logger().Info(logCtx, "step: 2 ,key: ", key, ", cache: ", cache, " ,err: ", err)
	if err != nil && err != redis.NilErr {
		return res, errno.Add("get-failed", errno.SysErr)
	}
	if cache == "" && cacheOld == "" {
		return res, errno.Add("un-login", errno.TokenErr)
	}
	// todo delete old->new
	tokenSli := make([]string, 0)
	if cache != "" {
		tokenSli = strings.Split(cache, "|")
	} else {
		tokenSli = strings.Split(cacheOld, "|")
	}
	//expireTime, _ := strconv.Atoi(tokenSli[2])
	if tokenSli[0] != token {
		log.Logger().Error(logCtx, "step: 3 , token is not same, cache token: ", tokenSli[0])
		return res, errno.Add("token-error", errno.TokenErr)
	}
	//if int64(expireTime) < time.Now().Unix() {
	//	return emptyUserInfo, errno.Add("token-expire", errno.TokenErr)
	//}
	res = Info{
		UId: uid,
	}
	return res, nil
}
