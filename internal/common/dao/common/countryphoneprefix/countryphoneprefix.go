package countryphoneprefix

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
)

type CountryPhonePrefix struct {
	ID    string `bson:"_id" json:"id"`
	Lang  string `bson:"lang" json:"lang"`
	Name  string `bson:"name" json:"name"`
	CCode string `bson:"c_code" json:"c_code"`
	Code  string `bson:"code" json:"code"`
	Sort  int16  `bson:"sort" json:"sort"`
}

func New() *CountryPhonePrefix {
	return new(CountryPhonePrefix)
}

func (c *CountryPhonePrefix) TableName() string {
	return "country_phone_prefix"
}

func (c *CountryPhonePrefix) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(c.TableName(), rpMode)
}

func (c *CountryPhonePrefix) GetList(lang string) ([]CountryPhonePrefix, error) {
	var data []CountryPhonePrefix
	where := bson.M{"lang": lang}
	err := c.Collection(mongo.PrimaryMode).Where(where).FindMany(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
