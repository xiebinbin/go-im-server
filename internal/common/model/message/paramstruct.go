package message

import "imsdk/internal/common/model/forward"

const (
	QueueDataTypeOrigin           = 1
	QueueDataTypeRebuild          = 2
	AsyncNum                      = 100
	PushTextType            int64 = 10020
	MsgContentTypeHollow          = 18
	MsgContentTypeDelHollow       = 81
	MsgContentTypeReport          = 82 // read report

	MsgTypeText       = 1
	MsgTypeImage      = 2
	MsgTypeVoice      = 3
	MsgTypeVideo      = 4
	MsgTypeAttachment = 5
	MsgTypeRedPacket  = 6
	MsgTypeTransfer   = 7
	MsgTypeCalls      = 8
	MsgTypePay        = 11
	MsgTypeCmd        = 18
)

var NotAddUnreadCountType = []int{MsgContentTypeHollow, MsgContentTypeDelHollow, MsgContentTypeReport}

type QueueData struct {
	QueueType int    `json:"type,omitempty"`
	Data      []byte `json:"data,omitempty"`
}

type MsgItem = forward.MsgItem
