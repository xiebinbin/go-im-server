package message

import (
	"imsdk/internal/common/dao/message/usermessage/unreadstock"
)

type UnreadStockParams struct {
	Sequence int64  `json:"sequence"`
	Num      uint16 `json:"num"`
}

func UpdateUnreadStockInfo(uid string, params UnreadStockParams) (int64, error) {
	modifiedCount, err := unreadstock.New().Update(uid, params.Sequence, params.Num)
	return modifiedCount, err
}
