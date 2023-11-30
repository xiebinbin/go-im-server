package noticesetting

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
)

type Setting struct {
	ID        string `bson:"_id" json:"id"`
	Lang      string `bson:"lang" json:"lang"`
	Action    string `bson:"action" json:"action"`
	Role      string `bson:"role" json:"role"`
	Type      int8   `bson:"type" json:"type"`
	Content   string `bson:"content" json:"content"`
	CreatedAt int64  `bson:"create_time" json:"create_time"`
	UpdatedAt int64  `bson:"update_time" json:"update_time"`
}

func New() *Setting {
	return new(Setting)
}

func (s Setting) TableName() string {
	return "notice_setting"
}

func (s *Setting) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(s.TableName(), rpMode)
}

func (s Setting) GetSettingByRoleAndLang(action, role, lang string) (Setting, error) {
	var data Setting
	where := bson.M{"action": action, "role": role, "lang": lang}
	err := s.Collection(mongo.SecondaryPreferredMode).Where(where).FindOne(&data)
	if err != nil {
		return Setting{}, err
	}
	return data, nil
}

func (s Setting) AddSetting(data []Setting) error {
	s.Collection().Where(bson.M{"_id": bson.M{"$ne": ""}}).Delete()
	_, err := s.Collection().InsertMany(&data)
	return err
}

func (s Setting) AddMany(data []Setting) error {
	if data == nil {
		return errors.New("nil data")
	}
	_, err := s.Collection().InsertMany(&data)
	return err
}

func (s Setting) DelMany(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	where := bson.M{"_id": bson.M{"$in": ids}}
	_, err := s.Collection().Where(where).Delete()
	return err
}

func (s Setting) GetSettings() ([]Setting, error) {
	var data []Setting
	s.Collection().FindMany(&data)
	return data, nil
}
