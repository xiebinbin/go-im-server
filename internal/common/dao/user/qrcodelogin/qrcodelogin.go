package qrcodelogin

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
)

type QrCodeLogin struct {
	ID          string `bson:"_id" json:"id"`
	Uid         string `bson:"uid" json:"uid"` // creator
	Sk          string `bson:"sk" json:"sk"`   // 加密后的密钥
	PhonePrefix string `bson:"phone_prefix" json:"phone_prefix"`
	Name        string `bson:"name" json:"name"` //
	Avatar      string `bson:"avatar" json:"avatar"`
	Expire      int64  `bson:"expire" json:"expire"`
	Status      int    `bson:"status" json:"status"`
	Ip          string `bson:"ip" json:"ip,omitempty"`
	Country     string `bson:"country" json:"country,omitempty"`
	Region      string `bson:"region" json:"region,omitempty"`
	City        string `bson:"city" json:"city,omitempty"`
	Timezone    string `bson:"timezone" json:"timezone,omitempty"`
	BandInfo    string `bson:"band_info" json:"band_info,omitempty"`
	Band        string `bson:"band" json:"band,omitempty"`
	BandModel   string `bson:"band_model" json:"band_model,omitempty"`
	DateIdx     string `bson:"date_idx" json:"date_idx,omitempty"`
	Date        int    `bson:"date" json:"date,omitempty"`
	CreateAt    int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt   int64  `bson:"update_time" json:"update_time,omitempty"`
	Version     string `bson:"version" json:"version,omitempty"`
	Os          string `bson:"os" json:"os,omitempty"`
	IsLock      int8   `bson:"is_lock" json:"is_lock"`
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
