package chat

import (
	"context"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao/chat"
	"imsdk/internal/common/dao/chat/members"
	"imsdk/internal/common/dao/message/detail"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/sdk"
)

type CreateParams struct {
	Id       string    `json:"id"`
	Creator  string    `json:"creator"`
	OwnerUid string    `json:"owner_uid"`
	UIds     []string  `json:"address_list"`
	Type     chat.Type `json:"type"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar,omitempty"`
	Intro    string    `json:"intro,omitempty"`
}

type CreateResp struct {
	UIds   []string  `json:"address_list"`
	Type   chat.Type `json:"type"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar,omitempty"`
	Intro  string    `json:"intro,omitempty"`
}

type GetMyChatParams struct {
	UId string `json:"uid"`
}

type DeleteChatRequest struct {
	ChatIds []string `json:"chat_ids"`
	UId     string   `json:"uid"`
}

type GetChatParams struct {
	ChatIds []string `json:"chat_ids"`
}

type GetChatResponse struct {
	ChatIds []string `json:"ids"`
	UId     string   `json:"uid"`
}

type GetMyChatResponse struct {
	ID               string    `bson:"_id" json:"id"`
	GId              string    `bson:"gid" json:"group_id"`
	CreatorUId       string    `bson:"creator" json:"creator,omitempty"`
	OwnerUId         string    `bson:"owner" json:"owner,omitempty"`
	Name             string    `bson:"name" json:"name"`
	Avatar           string    `bson:"avatar" json:"avatar"`
	Type             chat.Type `bson:"type" json:"type,omitempty"`
	Total            int64     `bson:"total" json:"total,omitempty"`
	LastReadSequence int64     `bson:"last_read_sequence" json:"last_read_sequence"`
	LastSequence     int64     `bson:"last_sequence" json:"last_sequence"`
	LastTime         int64     `bson:"last_time" json:"last_time"`
	CreatedAt        int64     `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt        int64     `bson:"update_time" json:"update_time,omitempty"`
}

const (
	MemExceedMax = 100100
)

func GetMemberUIds() ([]string, error) {
	res := make([]string, 0)
	return res, nil
}

func DeleteChat(ctx context.Context, request DeleteChatRequest) error {
	// delete My chat and message
	_, err := detail.New().DeleteMyByChatIds(request.UId, request.ChatIds)
	if err != nil {
		return err
	}
	return nil
}

func Create(ctx context.Context, params CreateParams) (CreateResp, error) {
	// todo if group exist
	logCtx := log.WithFields(ctx, map[string]string{"action": "createChat"})
	var res CreateResp
	t := funcs.GetMillis()
	creator := params.Creator
	allMemIds := []string{creator}
	allMemIds = append(allMemIds, params.UIds...)
	allMemDupIds := funcs.RemoveDuplicatesAndEmpty(allMemIds)
	chatType := chat.TypeGroup
	if len(allMemDupIds) == 2 {
		chatType = chat.TypeSingle
		params.Id = funcs.CreateSingleChatId(allMemDupIds[0], allMemDupIds[1])
	}
	ownerUid := creator
	if params.OwnerUid != "" {
		ownerUid = params.OwnerUid
	}
	// save group detail data
	chatDetail := chat.Chat{
		ID:         params.Id,
		CreatorUId: creator,
		OwnerUId:   ownerUid,
		Name:       params.Name,
		Total:      int64(len(allMemDupIds)),
		Type:       chatType,
		Avatar:     params.Avatar,
		Status:     chat.StatusNormal,
		CreatedAt:  t,
		UpdatedAt:  t,
	}
	if err := chat.New().Save(chatDetail); err != nil && !mongoDriver.IsDuplicateKeyError(err) {
		log.Logger().Info(logCtx, "save group detail unsuccessfully, err: ", err)
		return res, errno.Add("failed to save chat", errno.Exception)
	}
	// save members of group data
	err := JoinChat(ctx, JoinChatParams{
		JoinType:  sdk.JoinGroupTypeInvite,
		InviteUID: creator,
		ChatId:    params.Id,
		Role:      members.RoleCommonMember,
		UIds:      allMemDupIds,
	})
	if err != nil {
		log.Logger().Error(logCtx, "join chat detail unsuccessfully, err: ", err)
		return res, errno.Add("failed to join chat", errno.Exception)
	}
	err = ChangeRole(ChangeRoleParams{
		UIds:   []string{creator},
		Role:   members.RoleOwner,
		ChatId: params.Id,
	})
	res = CreateResp{
		UIds: allMemDupIds,
	}
	return res, nil
}

func GetMyChat(ctx context.Context, params GetMyChatParams) ([]GetMyChatResponse, error) {
	ids, _ := GetMyChatIds(ctx, params.UId)
	data, _ := GetChatList(ctx, GetChatParams{
		ChatIds: ids,
	})
	var res []GetMyChatResponse
	if len(data) > 0 {
		for _, datum := range data {
			res = append(res, GetMyChatResponse{
				ID:               datum.ID,
				GId:              datum.GId,
				OwnerUId:         datum.OwnerUId,
				Avatar:           datum.Avatar,
				Name:             datum.Name,
				Type:             datum.Type,
				Total:            datum.Total,
				LastReadSequence: datum.LastReadSequence,
				LastSequence:     datum.LastSequence,
				LastTime:         datum.LastTime,
			})
		}
	}
	return res, nil
}

func GetIdListByUid(uid string) ([]string, error) {
	res := make([]string, 0)
	data, err := members.New().GetMyChatIdList(uid)
	if data == nil {
		return res, err
	}
	for _, v := range data {
		res = append(res, v.ChatId)
	}
	return res, err
}

func GetChatList(ctx context.Context, params GetChatParams) ([]chat.Chat, error) {
	//var err error
	res, err := chat.New().GetInfoByIds(params.ChatIds, "_id,name,avatar,type,last_read_number, last_number, last_time")
	if err != nil {
		return []chat.Chat{}, nil
	}
	return res, err
}
