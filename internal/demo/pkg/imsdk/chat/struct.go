package chat

import "imsdk/internal/demo/pkg/imsdk/resource"

type TypeChat int8
type TypeStatus uint8

const (
	Normal  TypeChat = 1
	Notice           = 2
	Payment          = 3

	StatusNormal TypeStatus = 1
	StatusBlock             = 2
	StatusDelete            = 3
)

type createChat struct {
	AChatId string          `json:"achat_id" binding:"required"`
	AUIds   []string        `json:"auids"`
	Creator string          `json:"creator"`
	Name    string          `json:"name"`
	Type    TypeChat        `json:"type"`
	Avatar  *resource.Image `json:"avatar"`
}

type joinChat struct {
	AChatId string   `json:"achat_id" binding:"required"`
	AUIds   []string `json:"auids"`
}

type removeChatMember struct {
	AChatId string   `json:"achat_id" binding:"required"`
	AUIds   []string `json:"auids"`
}

type updateChatName struct {
	AChatId string `json:"achat_id"`
	Name    string `json:"name"`
}

type changeMemberStatus struct {
	AChatId string     `json:"achat_id"`
	AUId    string     `json:"auid"`
	Status  TypeStatus `json:"status"`
	Reason  string     `json:"reason"`
}

type updateChatAvatar struct {
	AChatId string          `json:"achat_id" binding:"required"`
	Avatar  *resource.Image `json:"avatar"`
	//Avatar  Avatar `json:"avatar"`
}
type Avatar struct {
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	BucketId string `json:"bucketId"`
	FileType string `json:"file_type"`
	Text     string `json:"text"`
}
