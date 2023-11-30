package chat

import (
	"context"
	"fmt"
	"imsdk/internal/common/dao/chat"
	"imsdk/pkg/errno"
)

type UpdateAvatarParams struct {
	ChatId string `json:"id" binding:"required"`
	Avatar string `json:"avatar"`
}

func UpdateAvatar(ctx context.Context, params UpdateAvatarParams) error {
	uData := chat.Chat{
		Avatar: params.Avatar,
	}
	if err := chat.New().UpByID(params.ChatId, uData); err != nil {
		fmt.Println("UpdateAvatar err :", err, params)
		return errno.Add("fail", errno.DefErr)
	}
	return nil
}
