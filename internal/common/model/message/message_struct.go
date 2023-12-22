package message

import "imsdk/internal/common/pkg/base"

const (
	HandleUndo     = "undo"
	HandleUndoMany = "undo-many"
	DirectionUp    = "up"
	DirectionDown  = "down"
)

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
	Limit     int64  `json:"limit"`
	ChatId    string `json:"chat_id"`
	Sequence  int64  `json:"sequence"`
	Direction string `json:"direction"`
}

type BatchDeleteRequest struct {
	Ids []string `json:"ids" binding:"required"`
}

type RevokeBatchRequest struct {
	Ids []string `json:"ids" binding:"required"`
}

type RevokeByChatIdsRequest struct {
	ChatIds []string `json:"chat_ids" binding:"required"`
}

type DeleteByUIDRequest struct {
	UID string `json:"uid"`
}

type UserMessageIdsParams struct {
	Sequence int64 `json:"sequence" binding:"gte=0"`
}

type DeleteByChatIdsRequest struct {
	ChatIds []string `json:"chat_ids" binding:"required"`
}

type SendResp struct {
	Sequence    int64  `json:"sequence"`
	ID          string `json:"id"`
	FromAddress string `json:"from_uid"`
	Content     string `json:"content"`
	Type        int8   `json:"type"`
	CreateTime  int64  `json:"create_time"`
}
