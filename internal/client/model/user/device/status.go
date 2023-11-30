package device

import (
	"context"
	"fmt"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"strconv"
	"time"
)

type StatusRequest struct {
	Status    int    `json:"status"` // 1-front 0-back
	Timestamp int    `json:"timestamp"`
	DeviceId  string `json:"device_id"`
	Os        string `json:"os"`
	Ip        string `json:"ip"`
}

const (
	CacheTag = base.UserDeviceStatus
)

func UpdateDeviceStatus(ctx context.Context, request StatusRequest) error {
	uid := ctx.Value("address").(string)
	deviceId := ctx.Value(base.HeaderFieldDeviceId).(string)
	key := uid + deviceId
	value := strconv.Itoa(request.Timestamp) + strconv.Itoa(request.Status)
	if value == "" {
		redis.Client.HSet(CacheTag, key, value)
	} else {
		val, _ := strconv.Atoi(value)
		cacheVal, _ := strconv.Atoi(GetCacheDeviceStatus(uid, deviceId))
		if val > cacheVal {
			redis.Client.HSet(CacheTag, key, value)
		}
	}
	fmt.Println("UpdateDeviceStatus value:", uid, value, key)
	RecordUserActive(uid)
	//AddDeviceInfo()
	return nil
}

func GetCacheDeviceStatus(uid string, deviceId string) string {
	key := uid + deviceId
	if !redis.Client.HExists(CacheTag, key).Val() {
		return ""
	}
	res := redis.Client.HGet(CacheTag, key).Val()
	return res
}

func GetDeviceStatus(uid, deviceId string) int {
	cacheVal := GetCacheDeviceStatus(uid, deviceId)
	status := -1
	if cacheVal != "" {
		status, _ = strconv.Atoi(cacheVal[13:])
	}
	return status
}

func RecordUserActive(uid string) {
	date := funcs.GetDate()
	key := "chat:online:members:count:" + strconv.Itoa(int(date))
	redis.Client.PFAdd(key, uid)
	redis.Client.Expire(key, time.Second*31*86400)
}
