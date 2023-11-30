package copywriting

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
)

type CopyWriting struct {
	ID     string `bson:"_id" json:"id"`
	Lang   string `bson:"lang" json:"lang"`
	Action string `bson:"action" json:"action"`
	Text   string `bson:"text" json:"text"`
	Type   string `bson:"type" json:"type"`
}

const (
	TypeCommon  = "c"
	TypeMoments = "m"
	TypePush    = "p"
)

func New() *CopyWriting {
	return new(CopyWriting)
}

func (c CopyWriting) TableName() string {
	return "copywriting_v2"
}

func (c CopyWriting) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(c.TableName(), rpMode)
}

func (c CopyWriting) DelByAction(act string) error {
	_, err := c.collection().Where(bson.M{"action": act}).Delete()
	return err
}

func (c CopyWriting) GetList() []CopyWriting {
	var data []CopyWriting
	c.collection(mongo.SecondaryPreferredMode).FindMany(&data)
	return data
}

func (c CopyWriting) GetId(action, bizType, lang string) string {
	return funcs.Md516(action + bizType + lang)
}

func (c CopyWriting) AddMany(data []CopyWriting) error {
	if data == nil {
		return errors.New("nil data")
	}
	for k, v := range data {
		data[k].ID = c.GetId(v.Action, v.Type, v.Lang)
	}
	_, err := c.collection().InsertMany(&data)
	return err
}

func (c CopyWriting) GetMsgListByLang(action, lang string) []CopyWriting {
	res := make([]CopyWriting, 0)
	where := bson.M{"lang": lang, "action": action}
	c.collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&res)
	return res
}

func (c CopyWriting) GetText(action, bizType, lang string) string {
	id := c.GetId(action, bizType, lang)
	var data CopyWriting
	c.collection(mongo.SecondaryPreferredMode).Fields(dao.GetMongoFieldsBsonByString("text")).FindByID(id, &data)
	return data.Text
}
