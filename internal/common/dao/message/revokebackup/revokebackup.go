package revokebackup

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
)

type RevokeBackup struct {
	ID       string `bson:"_id" json:"id"`
	SenderId string `bson:"sender_id" json:"sender_id"`
	Content  string `bson:"content" json:"content"`
	CreateAt int64  `bson:"create_time" json:"create_time"`
}

func (r *RevokeBackup) TableName() string {
	return "message_revoke_backup"
}

func New() *RevokeBackup {
	return new(RevokeBackup)
}

func (r *RevokeBackup) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(r.TableName(), rpMode)
}

func (r *RevokeBackup) DeleteByIds(ids []string) (int64, error) {
	where := bson.M{"_id": bson.M{"$in": ids}}
	i, err := r.collection().Where(where).Delete()
	if err != nil {
		return 0, err
	}
	return i, err
}

func (r *RevokeBackup) SaveMany(params []RevokeBackup) ([]interface{}, error) {
	i, err := r.collection().InsertMany(&params)
	if err != nil {
		return nil, err
	}
	return i.InsertedIDs, nil
}
