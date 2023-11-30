package apk

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"strconv"
)

type Apk struct {
	ID        string                 `bson:"_id" json:"id"`
	Desc      string                 `bson:"desc" json:"desc,omitempty"`
	Version   string                 `bson:"version" json:"version"`
	Ver       int64                  `bson:"ver" json:"ver,omitempty"`
	Apk       map[string]interface{} `bson:"apk" json:"apk"`
	Type      int8                   `bson:"type" json:"type,omitempty"`
	Status    int8                   `bson:"status" json:"status,omitempty"`
	CreatedAt int64                  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt int64                  `bson:"update_time" json:"update_time,omitempty"`
}

const ()

func New() *Apk {
	return new(Apk)
}

func (a Apk) TableName() string {
	return "app_version_apk"
}

func GetId(version string, osType int8) string {
	return funcs.Md516(version + strconv.Itoa(int(osType)))
}

func (a Apk) Init() error {
	ctx := context.Background()
	indexModel := mongoDriver.IndexModel{
		Keys: bson.M{"ver": -1},
	}
	_, err := a.collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (a Apk) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database("tmm").SetTable(a.TableName(), rpMode)
}

func (a Apk) Add(data Apk) error {
	return a.collection().Where(bson.M{"_id": data.ID}).Upsert(&data).Err()
}

func (a Apk) GetDetail(id, fields string) (Apk, error) {
	var res Apk
	err := a.collection().Where(bson.M{"_id": id}).Fields(dao.GetMongoFieldsBsonByString(fields)).FindOne(&res)
	return res, err
}

func (a Apk) Update(data Apk) error {
	where := bson.M{"_id": data.ID}
	_, err := a.collection().Where(where).UpdateOne(&data)
	return err
}

func (a Apk) Lists(page, row int64, where, sort bson.M, fields string) ([]Apk, error) {
	data := make([]Apk, 0)
	skip := (page - 1) * row
	err := a.collection().Where(where).Sort(sort).Limit(row).Skip(skip).
		Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, err
}

func (a Apk) Count(where bson.M) int64 {
	count := a.collection().Where(where).Count()
	return count
}

func (a Apk) GetLatestByType(osType int) (Apk, error) {
	var res Apk
	err := a.collection().Where(bson.M{"type": osType}).Sort(bson.M{"ver": -1}).FindOne(&res)
	return res, err
}

func (a Apk) GetLatestByTypes(osTypes []int) ([]Apk, error) {
	var res []Apk
	where := bson.M{"type": bson.M{"$in": osTypes}}
	err := a.collection().Where(where).Sort(bson.M{"ver": -1}).Limit(int64(len(osTypes))).FindMany(&res)
	return res, err
}
