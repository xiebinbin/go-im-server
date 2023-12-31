package group

import (
	"context"
	"fmt"
	"imsdk/internal/common/dao/group/detail"
	"imsdk/internal/common/dao/group/members"
	"imsdk/internal/common/dao/user"
)

type ApplyListResponse struct {
	ID        string `bson:"_id" json:"id"`
	GID       string `bson:"gid" json:"gid"`
	UID       string `bson:"uid" json:"uid"`
	Avatar    string `bson:"avatar" json:"avatar"`
	Name      string `bson:"name" json:"name"`
	EncKey    string `bson:"enc_key" json:"enc_key"`
	Role      uint8  `bson:"role" json:"role"`
	Status    int8   `bson:"status" json:"status"`
	CreatedAt int64  `bson:"created_at" json:"create_time"`
}

func MyApplyList(ctx context.Context, uid string, request IdsRequest) ([]members.ApplyRes, error) {
	dao := members.New()
	groupMembers, _ := dao.GetMyGroupByIds(uid, request.Ids)
	ids := make([]string, 0)
	memMap := make(map[string]members.ApplyRes)
	if len(groupMembers) > 0 {
		for _, v := range groupMembers {
			ids = append(ids, v.GID)
			memMap[v.GID] = v
		}
	}
	res := make([]members.ApplyRes, 0)
	items, err := detail.New().GetInfoByIds(ids)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return res, nil
	}
	gidSlice := make([]string, 0)
	for _, item := range items {
		gidSlice = append(gidSlice, item.ID)
	}
	groups, er := detail.New().GetInfoByIds(gidSlice)
	if er != nil || len(groups) == 0 {
		return res, er
	}
	for _, v := range groups {
		item := members.ApplyRes{
			ID:        memMap[v.ID].ID,
			GID:       v.ID,
			Status:    memMap[v.ID].Status,
			Avatar:    v.Avatar,
			Name:      v.Name,
			EncKey:    memMap[v.ID].EncKey,
			CreatedAt: v.CreatedAt,
		}
		res = append(res, item)
	}
	return res, err
}

func ApplyList(ctx context.Context, uid string, request IdsRequest) ([]ApplyListResponse, error) {
	dao := members.New()
	groupMembers, _ := dao.GetListsByStatusAndRole(uid, request.Ids, []int8{}, []int8{int8(members.RoleOwner), int8(members.RoleAdministrator)})
	fmt.Println("groupMembers: ", groupMembers)
	ids := make([]string, 0)
	if len(groupMembers) > 0 {
		for _, v := range groupMembers {
			ids = append(ids, v.GID)
		}
	}
	applyItems, _ := dao.GetApplyByGIds(ids, []int8{members.StatusIng, members.StatusRefuse})
	fmt.Println("applyItems: ", applyItems)
	res := make([]ApplyListResponse, 0)
	fmt.Println("items: ", ids, applyItems)
	if len(applyItems) == 0 {
		return []ApplyListResponse{}, nil
	}
	uidSlice := make([]string, 0)
	for _, item := range applyItems {
		uidSlice = append(uidSlice, item.UID)
	}
	memMap, err := user.New().GetInfoByIds(uidSlice)
	for _, v := range applyItems {
		item := ApplyListResponse{
			ID:        v.ID,
			GID:       v.GID,
			UID:       v.UID,
			Status:    v.Status,
			Avatar:    memMap[v.UID].Avatar,
			Name:      memMap[v.UID].Name,
			Role:      v.Role,
			EncKey:    v.EncKey,
			CreatedAt: v.CreatedAt,
		}
		res = append(res, item)
	}
	return res, err
}
