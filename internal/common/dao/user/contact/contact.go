package contact

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
)

type Contact struct {
	ID         string `bson:"_id" json:"id"`
	UId        string `bson:"uid" json:"uid"`
	ObjUId     string `bson:"obj_uid" json:"obj_uid"`
	Alias      string `bson:"alias" json:"alias"`
	RemarkText string `bson:"remark_text" json:"remark_text"`
	RemarkImg  string `bson:"remark_img" json:"remark_img"`
	Tag        string `bson:"tag" json:"tag"`
	Phone      string `bson:"phone" json:"phone"`
	HideIt     int8   `bson:"hide_it" json:"hide_it"`
	HideMe     int8   `bson:"hide_me" json:"hide_me"`
	IsTop      int8   `bson:"is_top" json:"is_top"`
	CreatedAt  int64  `bson:"itime" json:"create_time"`
	UpdatedAt  int64  `bson:"utime" json:"update_time"`

	Phones    []map[string]interface{} `bson:"phones" json:"phones"`
	Relations []map[string]interface{} `bson:"relations" json:"relations"`
	Emails    []map[string]interface{} `bson:"emails" json:"emails"`
	Dates     []map[string]interface{} `bson:"dates" json:"dates"`
	Companies []map[string]interface{} `bson:"companies" json:"companies"`
	Schools   []map[string]interface{} `bson:"schools" json:"schools"`
	Address   []map[string]interface{} `bson:"addrs" json:"addrs"`
}
type Phone struct {
	ID  string `json:"id"`
	Val string `json:"val"`
}

type Relation struct {
	ID  string `json:"id"`
	UId string `json:"uid"`
	Val string `json:"val"`
}

func New() *Contact {
	return new(Contact)
}

func (c Contact) TableName() string {
	return "user_contact"
}
func GetId(uid, objId string) string {
	return funcs.Md516(uid + objId)
}

func (c *Contact) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(c.TableName(), rpMode)
}

func (c Contact) GetContactsInfo(uid string, targetIds []string) (data []Contact, err error) {
	ids := make([]string, 0)
	for _, id := range targetIds {
		ids = append(ids, GetId(uid, id))
	}
	where := bson.M{"_id": bson.M{"$in": ids}}
	c.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	return
}

func (c Contact) UpdateRemark(uid, objUId string, upData map[string]interface{}) (int64, error) {
	id := GetId(uid, objUId)
	where := bson.M{"_id": id}
	res, err := c.Collection().Where(where).UpdateOne(&upData)
	return res.MatchedCount, err
}

func (c Contact) UpdateById(id string, upData map[string]interface{}) (int64, error) {
	where := bson.M{"_id": id}
	res, err := c.Collection().Where(where).UpdateOne(&upData)
	return res.MatchedCount, err
}

func (c Contact) AddRemark(addData map[string]interface{}) (string, error) {
	res, err := c.Collection().InsertOne(addData)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(string), nil
}

func (c Contact) GetAlias(uid, objUId string) string {
	var data Contact
	id := GetId(uid, objUId)
	where := bson.M{"_id": id}
	c.Collection(mongo.SecondaryPreferredMode).Where(where).FindOne(&data)
	return data.Alias
}

func (c Contact) GetByIds(ids []string) []Contact {
	var data []Contact
	where := bson.M{"_id": bson.M{"$in": ids}}
	c.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	return data
}

func (c Contact) GetAllLists() []Contact {
	var data []Contact
	where := bson.M{}
	c.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	return data
}
