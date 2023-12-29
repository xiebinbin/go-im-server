package user

import (
	"context"
	"errors"
	"fmt"
	"imsdk/internal/common/dao/user"
	"imsdk/internal/common/dao/user/infov2/info"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/config"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"

	"go.mongodb.org/mongo-driver/mongo"
)

type GetUIdsParams struct {
	Ids []string `json:"ids"`
}

type GetAuthParams struct {
	AK       string `json:"ak" binding:"required"`
	AUId     string `json:"auid" binding:"required"`
	AuthCode string `json:"authcode" binding:"required"`
	Env      string `json:"env"`
}

type AuthCodeParams struct {
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Ver       int64  `json:"ver"`
	Signature string `json:"signature"`
}

type GetAuthResponse struct {
	Id           string `json:"id"`
	Token        string `json:"token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func GetPubKey(ctx context.Context) string {
	pk, _ := config.GetConfigAk()
	return pk
}

func RegUser(ctx context.Context) (user.User, error) {
	address := ctx.Value(base.HeaderFieldUID).(string)
	pubKey := ctx.Value(base.HeaderFieldPubKey).(string)
	addData := user.User{
		ID:        address,
		PubKey:    pubKey,
		Avatar:    "https://img1.baidu.com/it/u=3709586903,1286591012&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500",
		Name:      funcs.GetRandString(6),
		Gender:    2,
		CreatedAt: funcs.GetMillis(),
	}
	//err := user.New().Init()
	err := user.New().Add(addData)
	if err != nil {
		return user.User{}, errno.Add("sys-err", errno.SysErr)
	}
	return addData, nil
}

func GetUserErr(uid string) error {
	uInfo, _ := user.New().GetInfoById(uid)
	fmt.Println("uInfo:", uInfo, uInfo.Status)
	if uInfo.Status == info.StatusDelete {
		return errno.Add("user-status-delete", errno.UserDelete)
	} else if uInfo.Status == info.StatusForbid {
		return errno.Add("user-status-forbid", errno.UserUnavailable)
	}
	return nil
}

func IsRegister(ctx context.Context) (user.User, error) {
	uid := ctx.Value(base.HeaderFieldUID).(string)
	uInfo, err := user.New().GetInfoById(uid)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, errno.Add("sys-err", errno.SysErr)
	}
	return uInfo, nil
}
