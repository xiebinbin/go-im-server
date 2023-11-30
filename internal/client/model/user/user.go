package user

import (
	"context"
	"imsdk/internal/common/dao/user"
	user2 "imsdk/internal/common/model/user"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/redis"
	"time"
)

func GetUserInfo(ctx context.Context, response GetUserInfoRequest) ([]user.ListResponse, error) {
	data, err := user.New().GetByIDs(response.Ids)
	if err != nil {
		return data, errno.Add("unknown error", errno.DefErr)
	}
	return data, err
}

func GetConnectInfo(ctx context.Context, params user2.GetConnectInfoRequest) ([]GetConnectInfoResponse, error) {
	res := make([]GetConnectInfoResponse, 0)
	ak := ctx.Value("ak").(string)
	if err := verifyGetConnectInfo(ctx, params); err != nil {
		return res, err
	}
	uids, auidReflectUid := make([]string, 0), make(map[string]string, 0)
	for _, id := range params.AUIds {
		uid := base.CreateUId(id, ak)
		auidReflectUid[id] = uid
		uids = append(uids, uid)
	}
	params.UIds = uids
	connData, _ := user2.GetConnectInfo(ctx, params)
	for _, i := range params.AUIds {
		res = append(res, GetConnectInfoResponse{
			AUId:     i,
			IsOnline: connData[auidReflectUid[i]],
		})
	}
	return res, nil
}

func verifyGetConnectInfo(ctx context.Context, params user2.GetConnectInfoRequest) error {
	if len(params.UIds) > user2.ConnectInfoLimit {
		return errno.Add("", user2.ErrorExceedTimesLimit)
	}
	uid := ctx.Value(base.HeaderFieldUID).(string)
	cacheTag := "sdk:connect:info:" + uid
	if v, _ := redis.Client.Get(cacheTag).Result(); v == "t" {
		return errno.Add("", user2.ErrorExceedQuantityLimit)
	}
	redis.Client.SetNX(cacheTag, "t", time.Second*10)
	return nil
}

func UpdateName(ctx context.Context, request UpdateNameRequest) error {
	uid := ctx.Value(base.HeaderFieldUID).(string)
	err := user.New().Update(uid, user.User{
		Name: request.Name,
	})
	if err != nil {
		return errno.Add("sys-err", errno.SysErr)
	}
	return nil
}

func UpdateAvatar(ctx context.Context, request UpdateAvatarRequest) error {
	uid := ctx.Value(base.HeaderFieldUID).(string)
	err := user.New().Update(uid, user.User{
		Avatar: request.Avatar,
	})
	if err != nil {
		return errno.Add("sys-err", errno.SysErr)
	}
	return nil
}
