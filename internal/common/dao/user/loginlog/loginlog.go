package loginlog

import (
	"imsdk/pkg/database/mongo"
)

type LoginLog struct {
	ID         string `bson:"_id" json:"id"`
	UId        string `bson:"uid" json:"uid"`
	Os         string `bson:"os" json:"os"`
	DeviceName string `bson:"device_name" json:"device_name"`
	DeviceId   string `bson:"device_id" json:"device_id"`
	CreatedAt  int64  `bson:"create_time" json:"create_time"`
	UpdatedAt  int64  `bson:"update_time" json:"update_time"`
}

func New() *LoginLog {
	return new(LoginLog)
}

func (c LoginLog) TableName() string {
	return "login_log"
}

func (c LoginLog) AddLogin(addData LoginLog) (string, error) {
	res, err := mongo.Database().SetTable(c.TableName()).InsertOne(addData)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(string), nil
}
