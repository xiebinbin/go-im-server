package chat

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"imsdk/internal/common/dao/chat"
	"imsdk/internal/common/pkg/config"
	"imsdk/internal/demo/model"
	"imsdk/internal/demo/pkg"
	"imsdk/internal/demo/pkg/imsdk"
	imsdkChat "imsdk/internal/demo/pkg/imsdk/chat"
	"imsdk/internal/demo/pkg/imsdk/resource"
	"imsdk/pkg/funcs"
)

type CreateParams struct {
	AChatId string                 `json:"achat_id" binding:"required"`
	AUIds   []string               `json:"auids"`
	Type    chat.Type              `json:"type"`
	Creator string                 `json:"creator"`
	Name    string                 `json:"name"`
	Avatar  map[string]interface{} `json:"avatar"`
}

type JoinParams struct {
	AChatId string   `json:"achat_id" binding:"required"`
	AUIds   []string `json:"auids"`
}

type RemoveParams struct {
	AChatId string   `json:"achat_id" binding:"required"`
	AUIds   []string `json:"auids"`
}

type GetChatMemberListParams struct {
	AChatIds []string `json:"achat_ids" binding:"required"`
}

func NewChatClient() *imsdkChat.Client {
	conf, _ := config.GetIMSdkKey()
	chatClient := imsdkChat.NewClient(&imsdkChat.Options{
		Credentials: imsdk.NewStaticCredentials(conf.AK, conf.SK),
		Model:       funcs.GetEnv(),
	})
	return chatClient
}

func CreateChat(ctx context.Context, params CreateParams) error {
	//params.AChatId = funcs.CreateSingleChatId("2x2qlr88wdcz", "1c51f241d167")
	//params.AUIds = []string{"2x2qlr88wdcz", "1c51f241d167"}
	fmt.Println("params.AChatId:", params.AChatId)
	err := NewChatClient().CreateChat(ctx, params.AChatId, params.AUIds, params.Creator)
	NewChatClient().UpdateName(ctx, params.AChatId, params.Name)
	NewChatClient().UpdateAvatar(ctx, params.AChatId, resource.NewOriginImage(params.Avatar))
	if err != nil {
		return err
	}
	return nil
}

func CreateNoticeChat(ctx context.Context, params CreateParams) error {
	err := NewChatClient().CreateNoticeChat(ctx, params.AChatId, params.Creator, params.Name)
	NewChatClient().UpdateAvatar(ctx, params.AChatId, resource.NewOriginImage(params.Avatar))
	if err != nil {
		return err
	}
	return nil
}

func JoinChat(ctx context.Context, params JoinParams) error {
	err := NewChatClient().JoinChat(ctx, params.AChatId, params.AUIds...)
	if err != nil {
		return err
	}
	return nil
}

func RemoveMember(ctx context.Context, params RemoveParams) error {
	err := NewChatClient().RemoveMember(ctx, params.AChatId, params.AUIds...)
	if err != nil {
		return err
	}
	return nil
}

func GetChatList(ctx *gin.Context) (pkg.CurlResponse, error) {
	params := ""
	dataByte, _ := json.Marshal(params)
	return model.RequestIMServer("getChatList", string(dataByte))
}

func GetMemberByAChatIds(ctx *gin.Context, aChatIds []string) (pkg.CurlResponse, error) {
	params := map[string]interface{}{
		"achat_ids": aChatIds,
	}
	dataByte, _ := json.Marshal(params)
	return model.RequestIMServer("getMemberByAChatIds", string(dataByte))
}

func GetChatMemberList(ctx *gin.Context) error {
	var params GetChatMemberListParams
	dataByte, _ := json.Marshal(params)
	model.RequestIMServer("getMemberByAChatIds", string(dataByte))
	return nil
}
