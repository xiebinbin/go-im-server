package chat

import (
	"context"
	json "github.com/json-iterator/go"
	"imsdk/internal/demo/pkg/imsdk"
	"imsdk/internal/demo/pkg/imsdk/resource"
	"sync"
)

type Options struct {
	Model       imsdk.ModelType
	Credentials *imsdk.Credentials
}

type Client struct {
	options *Options
}

var (
	once   sync.Once
	client *Client
)

func NewClient(options *Options) *Client {
	once.Do(func() {
		client = &Client{
			options: options,
		}
	})
	return client
}

func (c *Client) CreateChat(ctx context.Context, aChatId string, aUIds []string, creator string) error {
	dataByte, _ := json.Marshal(createChat{
		AChatId: aChatId,
		AUIds:   aUIds,
		Type:    Normal,
		Creator: creator,
		//Avatar:  avatar,
	})
	data := NewRequest(ActionCreateChat, string(dataByte))
	_, err := RequestIMServer(ctx, c.options, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateNoticeChat(ctx context.Context, aChatId, creator, name string) error {
	dataByte, _ := json.Marshal(createChat{
		AChatId: aChatId,
		Name:    name,
		Type:    Notice,
		Creator: creator,
	})
	data := NewRequest(ActionCreateChat, string(dataByte))
	_, err := RequestIMServer(ctx, c.options, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) JoinChat(ctx context.Context, aChatId string, aUIds ...string) error {
	dataByte, _ := json.Marshal(joinChat{
		AChatId: aChatId,
		AUIds:   aUIds,
	})
	data := NewRequest(ActionJoinChat, string(dataByte))
	_, err := RequestIMServer(ctx, c.options, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) RemoveMember(ctx context.Context, aChatId string, aUIds ...string) error {
	dataByte, _ := json.Marshal(removeChatMember{
		AChatId: aChatId,
		AUIds:   aUIds,
	})
	data := NewRequest(ActionRemoveChatMember, string(dataByte))
	_, err := RequestIMServer(ctx, c.options, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateName(ctx context.Context, aChatId, name string) {
	dataByte, _ := json.Marshal(updateChatName{
		AChatId: aChatId,
		Name:    name,
	})
	data := NewRequest(ActionUpdateChatName, string(dataByte))
	_, err := RequestIMServer(ctx, c.options, data)
	if err != nil {
		return
	}
}

func (c *Client) UpdateAvatar(ctx context.Context, aChatId string, avatar *resource.Image) {
	dataByte, _ := json.Marshal(updateChatAvatar{
		AChatId: aChatId,
		Avatar:  avatar,
	})
	data := NewRequest(ActionUpdateChatAvatar, string(dataByte))
	_, err := RequestIMServer(ctx, c.options, data)
	if err != nil {
		return
	}
}

func (c *Client) ForbiddenSendMessage(ctx context.Context, aChatId, aUId string, status TypeStatus, reason string) {
	dataByte, _ := json.Marshal(changeMemberStatus{
		AChatId: aChatId,
		AUId:    aUId,
		Status:  status,
		Reason:  reason,
	})
	data := NewRequest(ActionChangeChatMemberStatus, string(dataByte))
	_, err := RequestIMServer(ctx, c.options, data)
	if err != nil {
		return
	}
}

func (c *Client) CancelForbidden(ctx context.Context, aChatId, auid string) {
	dataByte, _ := json.Marshal(changeMemberStatus{
		AChatId: aChatId,
		AUId:    auid,
		Status:  StatusNormal,
		Reason:  "-",
	})
	data := NewRequest(ActionChangeChatMemberStatus, string(dataByte))
	_, err := RequestIMServer(ctx, c.options, data)
	if err != nil {
		return
	}
}
