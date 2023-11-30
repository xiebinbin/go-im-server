package thirdprovider

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/sdk"
)

type StatusType = sdk.StatusType

type ThirdProvider struct {
	ID        string     `bson:"_id" json:"id"` // ak
	AK        string     `bson:"ak" json:"ak"`  // appId
	SK        string     `bson:"sk" json:"sk"`  // secretKey
	PK        string     `bson:"pk" json:"pk"`  // publicKey
	Status    StatusType `bson:"status" json:"status"`
	CreatedAt int64      `bson:"create_time" json:"create_time"`
	UpdatedAt int64      `bson:"update_time" json:"update_time"`
}

func (s ThirdProvider) TableName() string {
	return "third_provider"
}

func New() *ThirdProvider {
	return new(ThirdProvider)
}

func (s *ThirdProvider) getCollection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database("tmm_sdk_server").SetTable(s.TableName(), rpMode)
}

func (s *ThirdProvider) InsertOne(data ThirdProvider) error {
	_, err := s.getCollection().InsertOne(data)
	return err
}

func (s *ThirdProvider) UpsertOne(data ThirdProvider) error {
	err := s.getCollection().Where(bson.M{"_id": data.ID}).Upsert(&data)
	return err.Err()
}

func (s *ThirdProvider) GetByAKey(ak string) (ThirdProvider, error) {
	var data ThirdProvider
	err := s.getCollection().Where(bson.M{"ak": ak}).FindOne(&data)
	return data, err
}
