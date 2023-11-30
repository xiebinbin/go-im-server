package forward

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"imsdk/internal/common/dao/chat"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/log"
)

const (
	TopicSocketMessagePush = base.TopicSocketMessagePush
	TopicMsgToOffline      = base.TopicMsgToOffline
)

type PushMessageParams struct {
	Host          string         `json:"host"` // used by socket forward (online)
	Cmd           string         `json:"cmd"`
	Uid           string         `json:"uid"`
	Data          MsgItem        `json:"items"`
	ChatInfo      chat.Chat      `json:"chat_info,omitempty"`
	OfflineData   OfflineRequest `json:"offline,omitempty"`
	SenderInfo    SenderInfo     `json:"sender_info,omitempty"`
	Devices       []string       `json:"device_ids,omitempty"`
	Ak            string         `json:"ak"`
	ReqId         string         `json:"req_id"`
	NoIncrUnread  bool           `json:"no_incr_unread"`  // business message (server send)
	AppCtx        string         `json:"app_ctx"`         // business message (server send)
	NoPushOffline bool           `json:"no_push_offline"` // business message (server send)
}

// ProduceSocketMessage
// Add the request data to be pushed to other socket servers into the queue
func ProduceSocketMessage(ctx context.Context, params PushMessageParams) error {
	p, _ := json.Marshal(params)
	//return kafka.ProduceMsg(ctx, TopicSocketMessagePush, p)
	fmt.Println("p", p)
	return nil
}

// PushMessageToSocket
// process in consumer
// Take the message to be pushed from the queue and forward it to the user by the socket.
// If it fails, it will join the offline forward queue
func PushMessageToSocket(ctx context.Context, val []byte) error {
	var data PushMessageParams
	err := json.Unmarshal(val, &data)
	logCtx := log.WithFields(ctx, map[string]string{"action": "PushMessageToSocket"})
	if err != nil {
		log.Logger().Error(logCtx, "failed to decode data , err: ", err)
		return err
	}
	socketData := SocketData{
		Cmd:  data.Cmd,
		Data: data.Data,
	}
	oriSocketData := ""
	if data.Cmd == CmdApplicationCmd {
		socketData.Data = data.Data.Content
		oriSocketData = data.Data.Content
	}
	offlineDeviceIds, er := PushMessageToUserSocketWithApi(ctx, data.Host, data.Uid, data.Devices, socketData, oriSocketData)
	if er != nil {
		log.Logger().Error(logCtx, "failed to forward message to user socket with api , err: ", er)
		return err
	}
	log.Logger().Info(logCtx, "PushMessageToSocket offlineDeviceIds:", offlineDeviceIds, data.Uid)
	if len(offlineDeviceIds) == 0 || data.NoPushOffline {
		return nil
	}

	pushParams := PushMessageParams{
		Cmd: data.Cmd,
		Uid: data.Uid,
		//Data:     data.Data.(MsgItem),
		Data:     data.Data,
		ChatInfo: data.ChatInfo,
		Devices:  offlineDeviceIds,
		Ak:       data.Ak,
	}
	return ProduceMessageToOfflineMq(ctx, pushParams)
}

func ProduceMessageToOfflineMq(ctx context.Context, pushParams PushMessageParams) error {
	p, _ := json.Marshal(pushParams)
	//return kafka.ProduceMsg(ctx, TopicMsgToOffline, p)
	fmt.Println("p", p)
	return nil
}
