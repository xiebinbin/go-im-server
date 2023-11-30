package maxsequence

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
)

type MaxSequence struct {
	ID  string `bson:"_id" json:"_id"` // message id
	Seq int64  `bson:"seq" json:"seq"`
}

func New() *MaxSequence {
	return new(MaxSequence)
}

func (m MaxSequence) TableName() string {
	return "message_max_sequence"
}

func (m MaxSequence) getCollection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(m.TableName(), rpMode)
}

func (m MaxSequence) UpdateSeq(uid string, seq int64) {
	data := bson.M{"$set": bson.M{"_id": uid, "seq": seq}}
	m.getCollection().Where(bson.M{"_id": uid}).UpsertByBson(data)
}

func (m MaxSequence) GetSeq(uid string) (data MaxSequence) {
	where := bson.M{"_id": uid}
	m.getCollection().Where(where).FindOne(&data)
	return
}
