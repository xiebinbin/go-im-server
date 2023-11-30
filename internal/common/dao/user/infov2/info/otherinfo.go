package info

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/errno"
)

type OtherInfo struct {
	Base
}

func NewOther() *OtherInfo {
	return new(OtherInfo)
}

func (u *OtherInfo) TableName() string {
	return "user"
}

func (u *OtherInfo) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(u.TableName(), rpMode)
}

// Init : init user table (create index: auid)
func (u *OtherInfo) Init() error {
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

func (u *OtherInfo) GetAUId(uid string) (string, error) {
	var data UserInfo
	err := u.collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": uid}).FindOne(&data)
	if err != nil {
		return "", err
	}
	return data.AUId, nil
}

func (u *OtherInfo) GetAUIds(uIds []string) (map[string]string, error) {
	var data []UserInfo
	result := make(map[string]string)
	err := u.collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": uIds}}).FindMany(&data)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		result[v.ID] = v.AUId
	}
	return result, nil
}

func (u *OtherInfo) GetSliceAUIds(uIds []string) ([]string, error) {
	var data []UserInfo
	result := make([]string, 0)
	err := u.collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": uIds}}).FindMany(&data)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		result = append(result, v.AUId)
	}
	return result, nil
}

func (u *OtherInfo) Add(data UserInfo) error {
	_, err := u.collection().InsertOne(&data)
	if err != nil {
		if dao.DataIsSaveSuccessfully(err) {
			return nil
		}
		return err
	}
	return nil
}

func (u *OtherInfo) Update(uid string, data UserInfo) error {
	_, err := u.collection().Where(bson.M{"_id": uid}).UpdateOne(&data)
	if err != nil {
		if dao.DataIsSaveSuccessfully(err) {
			return nil
		}
		return err
	}
	return nil
}

func (u *OtherInfo) GetInfoById(uid string) (UserInfo, error) {
	var data UserInfo
	err := u.collection().Where(bson.M{"_id": uid}).FindOne(&data)
	return data, err
}

func (u *OtherInfo) GetUidByAUIds(auid []string) (map[string]string, error) {
	var data []UserInfo
	result := make(map[string]string)
	err := u.collection(mongo.SecondaryPreferredMode).Where(bson.M{"auid": bson.M{"$in": auid}}).FindMany(&data)
	if err != nil {
		return result, err
	}
	for _, v := range data {
		result[v.AUId] = v.ID
	}
	return result, nil
}

func (u *OtherInfo) GetListByUIds(uids []string) ([]UserInfo, error) {
	data := make([]UserInfo, 0)
	err := u.collection().Where(bson.M{"_id": bson.M{"$in": uids}}).FindMany(&data)
	return data, err
}

func (u *OtherInfo) GetAll() ([]UserInfo, error) {
	var data []UserInfo
	err := u.collection().Where(bson.M{"status": StatusNormal}).FindMany(&data)
	return data, err
}

func (u *OtherInfo) VerifyUserStatus(uid string) error {
	info, _ := u.GetInfoById(uid)
	if info.Status == StatusDelete {
		return errno.Add("user-status-delete", errno.UserDelete)
	}
	if info.Status == StatusForbid {
		return errno.Add("user-status-forbid", errno.UserUnavailable)
	}
	return nil
}

func (u *OtherInfo) GetAllUIds(filters []Filter) ([]string, error) {
	where, err := formatWhere(filters)
	if err != nil {
		return nil, err
	}
	var data []UserInfo
	err = u.collection(mongo.SecondaryPreferredMode).Where(where).Fields("_id").FindMany(&data)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, v := range data {
		result = append(result, v.ID)
	}
	return result, err
}
