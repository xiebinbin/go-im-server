package dao

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mog "imsdk/pkg/database/mongo"
	"strings"
)

func GetMongoFieldsBsonByString(fields string) bson.M {
	if fields == "all" || fields == "" {
		return nil
	}
	fieldsSlice := strings.Split(fields, ",")
	var res = make(bson.M)
	for _, f := range fieldsSlice {
		f = strings.Replace(f, " ", "", -1)
		res[f] = 1
	}
	return res
}

func DataIsSaveSuccessfully(err error) bool {
	if err == nil {
		return true
	}
	return mongo.IsDuplicateKeyError(err)
}

func IsNoDocumentErr(err error) bool {
	return err == mongo.ErrNoDocuments
}

type MongoDb interface {
	DbName() string
	TableName() string
	Collection(mode ...mog.Mode) *mog.CollectionInfo
}

type MongoDbBase struct {
	MongoDb
}

func (b MongoDbBase) DbName() string {
	return "im"
}

func (b MongoDbBase) TableName() string {
	return ""
}

func (b *MongoDbBase) Collection(mode ...mog.Mode) *mog.CollectionInfo {
	rpMode := mog.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	fmt.Println(b.DbName())
	return mog.Database(b.DbName()).SetTable(b.TableName(), rpMode)
}
