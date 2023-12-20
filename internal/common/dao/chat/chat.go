package chat

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/redis"
	"imsdk/pkg/sdk"
	"time"
)

type Type = int8
type Chat struct {
	ID               string         `bson:"_id" json:"id"`
	TargetId         string         `bson:"target_id" json:"target_id"`
	CreatorUId       string         `bson:"creator" json:"creator,omitempty"`
	OwnerUId         string         `bson:"owner" json:"owner,omitempty"`
	Name             string         `bson:"name" json:"name"`
	Avatar           string         `bson:"avatar" json:"avatar"`
	Type             Type           `bson:"type" json:"type,omitempty"`
	Status           sdk.StatusType `bson:"status" json:"status,omitempty"`
	Total            int64          `bson:"total" json:"total,omitempty"`
	LastReadSequence int64          `bson:"last_read_sequence" json:"last_read_sequence"`
	LastSequence     int64          `bson:"last_sequence" json:"last_sequence"`
	LastTime         int64          `bson:"last_time" json:"last_time"`
	CreatedAt        int64          `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt        int64          `bson:"update_time" json:"update_time,omitempty"`
}

const (
	StatusNormal      = 1
	TypeSingle   int8 = 1
	TypeGroup    int8 = 2
)

func New() *Chat {
	return new(Chat)
}

func (c *Chat) TableName() string {
	return "chat"
}

func (c *Chat) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(c.TableName(), rpMode)
}

func (c *Chat) Add(data Chat) error {
	data.Status = sdk.StatusNormal
	_, err := c.Collection().InsertOne(&data)
	return err
}

func (c *Chat) GetByID(id string, fields string) (Chat, error) {
	var data Chat
	cacheTag := "imsdk:chatInfo:" + id
	info := redis.Client.Get(cacheTag).Val()
	if info == "" {
		err := c.Collection(mongo.SecondaryPreferredMode).Fields(dao.GetMongoFieldsBsonByString(fields)).FindByID(id, &data)
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
	return data, errors.New("chat not exist")
}

func (c *Chat) GetByAChatID(achatId string, fields string) (Chat, error) {
	var data Chat
	cacheTag := "imsdk:chatInfo:" + achatId
	info := redis.Client.Get(cacheTag).Val()
	if info == "" {
		err := c.Collection().Fields(fields).Where(bson.M{"achat_id": achatId}).FindOne(&data)
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
	return data, errors.New("chat not exist")
}

func (c *Chat) UpByID(id string, uData Chat) error {
	_, err := c.Collection().Where(bson.M{"_id": id}).UpdateOne(&uData)
	return err
}

func (c *Chat) UpMapByID(id string, uData map[string]interface{}) error {
	_, err := c.Collection().UpByID(id, uData)
	return err
}

func (c *Chat) GetInfoByIds(ids []string, fields string) ([]Chat, error) {
	data := make([]Chat, 0)
	err := c.Collection().Where(bson.M{"_id": bson.M{"$in": ids}}).
		Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, err
}

func (c *Chat) GetAChatIdByIds(Ids []string) (map[string]string, error) {
	var data []Chat
	result := make(map[string]string)
	err := c.Collection().Where(bson.M{"_id": bson.M{"$in": Ids}}).FindMany(&data)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		result[v.ID] = v.TargetId
	}
	return result, nil
}

func (c *Chat) GetIdByAChatIds(AChatIds []string) (map[string]string, error) {
	var data []Chat
	result := make(map[string]string)
	err := c.Collection().Where(bson.M{"achat_id": bson.M{"$in": AChatIds}}).FindMany(&data)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		result[v.TargetId] = v.ID
	}
	return result, nil
}

func (c *Chat) GetInfoById(id, fields string) (Chat, error) {
	var data Chat
	cacheTag := "imsdk:chatInfo:" + id
	info := redis.Client.Get(cacheTag).Val()
	if info == "" {
		err := c.Collection().Where(bson.M{"_id": id}).FindOne(&data)
		if err == nil {
			redis.Client.Set(cacheTag, data, time.Second*3)
		} else if errors.Is(err, mongoDriver.ErrNilDocument) {
			redis.Client.Set(cacheTag, "nil", time.Second)
		}
		return data, err
	}
	if info != "nil" {
		err := json.Unmarshal([]byte(info), &data)
		return data, err
	}
	return data, errors.New("chat not exist")
}

func (c *Chat) UpTotal(id string, num int) bool {
	uData := bson.M{"$inc": bson.M{"total": num}}
	err := c.Collection().Where(bson.M{"_id": id}).UpsertByBson(uData).Err()
	return err == nil
}

func (c *Chat) Save(data Chat) error {
	_, err := c.Collection().InsertOne(&data)
	return err
}

func (c *Chat) Upsert(data Chat) error {
	err := c.Collection().Where(bson.M{"_id": data.ID}).Upsert(&data)
	return err.Err()
}

func (c *Chat) GetCount(where bson.M) int64 {
	count := c.Collection().Where(where).Count()
	return count
}

func (c *Chat) GetListByCond(where bson.M) ([]Chat, error) {
	data := make([]Chat, 0)
	err := c.Collection().Where(where).FindMany(&data)
	return data, err
}

func (c *Chat) GetListByLimit(limit, offset int64, where bson.M) []Chat {
	var data []Chat
	c.Collection(mongo.SecondaryPreferredMode).Where(where).Sort(bson.M{"create_time": -1}).Limit(limit).Skip(offset).FindMany(&data)
	//Select("_id,avatar").Limit(limit).Offset(offset).Find(&data)
	return data
}
