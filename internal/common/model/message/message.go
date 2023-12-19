package message

import (
	"context"
	"errors"
	"fmt"
	json "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao"
	"imsdk/internal/common/dao/chat"
	"imsdk/internal/common/dao/chat/members"
	"imsdk/internal/common/dao/message/detail"
	"imsdk/internal/common/dao/message/usermessage"
	"imsdk/internal/common/model/forward"
	"imsdk/internal/common/pkg"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/config"
	"imsdk/pkg/app"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
)

const (
	HandleUndo     = "undo"
	HandleUndoMany = "undo-many"
)

type SendMessageParams struct {
	Mid           string                 `json:"mid"`
	ChatId        string                 `json:"chat_id"`
	SenderId      string                 `json:"from_uid"`
	ReceiveIds    []string               `json:"receive_ids"`
	Status        int8                   `json:"status"`
	Type          int8                   `json:"type"`
	Content       base.JsonString        `json:"content"`
	Extra         base.JsonString        `json:"extra"`
	Offline       map[string]interface{} `json:"offline"`
	NoIncrUnread  bool                   `json:"no_incr_unread"`  // business message (server send)
	AppCtx        string                 `json:"app_ctx"`         // business message (server send)
	NoPushOffline bool                   `json:"no_push_offline"` // business message (server send)
}

type SendMessageParamsTemp struct {
	Mid        string          `json:"mid"`
	Type       int8            `json:"type" binding:"required"`
	Content    base.JsonString `json:"content" binding:"required"`
	CreateTime int64           `json:"create_time" binding:"required"`
}

type GetMessageListRequest struct {
	Limit    int64  `json:"limit"`
	ChatId   string `json:"chat_id"`
	Sequence int64  `json:"sequence"`
}

type BatchDeleteRequest struct {
	Ids []string `json:"ids" binding:"required"`
}

type UserMessageIdsParams struct {
	Sequence int64 `json:"sequence" binding:"gte=0"`
}

type SendResp struct {
	Sequence    int64  `json:"sequence"`
	ID          string `json:"id"`
	FromAddress string `json:"from_uid"`
	Content     string `json:"content"`
	Type        int8   `json:"type"`
	CreateTime  int64  `json:"create_time"`
}

func Send(ctx context.Context, params SendMessageParams) (SendResp, error) {
	resp := SendResp{Sequence: 0}
	deData := []byte(ctx.Value("data").(string))
	json.Unmarshal(deData, &params)
	//if err != nil {
	//	return SendResp{}, err
	//}
	senderId := ctx.Value("uid").(string)
	params.SenderId = senderId
	logCtx := log.WithFields(ctx, map[string]string{"action": "sendMessage"})
	chatInfo, err := chat.New().GetInfoById(params.ChatId, "_id,type,total")
	fmt.Println("send message ,chat info :", chatInfo, "err: ", err, "  ,", chatInfo.Type)
	unknownErr := errno.Add("unknown error ", errno.SysErr)
	if err != nil {
		log.Logger().Error(logCtx, "failed to query chat info , err: ", err, params.ChatId)
		return resp, unknownErr
	}
	msgParamB, _ := json.Marshal(params)
	log.Logger().Info(logCtx, "send message params-model : ", string(msgParamB))
	_, er := members.New().GetByUidAndGid(params.SenderId, params.ChatId, "")
	if er != nil {
		if errors.Is(er, mongo.ErrNilDocument) {
			log.Logger().Error(logCtx, "user does not in chat , err: ", err)
			return resp, unknownErr
		}
	}
	// save message detail
	msgItem, er1 := saveMsgDetail(ctx, params)
	if er1 != nil {
		return resp, err
	}

	// to whom ?
	// priority:
	// 1: designated receiveUIds
	// 2: if chat members amount lte base.SyncProcessMaxAmount, send to all synchronously (get receive uid)
	// 3: active member
	// 4: put other to mq
	reqId := ctx.Value(base.HeaderFieldReqId).(string)
	mqData := WaitReprocessMsgMqData{
		ChatInfo:      chatInfo,
		MsgItem:       msgItem,
		AppCtx:        params.AppCtx,
		ReqId:         reqId,
		NoPushOffline: params.NoPushOffline,
		NoIncrUnread:  params.NoIncrUnread,
		ReceiverUsers: params.ReceiveIds,
	}
	log.Logger().Info(logCtx, "send msg -mqData: ", mqData)
	receiverLen := len(params.ReceiveIds)
	pushedUIds, _ := members.New().GetChatMembers(params.ChatId)
	log.Logger().Info(logCtx, "send message model : ", chatInfo.Total, base.SyncProcessMaxAmount, receiverLen, "chatInfo :", chatInfo)

	for _, uid := range pushedUIds {
		mqData.Uid = uid
		sequence, er2 := SpreadMsgAndReqPush(ctx, mqData, false)
		if er2 != nil {
			return resp, er1
		}
		if uid == senderId {
			resp.Sequence = sequence
		}
	}
	resp.Content = msgItem.Content
	resp.ID = msgItem.ID
	resp.Sequence = msgItem.Number
	resp.FromAddress = msgItem.SenderId
	resp.Type = msgItem.Type
	resp.CreateTime = msgItem.CreatedAt
	return resp, nil
}

