package qrcodelogin

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
)

type QrCodeLogin struct {
	ID          string                 `bson:"_id" json:"id"`
	Uid         string                 `bson:"uid" json:"uid"` // creator
	PhonePrefix string                 `bson:"phone_prefix" json:"phone_prefix"`
	FirstName   string                 `bson:"f_name" json:"f_name"` //
	LastName    string                 `bson:"l_name" json:"l_name"` //
	Avatar      map[string]interface{} `bson:"avatar" json:"avatar"`
	Expire      int64                  `bson:"expire" json:"expire"`
	Status      int                    `bson:"status" json:"status"`
}

const (
	StatusWaitScan       = 0
	StatusScan           = 1
	StatusLoginConfirmed = 2
)

const (
	ErrorUsed = 100001
)

func New() *QrCodeLogin {
	return new(QrCodeLogin)
}

func (q QrCodeLogin) TableName() string {
	return "qrcode_pc_login"
}

func (q QrCodeLogin) collection() *mongo.CollectionInfo {
	return mongo.Database("tmm").SetTable(q.TableName())
}

func (q QrCodeLogin) Add(data QrCodeLogin) error {
	_, err := q.collection().InsertOne(&data)
	return err
}
func (q QrCodeLogin) Save(data QrCodeLogin) error {
	return q.collection().Where(bson.M{"_id": data.ID}).Upsert(&data).Err()
}

func (q QrCodeLogin) GetByID(id string) (QrCodeLogin, error) {
	var data QrCodeLogin
	err := q.collection().FindByID(id, &data)
	return data, err
}

func (q QrCodeLogin) UpByID(id string, uData QrCodeLogin) error {
	_, err := q.collection().UpByID(id, &uData)
	return err
}
