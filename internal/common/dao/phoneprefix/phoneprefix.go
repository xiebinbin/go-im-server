package phoneprefix

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
)

type PhonePrefix struct {
	ID          string `bson:"_id" json:"id"`
	Lang        string `bson:"lang" json:"lang"`
	Name        string `bson:"name" json:"name"`
	Code        string `bson:"code" json:"code"`
	CountryCode string `bson:"c_code" json:"code"`
}

func New() *PhonePrefix {
	return new(PhonePrefix)
}

func (p PhonePrefix) TableName() string {
	return "country_phone_prefix"
}

func (p PhonePrefix) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database("tmm").SetTable(p.TableName(), rpMode)
}

func (p PhonePrefix) List(lang string) []PhonePrefix {
	var data []PhonePrefix
	where := bson.M{"lang": lang}
	p.collection(mongo.SecondaryPreferredMode).Where(where).FindOne(&data)
	return data
}