func dealReadStatus(ctx context.Context, senderId string, msgContentStr string, chatInfo chat.Chat) bool {
	tmpRes := make(map[string]interface{})
	err := json.Unmarshal([]byte(msgContentStr), &tmpRes)
	logCtx := log.WithFields(ctx, map[string]string{"action": "dealReadStatus"})
	if err == nil {
		mIds := make([]string, 0)
		for _, v := range tmpRes["mids"].([]interface{}) {
			mIds = append(mIds, v.(string))
		}
		_, upErr := usermessage.New().UpdateIsRead(senderId, mIds, usermessage.IsReadYes)
		if upErr != nil {
			log.Logger().Error(logCtx, "update read info err", err)
		}
	}
	return true
}

func saveMsgDetail(ctx context.Context, params SendMessageParams) (MsgItem, error) {
	t := funcs.GetMillis()
	status := params.Status
	if params.Mid == "" {
		params.Mid = funcs.CreateMsgId(params.SenderId)
	}
	if params.Type == 0 {
		params.Type = detail.TypeSystem
	}
	number, err := detail.GetSequence(base.MsgNumberCountersKey)
	addData := detail.Detail{
		ID:         params.Mid,
		ChatId:     params.ChatId,
		FromUID:    params.SenderId,
		Content:    params.Content,
		Extra:      params.Extra,
		Type:       params.Type,
		Sequence:   number,
		ReceiveIds: params.ReceiveIds,
		Status:     status,
		CreatedAt:  t,
		UpdatedAt:  t,
	}
	detailObj := detail.New()
	addData.Status, err = detailObj.Save(addData)
	if !dao.DataIsSaveSuccessfully(err) {
		return MsgItem{}, errno.Add("failed to save msg detail", errno.SaveDataFailed)
	}
	msgItem := MsgItem{
		ID:         addData.ID,
		ChatId:     addData.ChatId,
		SenderId:   addData.FromUID,
		Content:    addData.Content,
		Extra:      addData.Extra,
		Action:     addData.Action,
		Status:     addData.Status,
		Number:     number,
		Type:       addData.Type,
		CreatedAt:  addData.CreatedAt,
		ReceiveIds: addData.ReceiveIds,
		UpdatedAt:  addData.UpdatedAt,
		Offline:    params.Offline,
	}

	return msgItem, nil
}

func GetMsgInfo(ctx context.Context, msgIds []string) []detail.Detail {
	data := detail.New().GetDetails(msgIds)
	res := make([]detail.Detail, 0)
	ak, _ := config.GetConfigAk()
	if ak == app.OfficialAK {
		for _, datum := range data {
			res = append(res, datum)
		}
	}
	return data
}

func GetMessageList(ctx context.Context, request GetMessageListRequest) []detail.GetMessageListResponse {
	limit := request.Limit
	where := bson.M{}
	if request.Sequence != 0 {
		where = bson.M{
			"sequence": bson.M{"$lt": request.Sequence},
		}
	}
	data := detail.New().GetListByLimit(limit, where)
	return data
}

func BatchDelete(ctx context.Context, ids []string) error {
	_, err := detail.New().DeleteByIds(ids)
	if err == nil {
		return err
	}
	return nil
}

func DeleteByUID(ctx context.Context, uid string) error {
	_, err := detail.New().DeleteByUID(uid)
	if err == nil {
		return err
	}
	return nil
}

func DeleteSelfByIds(uid string, ids []string) bool {
	_, err := usermessage.New().Delete(uid, ids)
	if err == nil {
		return true
	}
	return false
}

