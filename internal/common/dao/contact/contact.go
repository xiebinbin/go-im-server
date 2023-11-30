package contact

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
)

type Contact struct {
	ID        string `bson:"_id" json:"id"`
	UID       string `bson:"uid" json:"uid"`
	Prefix    string `bson:"prefix" json:"prefix"`
	Phone     string `bson:"phone" json:"phone"`
	PhoneOri  string `bson:"phone_ori" json:"phone_ori"`
	Alias     string `bson:"alias" json:"alias"`
	CreatedAt int64  `bson:"create_time" json:"create_time"`
}

func New() *Contact {
	return new(Contact)
}

func GetId(uid, prefix, phone string) string {
	return funcs.SHA1Base64(uid + prefix + phone)
}

func (c Contact) TableName() string {
	return "user_contact_report"
}

func (c *Contact) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(c.TableName(), rpMode)
}

func (c Contact) AddMany(addData []Contact) ([]interface{}, error) {
	res, err := c.Collection().InsertMany(addData)
	if err != nil {
		return []interface{}{}, err
	}
	return res.InsertedIDs, nil
}

func (c Contact) Add(addData Contact) (string, error) {
	res, err := c.Collection().InsertOne(addData)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(string), nil
}

func (c Contact) GetById(id string) []Contact {
	var data []Contact
	where := bson.M{"_id": id}
	c.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	return data
}

func (c Contact) GetByIds(Ids []string) ([]Contact, error) {
	var res []Contact
	where := bson.M{"_id": bson.M{"$in": Ids}}
	err := c.Collection().Where(where).FindMany(&res)
	return res, err
}

func (c Contact) DeleteByUId(uid string) error {
	where := bson.M{"uid": uid}
	_, err := c.Collection().Where(where).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (c Contact) GetIsMineInfo(uid, prefix, phone string) Contact {
	var data Contact
	where := bson.M{"uid": uid, "prefix": prefix, "phone": phone}
	err := c.Collection(mongo.SecondaryPreferredMode).Where(where).FindOne(&data)
	if err != nil {

	}
	return data
}
