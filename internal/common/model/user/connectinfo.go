package user

import (
	"context"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/redis"
)

type GetConnectInfoRequest struct {
	UIds  []string `json:"uids"`
	AUIds []string `json:"auids"`
	UId   string   `json:"uid"`
}

const (
	ConnectInfoLimit         = 50
	ErrorExceedQuantityLimit = 200001
	ErrorExceedTimesLimit    = 200002
)

func GetConnectInfo(ctx context.Context, request GetConnectInfoRequest) (map[string]int8, error) {
	res := make(map[string]int8, 0)
	for _, id := range request.UIds {
		cacheTag := base.RedisUserLatestConn + id
		isOnline := 0
		if redis.Client.Exists([]string{cacheTag}...).Val() == 1 {
			isOnline = 1
		}
		res[id] = int8(isOnline)
	}
	return res, nil
}
