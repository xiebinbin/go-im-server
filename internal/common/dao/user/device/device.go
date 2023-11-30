package device

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
)

type Device struct {
	ID        string `bson:"_id" json:"id"` //deviceId+os+uid
	UID       string `bson:"uid" json:"uid"`
	DeviceId  string `bson:"device_id" json:"device_id"`
	Os        string `bson:"os" json:"os,omitempty"`
	CreateAt  int64  `bson:"create_time" json:"create_time"`
	UpdatedAt int64  `bson:"update_time" json:"update_time"`
}

func New() *Device {
	return new(Device)
}

func (d *Device) TableName() string {
	return "user_device"
}

func (d *Device) collection() *mongo.CollectionInfo {
	return mongo.Database().SetTable(d.TableName())
}

func GetId(uid, deviceId, os string) string {
	return funcs.Md516(uid + deviceId + os)
}

func (d *Device) InsertOne(addData Device) (string, error) {
	res, err := d.collection().InsertOne(&addData)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(string), nil
}

func (d *Device) UpsertOne(data Device) error {
	return d.collection().Where(bson.M{"_id": data.ID}).Upsert(&data).Err()
}

func (d *Device) GetDeviceIdsByUID(uid string) ([]Device, error) {
	data := make([]Device, 0)
	err := d.collection().Where(bson.M{"uid": uid}).FindMany(&data)
	return data, err
}
