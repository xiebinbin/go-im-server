package usersinglesetting

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
)

type Setting struct {
	ID           string                 `bson:"_id" json:"id"`
	ChatID       string                 `bson:"chat_id" json:"chat_id"` // conversation id
	UID          string                 `bson:"uid" json:"uid"`
	IsTop        uint8                  `bson:"is_top" json:"is_top"`
	HideSequence int64                  `bson:"hide_sequence" json:"hide_sequence"`
	IsMuteNotify uint8                  `bson:"is_mute_notify" json:"is_mute_notify"`
	MuteTime     int64                  `bson:"mute_time" json:"mute_time"`
	TopTime      int64                  `bson:"top_time" json:"top_time"`
	Background   map[string]interface{} `bson:"background" json:"-"`
	CreatedAt    int64                  `bson:"create_time" json:"create_time"`
	UpdatedAt    int64                  `bson:"update_time" json:"update_time"`
}

var DefSetting = Setting{
	IsTop:        0,
	HideSequence: -1,
	IsMuteNotify: 0,
	MuteTime:     0,
	TopTime:      0,
	Background:   nil,
}

type SetDataResponse struct {
	UID          string `bson:"uid" json:"uid"`
	IsMuteNotify uint8  `bson:"is_mute_notify" json:"is_mute_notify"`
}

func (s Setting) TableName() string {
	return "conversation_user_setting_single"
}

func New() *Setting {
	return new(Setting)
}

func (s *Setting) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(s.TableName(), rpMode)
}

func (s *Setting) Init() error {
	ctx := context.Background()
	indexModel := mongoDriver.IndexModel{
		Keys: bson.M{"chat_id": 1},
	}
	_, err := s.Collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func GetId(uid, chatId string) string {
	return funcs.Md516(uid + chatId)
}

func (s *Setting) GetSetting(uid, chatId string) (Setting, error) {
	id := GetId(uid, chatId)
	var data Setting
	err := s.Collection(mongo.SecondaryPreferredMode).FindByID(id, &data)
	return data, err
}

func (s *Setting) GetMuteSettingsForChatId(chatId string) map[string]uint8 {
	ResData := make([]SetDataResponse, 0)
	where := bson.M{"chat_id": chatId}
	s.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&ResData)
	data := make(map[string]uint8, 0)
	if ResData != nil {
		for _, v := range ResData {
			data[v.UID] = v.IsMuteNotify
		}
	}
	return data
}

func (s Setting) Add(data Setting) (string, error) {
	res, err := s.Collection().InsertOne(data)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(string), nil
}

func (s Setting) Update(id string, uData Setting) bool {
	res, err := s.Collection().UpByID(id, uData)
	if err == nil && res.ModifiedCount == 1 {
		return true
	}
	return false
}

func (s Setting) UpdateValue(id string, field string, value interface{}) bool {
	uData := map[string]interface{}{
		field: value,
	}
	res, err := s.Collection().UpByID(id, uData)
	if err == nil && res.ModifiedCount == 1 {
		return true
	}
	return false
}

func (s Setting) UpdateByMap(id string, uData map[string]interface{}) error {
	_, err := s.Collection().UpByID(id, uData)
	return err
}

func (s Setting) GetSettings(uid string, chatId []string) []Setting {
	data := make([]Setting, 0)
	ids := make([]string, 0)
	for _, v := range chatId {
		ids = append(ids, GetId(uid, v))
	}
	where := bson.M{"_id": bson.M{"$in": ids}}
	s.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	return data
}

func (s Setting) AddMany(data []Setting) interface{} {
	res, _ := s.Collection().InsertMany(data)
	return res.InsertedIDs
}

func (s *Setting) GetCount() int64 {
	count := s.Collection().Count()
	return count
}

func (s *Setting) GetListByLimit(limit, offset int64) []Setting {
	var data []Setting
	fmt.Println("offset:", offset)
	s.Collection().Limit(limit).Skip(offset).FindMany(&data)
	//Select("_id,avatar").Limit(limit).Offset(offset).Find(&data)
	return data
}
