package detail

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"time"
)

type Detail struct {
	ID          string `bson:"_id" json:"id"`
	Creator     string `bson:"creator" json:"creator,omitempty"`
	Owner       string `bson:"owner" json:"owner,omitempty"`
	Name        string `bson:"name" json:"name,omitempty"`
	Notice      string `bson:"notice" json:"notice,omitempty"`
	Pub         string `bson:"pub" json:"pub,omitempty"`
	NoticeId    string `bson:"notice_md5" json:"notice_md5,omitempty"`
	Desc        string `bson:"desc" json:"desc,omitempty"`
	DescMd5     string `bson:"desc_md5" json:"desc_md5,omitempty"`
	Avatar      string `bson:"avatar" json:"avatar,omitempty"`
	Cover       string `bson:"cover" json:"cover,omitempty"`
	Status      int8   `bson:"status" json:"status,omitempty"`
	Total       int    `bson:"total" json:"total,omitempty"`
	MemberLimit int    `bson:"member_limit" json:"member_limit,omitempty"`
	CreatedAt   int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt   int64  `bson:"update_time" json:"update_time,omitempty"`
}

type ListResponse struct {
	ID          string `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Avatar      string `bson:"avatar" json:"avatar"`
	Notice      string `bson:"notice" json:"notice"`
	Desc        string `bson:"desc" json:"desc"`
	Total       int    `bson:"total" json:"total"`
	MemberLimit int    `bson:"member_limit" json:"member_limit"`
	Owner       string `bson:"owner" json:"owner"`
	NoticeMd5   string `bson:"notice_md5" json:"notice_md5"`
	Pub         string `bson:"pub" json:"pub"`
	CreatedAt   int64  `bson:"create_time" json:"create_time,omitempty"`
}

const (
	StatusDel       = -1
	StatusForbidden = -2
	StatusYes       = 1
	TotalMax        = 300
)

func (d Detail) TableName() string {
	return "group_detail"
}

func New() *Detail {
	return new(Detail)
}

func (d Detail) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(d.TableName(), rpMode)
}

func (d Detail) Add(data Detail) error {
	data.Owner = data.Creator
	data.Status = StatusYes
	_, err := d.Collection().InsertOne(&data)
	return err
}

func (d Detail) GetByID(id string, fields string) (Detail, error) {
	var data Detail
	cacheTag := "chat:groupInfo:" + id
	info := redis.Client.Get(cacheTag).Val()
	if info == "" {
		err := d.Collection(mongo.SecondaryPreferredMode).Fields(dao.GetMongoFieldsBsonByString(fields)).FindByID(id, &data)
		if err == nil {
			redis.Client.Set(cacheTag, data, time.Second*5)
		} else if errors.Is(err, mongoDriver.ErrNilDocument) {
			redis.Client.Set(cacheTag, "nil", time.Second)
		}
		return data, err
	}
	if info != "nil" {
		err := json.Unmarshal([]byte(info), &data)
		return data, err
	}
	return data, errors.New("group not exist")
}

func (d Detail) UpByID(id string, uData Detail) error {
	_, err := d.Collection().UpByID(id, uData)
	return err
}

func (d Detail) UpMapByID(id string, uData map[string]interface{}) error {
	_, err := d.Collection().UpByID(id, uData)
	return err
}

func (d Detail) GetInfoByIds(ids []string) ([]ListResponse, error) {
	var data []ListResponse
	d.Collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": ids}}).FindMany(&data)
	return data, nil
}

func (d Detail) UpTotal(id string, num int) bool {
	uData := bson.M{"$inc": bson.M{"total": num}}
	err := d.Collection().Where(bson.M{"_id": id}).UpsertByBson(uData).Err()
	return err == nil
}

func (d Detail) OwnerQuit(gid, newOwner string) error {
	uData := bson.M{"owner_uid": newOwner, "update_time": funcs.GetMillis()}
	_, err := d.Collection().UpByID(gid, uData)
	return err
}

func (d Detail) GetOldGroup() []Detail {
	var data []Detail
	where := bson.M{"is_old": StatusYes}
	err := d.Collection().Where(where).FindMany(&data)
	if err != nil {
		return nil
	}
	return data
}

func (d Detail) Save(data Detail) error {
	_, err := d.Collection().InsertOne(&data)
	return err
}

func (d Detail) SendTransMsg(gid string) error {
	uData := bson.M{"is_send": 1}
	_, err := d.Collection().UpByID(gid, uData)
	return err
}

func (d Detail) GetListByLimit(limit, offset int64) []Detail {
	var data []Detail
	err := d.Collection().Sort(bson.M{"create_time": 1}).Limit(limit).Skip(offset).FindMany(&data)
	if err != nil {
		return nil
	}
	//Select("_id,avatar").Limit(limit).Offset(offset).Find(&data)
	return data
}

func (d Detail) GetCount() int64 {
	count := d.Collection().Count()
	return count
}

func (d Detail) GetAll() []Detail {
	var data []Detail
	err := d.Collection().Fields("_id").FindMany(&data)
	if err != nil {
		return nil
	}
	return data
}
