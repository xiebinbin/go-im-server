package region

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
)

type Region struct {
	ID    string   `bson:"_id" json:"id,omitempty"`
	Lang  string   `bson:"lang" json:"lang,omitempty"`
	Name  string   `bson:"name" json:"name,omitempty"`
	Code  string   `bson:"code" json:"code,omitempty"`
	PId   string   `bson:"pid" json:"pid,omitempty"`
	PIds  []string `bson:"pids" json:"pids,omitempty"`
	Level int8     `bson:"level" json:"level,omitempty"`
}

func New() *Region {
	return new(Region)
}

func (r *Region) TableName() string {
	return "region"
}

func (r *Region) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(r.TableName(), rpMode)
}

func (r *Region) GetNameById(id string) Region {
	var data Region
	where := bson.M{"_id": id}
	mongo.Database().SetTable(r.TableName(), mongo.SecondaryPreferredMode).Where(where).FindOne(&data)
	return data
}

func (r *Region) GetCountryList(lang string) ([]Region, error) {
	var data []Region
	where := bson.M{"lang": lang, "level": bson.M{"$in": []int8{1, 2}}}
	err := r.Collection(mongo.PrimaryMode).Where(where).FindMany(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *Region) GetListByLang(lang string, fields string) ([]Region, error) {
	var data []Region
	where := bson.M{"lang": lang}
	err := r.Collection(mongo.PrimaryMode).Fields(dao.GetMongoFieldsBsonByString(fields)).
		Where(where).FindMany(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
