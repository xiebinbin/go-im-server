package active

import (
	"fmt"
	"imsdk/internal/common/dao/chat/members"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"strconv"
)

const OnlineBTTime = 30

type ChatActiveRequest struct {
	ChatId    string `json:"chat_id" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

func SaveActiveChat(uid string, request ChatActiveRequest) {
	if !members.New().IsExist(uid, request.ChatId) {
		return
	}
	if request.Timestamp < funcs.GetTimeSecs() {
		return
	}
	cacheSetTag := base.UserActiveChatSet
	lastChatId, _ := GetActiveChatByUid(uid)
	if lastChatId != "" && lastChatId != request.ChatId {
		redis.Client.SRem(cacheSetTag+lastChatId, uid)
	}
	SetActiveByUid(uid, request.ChatId)
	redis.Client.SAdd(cacheSetTag+request.ChatId, uid)
}

func DelActiveUserByChatId(chatId, uid string) {
	cacheTag := base.UserActiveChatSet
	err := redis.Client.SRem(cacheTag+chatId, uid).Err()
	//testRes := GetActiveUserByChatId(chatId)
	//fmt.Println("DelActiveUserByChatId : ", chatId, uid, err, testRes)
	if err != nil {
		fmt.Println("DelActiveUserByChatId err:", err)
		return
	}
	return
}

func GetActiveUserByChatId(chatId string) []string {
	cacheTag := base.UserActiveChatSet
	res, _ := redis.Client.SMembers(cacheTag + chatId).Result()
	return res
}

func GetActiveChatByUid(uid string) (string, error) {
	cacheTag := base.UserActiveChatHash + ":" + uid
	res, err := redis.Client.HGet(cacheTag, "ac").Result()
	if err == redis.NilErr {
		return "", nil
	}
	return res, nil
}

func SetActiveByUid(uid, chatId string) {
	cacheTag := base.UserActiveChatHash + ":" + uid
	script := "local ret=redis.call('HSet',KEYS[1],ARGV[1],ARGV[2]);redis.call('expire',KEYS[1],ARGV[3]);return ret"
	redis.Client.Eval(script, []string{cacheTag}, "ac", chatId, 300)
}

func SetOnlineByUid(uid string, time int64) {
	if time == 0 {
		time = funcs.GetTimeSecs()
	}

	cacheTag := base.UserActiveChatHash + ":" + uid
	script := "local ret=redis.call('HSet',KEYS[1],ARGV[1],ARGV[2]);redis.call('expire',KEYS[1],ARGV[3]);return ret"
	redis.Client.Eval(script, []string{cacheTag}, "at", time, 300)

	members.New().UpdateUserOnlineTime(uid, time)
}

func UserOnlineStatus(uid string) bool {
	cacheTag := base.UserActiveChatHash + ":" + uid
	onlineTimeVal, err := redis.Client.HGet(cacheTag, "at").Result()
	if err == redis.NilErr {
		return false
	}

	onlineTime, err := strconv.ParseInt(onlineTimeVal, 10, 64)
	if err == nil && funcs.GetTimeSecs()-OnlineBTTime < onlineTime {
		return true
	}
	return false
}
