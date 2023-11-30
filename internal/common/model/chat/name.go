package chat

import (
	"context"
	"imsdk/internal/common/dao/chat"
	"imsdk/pkg/errno"
	"imsdk/pkg/log"
)

type UpdateNameParams struct {
	ChatId string `json:"chat_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

func UpdateName(ctx context.Context, params UpdateNameParams) error {
	logCtx := log.WithFields(ctx, map[string]string{"action": "UpdateChatName"})
	uData := map[string]interface{}{
		"name": params.Name,
	}
	if err := chat.New().UpMapByID(params.ChatId, uData); err != nil {
		log.Logger().Error(logCtx, err, params)
		return errno.Add("fail", errno.DefErr)
	}
	return nil

}