func DeleteSelfByChatIds(uid string, chatIds, exceptIds []string) bool {
	msgInfo := detail.New().GetExceptByChatIds(chatIds, exceptIds, "_id")
	var msgIds []string
	if len(msgInfo) > 0 {
		for _, v := range msgInfo {
			msgIds = append(msgIds, v.ID)
		}
		_, err := usermessage.New().Delete(uid, msgIds)
		if err != nil {
			return false
		}
	}

	if len(chatIds) > 0 {
		chatId := chatIds[0]
		_ = AddPrepareMsgByStatus(context.Background(), AddPrepareMsgByStatusReq{
			ChatId:   chatId,
			SenderId: uid,
			MIds:     exceptIds,
			Tag:      TagDelete,
		})
	}
	return true
}

const (
	TagDelete = 1
	TagRevoke = 2
)

type AddPrepareMsgByStatusReq struct {
	ChatId   string   `json:"chat_id"`
	SenderId string   `json:"sender_id"`
	MIds     []string `json:"mids"`
	Status   int8     `json:"status"`
	Tag      int8     `json:"tag"`
}

func AddPrepareMsgByStatus(ctx context.Context, params AddPrepareMsgByStatusReq) error {
	logCtx := log.WithFields(context.Background(), map[string]string{"action": "DeleteChat"})
	//log.Logger().Info(logCtx, "--2: ", params.MIds, params)
	if len(params.MIds) == 0 {
		return nil
	}
	t := funcs.GetMillis()
	chatId := params.ChatId
	if params.Tag == TagRevoke {
		params.Status = detail.StatusDelete
	}
	baseMsg := detail.Detail{
		ChatId:    chatId,
		FromUID:   params.SenderId,
		Status:    params.Status,
		CreatedAt: t,
		UpdatedAt: t,
	}
	detailDao := detail.New()
	userMessageDao := usermessage.New()
	//log.Logger().Info(logCtx, "--3: ", baseMsg)
	for _, id := range params.MIds {
		baseMsg.ID = id
		_, err := detailDao.Save(baseMsg)
		userMessageDao.Save(params.SenderId, id, chatId, usermessage.IsReadYes, usermessage.StatusDelete)
		if err != nil {
			log.Logger().Error(logCtx, "SpreadMsgAndReqPush: ", baseMsg, err)
			continue
		}
	}
	return nil
}

func DeleteSelfAll(uid string) bool {
	_, err := usermessage.New().DeleteSelfAll(uid)
	if err == nil {
		return true
	}
	return false
}

func SpreadMsgAndReqPush(ctx context.Context, data WaitReprocessMsgMqData, isMq bool) (int64, error) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "SpreadMsgToUser"})
	senderDeviceId := pkg.GetDeviceId(ctx)
	msgItem := data.MsgItem
	senderId := msgItem.SenderId
	msgItem.Sequence = 0

	isRead := usermessage.IsReadNo
	if data.Uid == senderId {
		isRead = usermessage.IsReadYes
	}
	userMsgDao := usermessage.New()
	sequence, err := userMsgDao.Save(data.Uid, msgItem.ID, msgItem.ChatId, int8(isRead), msgItem.Status)
	log.Logger().Info(logCtx, "SpreadMsgAndReqPush-1: ", data.Uid, msgItem.ID, msgItem.ChatId, sequence)
	defErr := errno.Add("failed to save user message", errno.DefErr)
	if err != nil {
		log.Logger().Error(logCtx, "uid: ", data.Uid, " failed to save user message , err: ", err)
		return 0, defErr
	}
	msgItem.Sequence = sequence

	// forward message to with socket
	senderInfo := forward.SenderInfo{
		SenderUid:      senderId,
		SenderDeviceId: senderDeviceId,
	}

	pushParams := forward.PushMessageParams{
		Cmd:           forward.CmdNoticeNewMsg,
		Uid:           data.Uid,
		Data:          msgItem,
		ChatInfo:      data.ChatInfo,
		SenderInfo:    senderInfo,
		Ak:            data.Ak,
		ReqId:         data.ReqId,
		AppCtx:        data.AppCtx,
		NoPushOffline: data.NoPushOffline,
	}
	if err := forward.PushMessageToUserSocketDirectly(ctx, pushParams, isMq); err != nil {
		return 0, err
	}
	return 0, nil
}
