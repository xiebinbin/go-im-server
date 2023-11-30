package message

import (
	"context"
	"imsdk/internal/common/pkg/config"
	"imsdk/pkg/funcs"
)

func HollowManSendMsg(ctx context.Context, params SendMessageParams, os string) error {
	hollowId, _ := config.GetHollowUId()
	mid := funcs.CreateMsgId(hollowId)
	if params.Mid != "" {
		mid = params.Mid
	}
	data := SendMessageParams{
		SenderId: hollowId,
		Mid:      mid,
		Type:     params.Type,
		ChatId:   params.ChatId,
		Content:  params.Content,
		Extra:    params.Extra,
	}
	_, err := Send(ctx, data)
	return err
}
