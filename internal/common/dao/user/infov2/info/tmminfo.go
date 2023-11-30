package info

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/errno"
)

type TmmInfo struct {
	Base
}

func NewTmm() *TmmInfo {
	return new(TmmInfo)
}

func (u *TmmInfo) TableName() string {
	return "user_info"
}

// Init : init user table (create index: auid)
func (u *TmmInfo) Init() error {
	isUnique := true
	opts := options.IndexOptions{Unique: &isUnique}
	indexModel := mongoDriver.IndexModel{
		Keys: bson.D{
			{"auid", -1},
		},
		Options: &opts,
	}

	ctx := context.Background()
	_, err := u.collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (u *TmmInfo) collection() *mongo.CollectionInfo {
	return mongo.Database().SetTable(u.TableName())
}

func (u *TmmInfo) GetAUId(uId string) (string, error) {
	return uId, nil
}

func (u *TmmInfo) GetAUIds(uIds []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, id := range uIds {
		result[id] = id
	}
	return result, nil
}

func (u *TmmInfo) GetInfoById(uid string) (UserInfo, error) {
	var data UserInfo
	err := u.collection().Where(bson.M{"_id": uid}).FindOne(&data)
	return data, err
}

func (u *TmmInfo) GetInfoByPhone(phone string) (UserInfo, error) {
	var data UserInfo
	err := u.collection().Where(bson.M{"phone": phone}).FindOne(&data)
	return data, err
}

func (u *TmmInfo) GetUidByAUIds(auid []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, id := range auid {
		result[id] = id
	}
	return result, nil
}

func (u *TmmInfo) GetSliceAUIds(uIds []string) ([]string, error) {
	result := make([]string, 0)
	for _, id := range uIds {
		result = append(result, id)
	}
	return result, nil
}

func (u *TmmInfo) GetListByUIds(uids []string) ([]UserInfo, error) {
	data := make([]UserInfo, 0)
	for _, uid := range uids {
		data = append(data, UserInfo{
			ID:   uid,
			AUId: uid,
		})
	}
	return data, nil
}

func (u *TmmInfo) GetAll() ([]UserInfo, error) {
	var data []UserInfo
	err := u.collection().Where(bson.M{"status": StatusNormal}).FindMany(&data)
	return data, err
}

func (u *TmmInfo) VerifyUserStatus(uid string) error {
	info, _ := u.GetInfoById(uid)
	if info.Status == StatusDelete {
		return errno.Add("user-status-delete", errno.UserDelete)
	}
	if info.Status == StatusForbid {
		return errno.Add("user-status-forbid", errno.UserUnavailable)
	}
	return nil
}

func (u *TmmInfo) GetAllUIds(filters []Filter) ([]string, error) {
	where, err := formatWhere(filters)
	if err != nil {
		return nil, err
	}
	var data []UserInfo
	err = u.collection().Where(where).Fields("_id").FindMany(&data)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, v := range data {
		result = append(result, v.ID)
	}
	return result, err
}
