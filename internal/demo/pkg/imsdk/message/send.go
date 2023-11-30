package message

import (
	"context"
	"encoding/json"
	"imsdk/internal/demo/pkg/errno"
	"imsdk/internal/demo/pkg/imsdk/message/msgtype"
	"imsdk/internal/demo/pkg/imsdk/utils"
)

type SendBase struct {
	AChatId    string              `json:"achat_id"`
	AMId       string              `json:"amid"`
	SenderId   string              `json:"sender_id"`
	Content    string              `json:"content"`
	ReceiveIds []string            `json:"receive_ids"`
	SendTime   int64               `json:"send_time"`
	Type       msgtype.MessageType `json:"type"`
}
type SendParams struct {
	*SendBase
	Extra     string                 `json:"extra"`
	Action    map[string]interface{} `json:"action"`
	Offline   map[string]interface{} `json:"offline"`
	ExtraInfo map[string]interface{} `json:"extra_info"`
}

type SendMessageRespData struct {
	ErrCode  int    `json:"err_code"`
	ErrMsg   string `json:"err_msg"`
	Sequence int64  `json:"sequence"`
	IsSuc    int8   `json:"is_suc"`
	Reason   string `json:"reason"`
}

type SendResp struct {
	Sequence     int64  `json:"sequence"`
	IsSuc        int    `json:"is_suc"`
	FailedReason string `json:"reason"`
}

const (
	MiddleMsgSendUID = "33b1hvbe9lxf"
)

const (
	ErrCodeBlockUser = 200004
)

func NewSendBase(aChatId, aMId, content, senderId string, messageType msgtype.MessageType, receiveIds ...string) *SendBase {
	return &SendBase{
		AChatId:    aChatId,
		AMId:       aMId,
		Content:    content,
		SenderId:   senderId,
		SendTime:   utils.GetMillis(),
		Type:       messageType,
		ReceiveIds: receiveIds,
	}
}

func (c *Client) SendMessage(ctx context.Context, params SendParams) (int64, error) {
	res, err := c.sendMessage(ctx, params)
	if err != nil {
		return 0, err
	}
	var sendRes SendMessageRespData
	itemsByte, _ := json.Marshal(res.Data)
	err = json.Unmarshal(itemsByte, &sendRes)
	if err != nil {
		return 0, err
	}
	if sendRes.ErrCode != 0 {
		return 0, errno.Add(sendRes.ErrMsg, sendRes.ErrCode)
	}

	if sendRes.IsSuc != 1 {
		if sendRes.Reason == "block-user" {
			return 0, errno.Add(sendRes.Reason, ErrCodeBlockUser)
		}
		return 0, errno.Add("err", sendRes.ErrCode)
	}
	return sendRes.Sequence, nil
}
