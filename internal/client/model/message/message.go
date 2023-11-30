package message

import (
	"context"
	_ "imsdk/internal/common/dao/common"
	_ "imsdk/internal/common/dao/conversation/usergroupsetting"
	_ "imsdk/internal/common/dao/conversation/usersinglesetting"
	"imsdk/internal/common/dao/message/usermessage"
	"imsdk/internal/common/dao/message/usermessage/maxsequence"
	_ "imsdk/internal/common/dao/user/setting"
	"imsdk/internal/common/pkg/base"
)

const (
	HandleUndo     = "undo"
	HandleUndoMany = "undo-many"
)

type SendMessageParams struct {
	Mid        string                 `json:"mid" binding:"required"`
	Type       uint16                 `json:"type" binding:"required"`
	ChatId     string                 `json:"chat_id" binding:"required,gt=3,lte=64"`
	Content    base.JsonString        `json:"content" binding:"required"`
	Action     map[string]interface{} `json:"action"`
	Extra      base.JsonString        `json:"extra"`
	SendTime   int64                  `json:"send_time" binding:"required"`
	SenderId   string                 `json:"sender_id"`
	ReceiveIds []string               `json:"receive_ids"`
	ExtraInfo  map[string]interface{} `json:"extra_info"`
}

type UserMessageInfoParams struct {
	Ids []string `json:"ids" binding:"required"`
}

type UserMessageIdsParams struct {
	Sequence int64 `json:"sequence" binding:"gte=0"`
}

func GetMsgIds(ctx context.Context, uid string, sequence int64) []usermessage.MsgIds {
	return usermessage.New().GetMsgIds(ctx, uid, sequence)
}

func GetMsgIdsV2(ctx context.Context, uid string, sequence int64) []usermessage.MsgIds {
	return usermessage.New().GetMsgIdsV2(ctx, uid, sequence)
}

func GetMessageStatus(uid string, ids []string) ([]usermessage.MsgIds, error) {
	data, err := usermessage.New().GetMsgStatus(uid, ids)
	if err != nil {
		return []usermessage.MsgIds{}, err
	}
	return data, nil
}

func GetMaxSequence(uid string) int64 {
	MaxSeqDao := maxsequence.New()
	return MaxSeqDao.GetSeq(uid).Seq
}
