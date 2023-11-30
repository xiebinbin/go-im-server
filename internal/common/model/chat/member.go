package chat

import (
	"context"
	"fmt"
	"imsdk/internal/common/dao/chat"
	"imsdk/internal/common/dao/chat/changelogs"
	"imsdk/internal/common/dao/chat/members"
	"imsdk/internal/common/model/active"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
)

type ChangeRoleParams struct {
	ChatId string           `json:"chat_id"`
	Role   members.RoleType `json:"role"`
	UIds   []string         `json:"UIds"`
}

type GetRoleCountParams struct {
	ChatId string           `json:"chat_id"`
	Role   members.RoleType `json:"role"`
}

type GetMembersParams struct {
	ChatId string           `json:"id"`
	Role   members.RoleType `json:"role"`
}

type MemberInfoParams struct {
	ChatId string   `json:"id"`
	UIds   []string `json:"uids"`
}

type ChangeStatusParams struct {
	ChatId string `json:"chat_id"`
	Status int8   `json:"status"`
	Reason string `json:"reason"`
	UId    string `json:"uid"`
}

type RemoveMembersParams struct {
	ChatId string   `json:"chat_id"`
	UIds   []string `json:"uids"`
}

type GetMemberByAChatIdsRequest struct {
	AChatIDs []string `json:"achat_ids" binding:"required"`
}

func ChangeRole(params ChangeRoleParams) error {
	if params.Role == members.RoleOwner && GetRoleCount(GetRoleCountParams{
		ChatId: params.ChatId,
		Role:   members.RoleOwner,
	}) > 0 {
		// todo
		//return errno.Add("", )
	}
	err := members.New().UpdateRoleByUIds(params.UIds, params.ChatId, params.Role)
	if err != nil {
		return errno.Add("change role err", errno.SysErr)
		// todo
	}
	return nil
}

func GetRoleCount(params GetRoleCountParams) int64 {
	return members.New().GetRoleCount(params.ChatId, params.Role)
}

func GetChatMemberIds(params GetMembersParams) ([]string, error) {
	data, err := members.New().GetChatMemberUIds(params.ChatId)
	if err != nil {
		// todo
		return data, err
	}
	return data, err
}

func GetMyChatIds(ctx context.Context, uid string) ([]string, error) {
	data, err := members.New().GetMyChatIdList(uid)
	if err != nil {
		// todo
	}
	res := make([]string, 0)
	for _, datum := range data {
		res = append(res, datum.ChatId)
	}
	return res, nil
}

func GetMemberByAChatIds(ctx context.Context, params GetMemberByAChatIdsRequest) (map[string][]string, error) {
	ak := ctx.Value("ak").(string)
	chatIds, uids := make([]string, 0), make([]string, 0)
	for _, id := range params.AChatIDs {
		chatIds = append(chatIds, base.CreateChatId(id, ak))
	}
	aChatInfo, _ := chat.New().GetAChatIdByIds(chatIds)
	chatMembers, res := make(map[string][]string, 0), make(map[string][]string, 0)
	if len(aChatInfo) == 0 {
		return res, nil
	}
	//for _, v := range aChatInfo {
	//	chatIds = append(chatIds, v)
	//}
	data, err := members.New().GetChatsMemberInfo(chatIds, "uid,chat_id")
	for _, datum := range data {
		uids = append(uids, datum.UID)
		chatMembers[aChatInfo[datum.ChatId]] = append(chatMembers[aChatInfo[datum.ChatId]], datum.UID)
	}
	auids := make(map[string]string, 0)
	fmt.Println("chatMembers---", chatMembers)
	if err != nil {
		// todo
		return chatMembers, err
	}
	for k, v := range chatMembers {
		var tmp []string
		for _, i2 := range v {
			if auid, ok := auids[i2]; ok {
				tmp = append(tmp, auid)
			}
		}
		res[k] = tmp
	}
	return res, err
}

func GetMemberByIds(params GetMembersParams) ([]string, error) {
	data, err := members.New().GetChatMemberUIds(params.ChatId)
	if err != nil {
		// todo
		return data, err
	}
	return data, err
}

func GetChatMembers(params GetMembersParams) {
	members.New().GetChatMember(params.ChatId, "uid, role, create_time")
}

