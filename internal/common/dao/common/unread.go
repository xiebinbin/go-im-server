package common

import (
	"fmt"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/redis"
	"strconv"
)

func GetAllBadge(uid string) int {
	sum := 0
	for _, v := range []string{base.MomentsUnreadKey, base.ImUnreadKey} {
		badge, _ := redis.Client.HGet(base.UnreadInfoHash, v+"_"+uid).Result()
		num, _ := strconv.Atoi(badge)
		fmt.Println(uid, "***** unread num ***********", v, "*********", num)
		sum += num
	}
	return sum
}
