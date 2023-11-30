package setting

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/internal/common/dao"
	"imsdk/pkg/app"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"reflect"
)

type Setting struct {
	ID              string `bson:"_id" json:"id"`
	FindById        uint8  `bson:"find_f_tid" json:"find_f_tid"`       // search by tmm id
	FindByPhone     uint8  `bson:"find_f_phone" json:"find_f_phone"`   // search by phone
	AddFromGroup    uint8  `bson:"add_f_group" json:"add_f_group"`     // added friend by group
	AddFromQrCode   uint8  `bson:"add_f_qrcode" json:"add_f_qrcode"`   // added friend by qr code
	AddFromCard     uint8  `bson:"add_f_card" json:"add_f_card"`       // added by card
	AddFromMoments  uint8  `bson:"add_f_moments" json:"add_f_moments"` //
	AddFromDiscover uint8  `bson:"add_f_disc" json:"add_f_disc"`
	AutoAddFriend   uint8  `bson:"auto_add_friend" json:"auto_add_friend"`
	Language        string `bson:"language" json:"language"`
	SubLang         string `bson:"sub_lang" json:"sub_lang"` // if language equal 1 (auto) ,sub_lang is system language
	CreatedAt       int64  `bson:"create_time" json:"create_time"`
	UpdatedAt       int64  `bson:"update_time" json:"update_time"`
}

const (
	AutoLang        = "1"
	StatusYes uint8 = 1
	StatusNo  uint8 = 2
)

var (
	DefaultSetting = Setting{
		FindById:        StatusYes,
		FindByPhone:     StatusYes,
		AddFromGroup:    StatusYes,
		AddFromQrCode:   StatusYes,
		AddFromCard:     StatusYes,
		AddFromMoments:  StatusYes,
		AddFromDiscover: StatusYes,
		AutoAddFriend:   StatusNo,
		Language:        AutoLang,
		SubLang:         "",
	}
	defSettingMap = map[string]interface{}{
		"find_f_tid":      StatusYes,
		"find_f_phone":    StatusYes,
		"add_f_group":     StatusYes,
		"add_f_qrcode":    StatusYes,
		"add_f_card":      StatusYes,
		"add_f_moments":   StatusYes,
		"add_f_disc":      StatusYes,
		"auto_add_friend": StatusNo,
		"language":        AutoLang,
		"sub_lang":        "",
	}
)

func (s *Setting) TableName() string {
	return "user_config"
}

func New() *Setting {
	return new(Setting)
}

func (s *Setting) getCollection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(s.TableName(), rpMode)
}

func (s *Setting) GetByID(id, fields string) (Setting, error) {
	var data Setting
	err := s.getCollection(mongo.SecondaryPreferredMode).FindByID(id, &data)
	return data, err
}

func (s *Setting) Add(data Setting) bool {
	_, err := s.getCollection().InsertOne(data)
	return err == nil
}

func (s *Setting) GetFieldName(field string) string {
	ins := Setting{}
	reflectIns := reflect.TypeOf(ins)
	if fieldObj, ok := reflectIns.FieldByName(field); ok {
		return fieldObj.Tag.Get("bson")
	}
	return ""
}

func (s *Setting) UpdateValue(id, field string, value interface{}) error {
	f := s.GetFieldName(field)
	if f == "" {
		return errors.New("field not exist")
	}
	uData := map[string]interface{}{
		f: value,
	}
	res, err := s.getCollection().UpByID(id, uData)
	if err == nil && res.ModifiedCount == 1 {
		return nil
	} else if err == nil && res.MatchedCount == 0 {
		data := defSettingMap
		data["_id"] = id
		data[f] = value
		_, err = s.getCollection().InsertOne(data)
		fmt.Println("init:config", data)
	}
	return err
}

func (s *Setting) GetByIDs(ids []string, fields string) []Setting {
	data := make([]Setting, 0)
	f := dao.GetMongoFieldsBsonByString(fields)
	err := s.getCollection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": ids}}).Fields(f).FindMany(&data)
	fmt.Println("setting--GetByIDs", ids, f, bson.M{"_id": bson.M{"$in": ids}}, data, err)
	return data
}

func (s *Setting) GetUsersLang(ids []string) map[string]string {
	var langConf struct {
		Languages []string `toml:"languages"`
		Default   string   `toml:"default"`
	}
	app.Config().Bind("global", "languages", &langConf)

	var res = make(map[string]string)
	data := s.GetByIDs(ids, "_id,language,sub_lang")
	if data == nil {
		for _, id := range ids {
			res[id] = langConf.Default
		}
		return res
	}
	for _, v := range data {
		lang := v.Language
		if lang == AutoLang {
			lang = v.SubLang
		}
		if !funcs.In(lang, langConf.Languages) {
			lang = langConf.Default
		}
		res[v.ID] = lang
	}
	return res
}

func (s *Setting) GetUserLang(id string) string {
	data, _ := s.GetByID(id, "_id,language,sub_lang")
	var langConf struct {
		Languages []string `toml:"languages"`
		Default   string   `toml:"default"`
	}
	app.Config().Bind("global", "languages", &langConf)
	if data.ID == "" {
		return langConf.Default
	}
	lang := data.Language
	if lang == AutoLang {
		lang = data.SubLang
	}
	if !funcs.In(lang, langConf.Languages) {
		lang = langConf.Default
	}
	return lang
}
