package language

import (
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
)

type Language struct {
	ID   string `bson:"_id" json:"id,omitempty"`
	Name string `bson:"name" json:"name,omitempty"`
	Rank int8   `bson:"rank" json:"rank,omitempty"`
}

func New() *Language {
	return new(Language)
}

func (l *Language) TableName() string {
	return "language"
}

func (l *Language) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(l.TableName(), rpMode)
}

func (l *Language) GetLanguageLists(fields string) ([]Language, error) {
	var data []Language
	err := l.Collection(mongo.PrimaryMode).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