func GetMembersInfo(params MemberInfoParams) ([]members.ChatMembersInfoRes, error) {
	return members.New().GetMembersInfo(params.ChatId, params.UIds)
}

func GetChatsMemberIds(chatIds []string) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)
	fields := "uid,chat_id"
	data, err := members.New().GetChatsMemberInfo(chatIds, fields)
	if data == nil {
		return nil, err
	}
	gRes := make(map[string][]string, 0)
	for _, v := range data {
		gRes[v.ChatId] = append(gRes[v.ChatId], v.UID)
	}
	for _, v := range chatIds {
		uIds := gRes[v]
		if gRes[v] == nil {
			uIds = []string{}
		}
		tmp := map[string]interface{}{
			"id":   v,
			"uids": uIds,
		}
		res = append(res, tmp)
	}
	return res, nil
}

func ChangeMemberStatus(ctx context.Context, params ChangeStatusParams) error {
	memDao := members.New()
	upData := members.Members{
		ID:     memDao.GetId(params.UId, params.ChatId),
		UID:    params.UId,
		ChatId: params.ChatId,
		Status: params.Status,
	}
	err := memDao.UpsertOne(upData)
	logCtx := log.WithFields(ctx, map[string]string{"action": "ChangeMemberStatus"})
	if err != nil {
		log.Logger().Error(logCtx, "err:", err)
		return errno.Add("change member status err", errno.SysErr)
		// todo
	}
	return nil
}

func RemoveMember(ctx context.Context, params RemoveMembersParams) error {
	logCtx := log.WithFields(ctx, map[string]string{"action": "RemoveMember"})
	if len(params.UIds) == 0 {
		return nil
	}
	memDao := members.New()
	ids := make([]string, 0)
	for _, id := range params.UIds {
		ids = append(ids, memDao.GetId(id, params.ChatId))
		fmt.Println("ids----", ids, id, params.ChatId)
	}
	count, err := memDao.DeleteMany(ids)
	log.Logger().Info(logCtx, "params:", params, count)
	if err != nil {
		log.Logger().Error(logCtx, err)
		return errno.Add("RemoveMember err", errno.SysErr)
		// todo
	}
	if count > 0 {
		chat.New().UpTotal(params.ChatId, -int(count))
		log.Logger().Info(logCtx, "DelActiveUserByChatId:", params.UIds, params.ChatId)
		for _, id := range params.UIds {
			err = changelogs.New().DelMemberInfo(params.ChatId, id)
			if err == nil {
				active.DelActiveUserByChatId(params.ChatId, id)
			} else {
				log.Logger().Info(logCtx, "DelActiveUserByChatId err:", err)
				continue
			}
		}
	}
	return nil
}

func GetChatActiveMember(chatId string) []string {
	activeUIds := active.GetActiveUserByChatId(chatId)
	mems, _ := members.New().GetByGidAndUids(chatId, activeUIds, "uid")
	if len(mems) > 0 {
		mIds := make(map[string]struct{})
		for _, mem := range mems {
			mIds[mem.UID] = struct{}{}
		}
		aUids := make([]string, 0)
		for _, mUid := range activeUIds {
			if _, ok := mIds[mUid]; ok {
				aUids = append(aUids, mUid)
			}
		}
		activeUIds = aUids
	} else {
		activeUIds = make([]string, 0)
	}
	return activeUIds
}

func GetChatMemberUIds(logCtx context.Context, chatId string) ([]string, error) {
	chatAllUIDs, err := members.New().GetChatMemberUIds(chatId)
	if err != nil {
		return chatAllUIDs, err
	}
	log.Logger().Info(logCtx, "GetChatMemberUIds - chat all uids : ", chatAllUIDs)
	//chatActives := active.GetActiveUserByChatId(chatId)
	chatActives := GetChatActiveMember(chatId)
	log.Logger().Info(logCtx, "GetChatMemberUIds - chat active uids: ", chatActives)
	if len(chatActives) > 0 {
		chatActives = append(chatActives, chatAllUIDs...)
		return funcs.RemoveDuplicatesAndEmpty(chatActives), nil
	}
	return chatAllUIDs, nil
}

func GetChatMemberUIdsWithoutOther(chatId string, uIds []string) ([]string, error) {
	return members.New().GetChatMemberUIds(chatId)
}
