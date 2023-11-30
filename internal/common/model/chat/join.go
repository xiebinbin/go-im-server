package chat

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao/chat"
	"imsdk/internal/common/dao/chat/members"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
)

type JoinChatParams struct {
	JoinType  members.JoinType `json:"join_type"`
	InviteUID string           `json:"invite_uid"`
	ChatId    string           `json:"chat_id"`
	Role      members.RoleType `json:"role"`
	UIds      []string         `json:"UIds"`
}

func JoinChat(ctx context.Context, params JoinChatParams) error {
	logCtx := log.WithFields(ctx, map[string]string{"action": "JoinChat"})
	if len(params.UIds) == 0 {
		return nil
	}
	memDao := members.New()
	t := funcs.GetMillis()
	for _, u := range params.UIds {
		item := members.Members{
			ID:        memDao.GetId(u, params.ChatId),
			ChatId:    params.ChatId,
			UID:       u,
			Role:      params.Role,
			JoinType:  params.JoinType,
			InviteUID: params.InviteUID,
			Status:    members.StatusNormal,
			CreatedAt: t,
			UpdatedAt: t,
		}
		if err := memDao.UpsertOne(item); err != nil && !mongo.IsDuplicateKeyError(err) {
			log.Logger().Info(logCtx, "DuplicateKey JoinChat Error : ", err)
			return errno.Add("duplicateKey join to chat", errno.Exception)
		}
	}
	totalMember := memDao.GetChatMembersCount(params.ChatId)
	if totalMember > 2 {
		err := chat.New().UpMapByID(params.ChatId, map[string]interface{}{
			"only_two": 0,
			"total":    totalMember,
		})
		if err != nil {
			log.Logger().Info(logCtx, "update chat only two fail : ", err)
			return errno.Add("update chat only two fail", errno.Exception)
		}
	}
	//changelogs.New().UpdateManyMemberInfo(params.ChatId, params.UIds)
	return nil
}
