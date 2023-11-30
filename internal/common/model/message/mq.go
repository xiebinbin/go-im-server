package message

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"imsdk/internal/common/dao/chat"
	chat2 "imsdk/internal/common/model/chat"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/log"
	"math"
)

const (
	TopicMsgToWaitReProcess          = base.TopicMsgToWaitReProcess
	WaitReprocessMsgTypeSingle uint8 = 1
	WaitReprocessMsgTypeChat   uint8 = 2
)

type WaitReprocessMsgMqData struct {
	Type          uint8     `json:"type"`
	ChatInfo      chat.Chat `json:"chat_info"`
	MsgItem       MsgItem   `json:"msg_item,omitempty"`
	PushedUsers   []string  `json:"pushed_users,omitempty"`
	ReceiverUsers []string  `json:"rec_users,omitempty"`
	ChunkUsers    []string  `json:"chunk_users,omitempty"`
	Uid           string    `json:"uid,omitempty"`
	AppCtx        string    `json:"app_ctx"`
	Ak            string    `json:"ak"`
	ReqId         string    `json:"req_id"`
	NoIncrUnread  bool      `json:"no_incr_unread"`
	NoPushOffline bool      `json:"no_push_offline"`
	HasBuild      bool      `json:"has_build"`
}

func ProduceMessageToWaitProcessMq(ctx context.Context, mqData WaitReprocessMsgMqData) error {
	mqData.Type = WaitReprocessMsgTypeChat
	p, _ := json.Marshal(mqData)
	//return kafka.ProduceMsg(ctx, TopicMsgToWaitReProcess, p)
	fmt.Println("p", p)
	return nil
}

func ProcessMessage(ctx context.Context, val []byte) error {
	var data WaitReprocessMsgMqData
	err := json.Unmarshal(val, &data)
	logCtx := log.WithFields(ctx, map[string]string{"action": "ProcessMessage", "data": string(val)})
	log.Logger().Info(logCtx, "start ---------->")
	if err != nil {
		log.Logger().Error(logCtx, "failed to decode data , err: ", err)
		return err
	}

	if data.Type == WaitReprocessMsgTypeSingle {
		log.Logger().Info(logCtx, "start ---------->  WaitReprocessMsgTypeSingle")
		_, err = SpreadMsgAndReqPush(ctx, data, true)
		return err
	}

	pushedUIds := make(map[string]struct{})
	if len(data.PushedUsers) > 0 {
		for _, uid := range data.PushedUsers {
			pushedUIds[uid] = struct{}{}
		}
	}
	log.Logger().Info(logCtx, "pushedUIds ---------->", pushedUIds)

	//var receiverUids []string
	receiverUids := data.ReceiverUsers
	mqReceiverLen := len(data.ReceiverUsers)
	if !data.HasBuild {
		//var extra struct {
		//	Filter []info.Filter `json:"filter"`
		//}
		//filter := make([]info.Filter, 0)
		//if data.MsgItem.Extra != "" {
		//	err = json.Unmarshal([]byte(data.MsgItem.Extra), &extra)
		//	if err != nil {
		//		return err
		//	}
		//	filter = extra.Filter
		//}
		receiverUids = []string{}
		fmt.Println("notice members : ", receiverUids)
		if err != nil {
			return err
		}
	} else if mqReceiverLen > 0 { // Receiver specified
		log.Logger().Info(logCtx, "mqReceiverLen > 0 ")
		if mqReceiverLen <= base.SyncProcessMaxAmount {
			for _, uid := range data.ReceiverUsers {
				if _, ok := pushedUIds[uid]; ok {
					continue
				}
				if err = ReprocessMessageSingle(ctx, uid, data); err != nil {
					log.Logger().Error(logCtx, "process message failed , err: ", err)
					return err
				}
			}
			return nil
		}
	} else {
		receiverUids, err = chat2.GetChatMemberUIds(logCtx, data.MsgItem.ChatId)
		log.Logger().Info(logCtx, "chat2.GetChatMemberUIds : ", receiverUids, err)
		if err != nil {
			log.Logger().Error(logCtx, "failed to query chat member , err: ", err)
			return err
		}
	}

	data.HasBuild = true
	allLen := len(receiverUids)
	groupAmount := int(math.Ceil(float64(allLen) / base.SyncProcessMaxAmount))
	log.Logger().Info(logCtx, "groupAmount : ", groupAmount)
	for i := 0; i < groupAmount; i++ {
		startIndex := i * base.SyncProcessMaxAmount
		endIndex := startIndex + base.SyncProcessMaxAmount
		if endIndex > allLen {
			endIndex = allLen
		}
		data.ReceiverUsers = receiverUids[startIndex:endIndex]
		fmt.Println("ProduceMessageToWaitProcessMq--->", data)
		err = ProduceMessageToWaitProcessMq(ctx, data)
		if err != nil {
			log.Logger().Error(logCtx, "group user again , err: ", err)
			return err
		}
	}

	return nil
}

func ReprocessMessageSingle(ctx context.Context, uid string, data WaitReprocessMsgMqData) error {
	//if _, ok := pushedUIds[data.Uid]; ok {
	//	// User who have already pushed successfully do not need to forward again
	//	return nil
	//}
	data.Uid = uid
	if _, err := SpreadMsgAndReqPush(ctx, data, true); err != nil {
		// If batch processing fails, put the single task into the queue again, wait for the next processing,
		// and continue to the next task in the batch processing instead of directly returning an error and blocking the process
		data.Type = WaitReprocessMsgTypeSingle
		p, _ := json.Marshal(data)
		//if err = kafka.ProduceMsg(ctx, TopicMsgToWaitReProcess, p); err != nil {
		//	return err
		//}
		fmt.Println("p", p)
	}
	return nil
}
