package group

import (
	"context"
	"imsdk/internal/common/dao/group/detail"
	"imsdk/internal/common/dao/group/members"
)

func MyApplyList(ctx context.Context, uid string, request IdsRequest) ([]members.ApplyRes, error) {
	dao := members.New()
	applyList, _ := dao.GetMyGroupByIds(uid, request.GroupIDs)
	ids := make([]string, 0)
	memMap := make(map[string]members.ApplyRes)
	if len(applyList) > 0 {
		for _, v := range applyList {
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

func ApplyList(ctx context.Context, uid string, request IdsRequest) ([]members.ApplyRes, error) {
	dao := members.New()
	groupMembers, _ := dao.GetMyGroupByIds(uid, request.GroupIDs)
	ids := make([]string, 0)
	memMap := make(map[string]members.Members)
	if len(groupMembers) > 0 {
		for _, v := range groupMembers {
			if v.Role == 1 || v.Role == 2 {
				ids = append(ids, v.GroupID)
				memMap[v.GroupID] = v
			}
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
