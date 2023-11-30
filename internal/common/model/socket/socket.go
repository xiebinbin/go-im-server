package socket

import (
	"context"
	"imsdk/pkg/errno"
	"imsdk/pkg/log"
	"imsdk/pkg/redis"
	"strings"
)

type WebsocketType = int8

const (
	clientConnServerCacheKey                = "imsdk:socket:conn:"
	clientConnPushUrlCacheKey               = "imsdk:pushurl:conn:"
	ForceDelSocketConnTraceId               = "forcedel"
	TypeWebSocket             WebsocketType = 0
	TypeStreamServer          WebsocketType = 1
)

type UserSocketConnection struct {
	Uid      string `json:"uid"`
	DeviceId string `json:"device_id"`
	Host     string `json:"host"`
	TraceId  string `json:"trace_id"`
}

type ConnectionList struct {
	Uid   string
	Hosts []UserSocketConnection
}

type SetUserPushUrlReq = UserSocketConnection

var (
	PushUrlInfo = make(map[string]string, 0)
)

func SetUserPushUrl(ctx context.Context, req SetUserPushUrlReq) (bool, error) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "SetUserPushUrl"})
	if v, ok := PushUrlInfo[req.Uid]; ok && v == req.Host {
		log.Logger().Info(logCtx, "-1:", req.Uid, v, req.Host)
		return true, nil
	}
	PushUrlInfo[req.Uid] = req.Host
	log.Logger().Info(logCtx, "-2:", req.Uid, req.Host)
	res, err := SaveUserPushUrl(req.Uid, req.DeviceId, req.Host, "")
	if err != nil || !res {
		log.Logger().Info(logCtx, "-3:", "err:", err, ":res:", res)
		return true, errno.Add("connect fail", errno.SysErr)
	}
	return res, err
}

func getCacheKey(uid string) string {
	return clientConnServerCacheKey + ":" + uid
}

// SaveSocketConnection redis hash key: clientConnServerCacheKey+uid ,field: terType, value: host
func SaveSocketConnection(uid, deviceId, host, traceId string) (bool, error) {
	return redis.Client.HSet(getCacheKey(uid), deviceId, host+"|"+traceId).Result()
}

func SaveUserPushUrl(uid, deviceId, host, traceId string) (bool, error) {
	cacheTag := clientConnPushUrlCacheKey + ":" + uid
	return redis.Client.HSet(cacheTag, deviceId, host+"|"+traceId).Result()
}

func GetUserPushUrl(uid, deviceId string) (string, error) {
	cacheTag := clientConnPushUrlCacheKey + ":" + uid
	host, err := redis.Client.HGet(cacheTag, deviceId).Result()
	if err == nil {
		info := strings.Split(host, "|")
		return info[0], nil
	} else if err == redis.NilErr {
		return "", nil
	}
	return "", err
}

func DelSocketConnection(ctx context.Context, uid, deviceId, traceId string) (int64, error) {
	key := getCacheKey(uid)
	host, err := redis.Client.HGet(key, deviceId).Result()
	logCtx := log.WithFields(ctx, map[string]string{"action": "DelSocketConnection"})
	if err == nil {
		log.Logger().Info(logCtx, "DelSocketConnection-1: ", host)
		info := strings.Split(host, "|")
		log.Logger().Info(logCtx, "DelSocketConnection-2: ", info)
		if traceId == ForceDelSocketConnTraceId || len(info) < 2 || info[1] == traceId {
			res, er := redis.Client.HDel(key, deviceId).Result()
			return res, er
		}
		return 0, nil
	} else if err == redis.NilErr {
		return 1, nil
	}
	return 0, err
}

func GetSocketConnectionHost(uid, deviceId string) (string, error) {
	host, err := redis.Client.HGet(getCacheKey(uid), deviceId).Result()
	if err == nil {
		info := strings.Split(host, "|")
		return info[0], nil
	} else if err == redis.NilErr {
		return "", nil
	}
	return "", err
}

func GetUserConnections(uid string) []UserSocketConnection {
	list := make([]UserSocketConnection, 0)
	connects, _ := redis.Client.HGetAll(getCacheKey(uid)).Result()
	if connects != nil {
		for deviceId, host := range connects {
			item := UserSocketConnection{
				DeviceId: deviceId,
				Uid:      uid,
				Host:     dealHostStr(host),
			}
			list = append(list, item)
		}
	}
	return list
}

func GetTerminalConnection(uids []string) []ConnectionList {
	list := make([]ConnectionList, 0)
	for _, uid := range uids {
		connects, _ := redis.Client.HGetAll(getCacheKey(uid)).Result()
		if connects != nil {
			userList := make([]UserSocketConnection, 0)
			for deviceId, host := range connects {
				item := UserSocketConnection{
					DeviceId: deviceId,
					Uid:      uid,
					Host:     dealHostStr(host),
				}
				userList = append(userList, item)
			}
			list = append(list, ConnectionList{Uid: uid, Hosts: userList})
		}
	}
	return list
}

func dealHostStr(str string) string {
	info := strings.Split(str, "|")
	return info[0]
}
