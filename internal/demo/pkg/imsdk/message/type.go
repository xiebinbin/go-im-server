package message

import (
	"context"
	"encoding/json"
	"imsdk/internal/demo/pkg/imsdk/message/msgtype"
)

func (c *Client) SendRedPacketMessage(ctx context.Context, aChatId, aMId string, text *msgtype.RedPack, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(text)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeRedPacket, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendTransferMessage(ctx context.Context, aChatId, aMId string, text *msgtype.Transfer, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(text)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeTransfer, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendMomentsMessage(ctx context.Context, aChatId, aMId string, text *msgtype.Moments, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(text)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeMoments, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}
