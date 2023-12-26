package group

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao/group/app"
	"imsdk/pkg/funcs"
)

type AddAppRequest struct {
	GroupID string `json:"gid" binding:"required"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Desc    string `json:"desc"`
	Url     string `json:"url"`
	Sort    int    `json:"sort"`
}

type UpdateAppRequest struct {
	ID      string `json:"id" binding:"required"`
	GroupID string `json:"gid"`  // 群id
	Name    string `json:"name"` // 名称
	Icon    string `json:"icon"` // 图标
	Desc    string `json:"desc"` // 描述
	Url     string `json:"url"`  // 地址
	Sort    int    `json:"sort"` // 排序
}

type DeleteByGIdsRequest struct { // 删除群组下的所有应用
	GIds []string `json:"gids" binding:"required"`
}

func AddApp(ctx context.Context, uid string, request AddAppRequest) error {
	t := funcs.GetMillis()
	err := app.New().Add(app.App{
		ID:        funcs.UniqueId16(),
		GroupID:   request.GroupID,
		Name:      request.Name,
		Icon:      request.Icon,
		Desc:      request.Desc,
		Url:       request.Url,
		Sort:      request.Sort,
		Status:    app.StatusYes,
		CreatedAt: t,
		UpdatedAt: t,
	})
	if err != nil {
		return err
	}
	return nil
}

func AppList(ctx context.Context, uid string, request IdsRequest) ([]app.App, error) {
	data, err := app.New().GetByIds(request.Ids, []int{app.StatusYes, app.StatusForbidden})
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return []app.App{}, err
	}
	return data, nil
}

func AppInfo(ctx context.Context, uid string, request IdsRequest) ([]app.App, error) {
	data, err := app.New().GetByIds(request.Ids, []int{app.StatusYes, app.StatusForbidden})
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return []app.App{}, err
	}
	return data, nil
}

func UpdateApp(ctx context.Context, uid string, request UpdateAppRequest) error {
	err := app.New().UpdateById(request.ID, app.App{GroupID: request.GroupID,
		Name:      request.Name,
		Icon:      request.Icon,
		Desc:      request.Desc,
		Url:       request.Url,
		Sort:      request.Sort,
		UpdatedAt: funcs.GetMillis(),
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteByIds(ctx context.Context, uid string, request IdsRequest) error {
	err := app.New().DeleteByIds(request.Ids)
	if err != nil {
		return err
	}
	return nil
}

func DeleteByGIds(ctx context.Context, uid string, request DeleteByGIdsRequest) error {
	err := app.New().DeleteByGroupIds(request.GIds)
	if err != nil {
		return err
	}
	return nil
}
