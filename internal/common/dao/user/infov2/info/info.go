package info

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/internal/common/dao/message/usermessage"
)

const (
	StatusNormal = 1
	StatusDelete = -2
	StatusForbid = -1
)

type UserInfo struct {
	ID        string `bson:"_id" json:"id"`
	AUId      string `bson:"auid" json:"auid,omitempty"`
	CreatedAt int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt int64  `bson:"update_time" json:"update_time,omitempty"`
	Status    int8   `bson:"status" json:"status"`
	//Avatar      map[string]interface{} `bson:"avatar" json:"avatar"`
	//PhonePrefix string                 `bson:"phone_prefix" json:"phone_prefix,omitempty"`
	//Phone       string                 `bson:"phone" json:"phone,omitempty"`
}

type Filter struct {
	Field string `json:"key"`
	Val   string `json:"value"`
}

type Info interface {
	GetAUId(uid string) (string, error)
	GetAUIds(uIds []string) (map[string]string, error)
	GetSliceAUIds(uIds []string) ([]string, error)
	Add(info UserInfo) error
	Update(uid string, info UserInfo) error
	GetInfoById(uid string) (UserInfo, error)
	GetInfoByPhone(phone string) (UserInfo, error)
	GetUidByAUIds(auid []string) (map[string]string, error)
	GetListByUIds(uid []string) ([]UserInfo, error)
	VerifyUserStatus(uid string) error
	GetAll() ([]UserInfo, error)
	GetAllUIds(filter []Filter) ([]string, error)
}

type Base struct {
	Info
}

func (u *Base) Add(data UserInfo) error {
	return nil
}

func (u *Base) Update(uid string, data UserInfo) error {
	return nil
}

func formatWhere(filters []Filter) (bson.M, error) {
	where := bson.M{}
	if len(filters) == 0 {
		return where, nil
	}
	for _, v := range filters {
		if v.Field == "phone_prefix" {
			where["phone_prefix"] = v.Val
		}
	}
	where["status"] = StatusNormal
	return where, nil
}

func Init(uid string) error {
	return usermessage.New().Init(uid)
}
