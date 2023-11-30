package initialization

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	chat2 "imsdk/internal/common/dao/chat"
	"imsdk/internal/common/dao/chat/members"
	"imsdk/internal/common/dao/message/usermessage"
	"imsdk/pkg/errno"
	"imsdk/pkg/log"
	"math"
	"strings"
)

func DealChatManyMember(ctx context.Context) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "DealChatManyMember"})
	where := bson.M{}
	count := chat2.New().GetCount(where)
	//count = 10
	row := 1000
	page := int(math.Ceil(float64(count) / float64(row)))
	fmt.Println("count:", count, page)
	log.Logger().Info(logCtx, "count:", count, ":page total:", page)
	for i := 0; i < page; i++ {
		data := chat2.New().GetListByLimit(int64(row), int64(i*row), where)
		log.Logger().Info(logCtx, "deal page ing-------:", i)
		if len(data) > 0 {
			for _, v := range data {
				uIds, _ := members.New().GetChatMembers(v.ID)
				if len(uIds) == 2 {
					continue
				}
				userMsgDao := usermessage.New()
				log.Logger().Info(logCtx, "change isread ing....... chatId:", v.ID, ":msgId:", v.ID, ":uids:", uIds)
				for _, id := range uIds {
					_, err := userMsgDao.UpdateIsReadAll(id, 1)
					if err != nil {
						log.Logger().Error(logCtx, "chatId:", v.ID, ":uid:", id, ":mids:", []string{v.ID}, ":err:", err)
					}
				}
			}
		}
		log.Logger().Info(logCtx, "deal page success ********** :", i)
	}
}

type ChatIDDetail struct {
	Type     uint8
	SenderId string // only available for Single conversation
	TargetId string // Single conversation: receive user id | group conversation: group id
}

func GetChatIDDetail(chatId, senderId string) (ChatIDDetail, error) {
	//hollowId := "33b1hvbe9lxf"
	var result ChatIDDetail
	convInfo := strings.Split(chatId, "_")
	if len(convInfo) > 0 && convInfo[0] == "s" {
		if len(convInfo) < 3 {
			return result, errno.Add("params-format-error", errno.DefErr)
		}
		fmt.Println("convInfo----", convInfo)
		targetId := convInfo[1]
		if senderId == convInfo[1] {
			targetId = convInfo[2]
		}
		result = ChatIDDetail{
			Type:     1,
			SenderId: senderId,
			TargetId: targetId,
		}
		return result, nil
	}
	return result, nil
}

func dealReadStatus(ctx context.Context, senderId string, msgContentStr string) bool {
	tmpRes := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(msgContentStr), &tmpRes)
	logCtx := log.WithFields(ctx, map[string]string{"action": "dealReadStatus"})
	if err == nil {
		mIds := make([]string, 0)
		for _, v := range tmpRes["mids"].([]interface{}) {
			mIds = append(mIds, v.(string))
		}
		userMsgDao := usermessage.New()
		_, upErr := userMsgDao.UpdateIsRead(senderId, mIds, usermessage.IsReadYes)
		if upErr != nil {
			log.Logger().Error(logCtx, "update read info err", err)
		}
	}
	return true
}
