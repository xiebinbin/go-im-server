package phone

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
)

type Phone struct {
	ID          string `bson:"_id" json:"id"`
	PhonePrefix string `bson:"phone_prefix" json:"phone_prefix"`
	Phone       string `bson:"phone" json:"phone"`
	Status      int8   `bson:"status" json:"status"`
	CreatedAt   int64  `bson:"create_time" json:"create_time"`
	UpdatedAt   int64  `bson:"update_time" json:"update_time"`
}

func (p Phone) TableName() string {
	return "user_phone"
}

func New() *Phone {
	return new(Phone)
}

func (p Phone) Collection() *mongo.CollectionInfo {
	return mongo.Database().SetTable(p.TableName())
}

func (p Phone) collection() *mongo.CollectionInfo {
	return mongo.Database().SetTable(p.TableName())
}

func (p Phone) Init() error {
	ctx := context.Background()
	isUnique := true
	opts := options.IndexOptions{Unique: &isUnique}
	indexModel := mongoDriver.IndexModel{
		Keys:    []bson.M{{"Phone": 1}},
		Options: &opts,
	}
	_, err := p.collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (p Phone) GetByPhone(phonePrefix, phone string) (data Phone, err error) {
	where := bson.M{"phone_prefix": phonePrefix, "phone": phone}
	err = p.collection().Where(where).FindOne(&data)
	return data, err
}

func (p Phone) Add(data Phone) (Phone, error) {
	if _, err := p.collection().InsertOne(&data); err != nil {
		return Phone{}, err
	}
	return data, nil
}

func (p Phone) Save(data Phone) (Phone, error) {
	item, err := p.GetByPhone(data.PhonePrefix, data.Phone)
	if err != nil {
		if dao.IsNoDocumentErr(err) {
			return p.Add(data)
		}
		return Phone{}, err
	}
	return item, err
}

func (p Phone) SaveMany(data []Phone) ([]interface{}, error) {
	item, err := p.collection().InsertMany(&data)
	return item.InsertedIDs, err
}

func (p Phone) DelById(id string) error {
	where := bson.M{"_id": id}
	_, err := p.collection().Where(where).Delete()
	return err
}

func (p Phone) GetInfoByPhones(phones []string, fields string) ([]Phone, error) {
	var data []Phone
	where := bson.M{"phone": bson.M{"$in": phones}}
	err := p.collection().Where(where).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, err
}

func (p Phone) UpByMap(uid string, uData Phone) error {
	res, err := p.collection().UpByID(uid, uData)
	if err != nil {
		return err
	} else if res.ModifiedCount == 0 {
		return errors.New("update failed")
	}
	return nil
}
