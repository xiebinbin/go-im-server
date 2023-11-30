package user

import (
	"context"
	"imsdk/internal/common/pkg/config"
	"imsdk/internal/demo/pkg/imsdk"
	"imsdk/internal/demo/pkg/imsdk/common/user"
	"imsdk/pkg/funcs"
)

type GetConnectInfoRequest struct {
	AUIds []string `json:"auids" `
}

func NewUserClient() *user.Client {
	conf, _ := config.GetIMSdkKey()
	userClient := user.NewClient(&user.Options{
		Credentials: imsdk.NewStaticCredentials(conf.AK, conf.SK),
		Model:       funcs.GetEnv(),
	})
	return userClient
}

func GetConnectInfo(ctx context.Context, request GetConnectInfoRequest) (interface{}, error) {
	res, _ := NewUserClient().GetConnectInfo(ctx, request.AUIds)
	return res, nil
}
