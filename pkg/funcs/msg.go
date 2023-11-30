package funcs

import (
	"sort"
	"strconv"
	"strings"
)

const (
	ChatGroupType  = 2
	ChatSingleType = 1
)

func CreateMsgId(uid string) string {
	//md5(uid+timestamp_ms+6随机数+ “createmessage”)
	randStr := GetRandString(6)
	millis := strconv.FormatInt(GetMillis(), 10)
	return Md516(uid + millis + randStr + "createmessage")
}

func CreateSingleChatId(uid, targetId string) string {
	//uid(小)_uid(大)
	ids := []string{uid, targetId}
	sort.Strings(ids)
	return "s_" + Md516(ids[0]+"_"+ids[1])
}

func CreateGroupChatId(uid string) string {
	//g_md5(uid+time+随机数+"group")
	randStr := GetRandString(6)
	str := uid + randStr + "group"
	return "g_" + Md5Str(str)
}

func parseChatId(chatId string, senderId string) map[string]interface{} {
	chatIdInfo := strings.Split(chatId, "_")
	res := map[string]interface{}{
		"type":     0,
		"senderId": senderId,
		"targetId": "",
	}
	if chatIdInfo[0] == "s" {
		res["type"] = ChatSingleType
	} else if chatIdInfo[0] == "g" {
		res["type"] = ChatGroupType
	}
	i := 0
	chatIdInfo = append(chatIdInfo[:i], chatIdInfo[i+1:]...)
	sort.Strings(chatIdInfo)
	if len(chatIdInfo) == 1 {
		res["targetId"] = chatIdInfo[0]
	} else if len(chatIdInfo) == 2 {
		for _, v := range chatIdInfo {
			if v != senderId {
				res["targetId"] = v
			}
		}
	}
	return res
}
