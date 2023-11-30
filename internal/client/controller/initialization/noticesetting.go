package initialization

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/conversation/notice"
	"imsdk/internal/common/dao/conversation/noticesetting"
	noticesettingv2 "imsdk/internal/common/dao/conversation/noticesetting/v2"
	"imsdk/pkg/app"
	"imsdk/pkg/funcs"
	"imsdk/pkg/translate/google"
	"os"
)

func InitNoticeSetting(ctx *gin.Context) {
	filename := app.Config().GetPublicConfigDir() + "notice_setting.json"
	filePtr, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
		return
	}
	defer filePtr.Close()
	var settings []noticesetting.Setting
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&settings)

	//languages := []string{"ar", "cs", "tl", "de", "da", "el", "es", "sw", "fi", "fr", "hu", "id",
	//	"is", "it", "ja", "ko", "lt", "ms", "no", "pl", "pt", "ro", "ru", "sk", "sl", "sv", "th", "uk", "vi"}
	//languages := []string{"ar", "zh-TW", "cs", "id"}
	languages := []string{"id", "zh-TW"}
	//fmt.Println(languages)
	mills := funcs.GetMillis()
	for _, v := range settings {
		addData := make([]noticesetting.Setting, 0)
		ids := make([]string, 0)
		id := funcs.Md516(v.Action + v.Lang + v.Role)
		tmp := noticesetting.Setting{
			ID:        id,
			Lang:      v.Lang,
			Action:    v.Action,
			Role:      v.Role,
			Content:   v.Content,
			Type:      20,
			CreatedAt: mills,
			UpdatedAt: mills,
		}
		addData = append(addData, tmp)
		ids = append(ids, id)
		err = noticesetting.New().DelMany(ids)
		if err == nil {
			err = noticesetting.New().AddMany(addData)
			fmt.Println("err=====", err)
		}
		if v.Lang == "en" {
			enStr := replaceNoticeStr(v.Content)
			fmt.Println("enStr======", enStr)
			for _, v1 := range languages {
				text, _ := google.TranslV2(google.TransData{Text: enStr, TargetLang: v1})
				returnStr := returnNoticeStr(text)
				fmt.Println("returnStr======", v1, returnStr)
				lang := v1
				if v1 == "zh-TW" {
					lang = "zh-hant"
				} else if v1 == "tl" {
					lang = "fil"
				}
				addDataV2 := make([]noticesetting.Setting, 0)
				ids2 := make([]string, 0)
				id2 := funcs.Md516(v.Action + lang + v.Role)
				tmp2 := noticesetting.Setting{
					ID:        id2,
					Lang:      lang,
					Action:    v.Action,
					Role:      v.Role,
					Content:   returnStr,
					Type:      20,
					CreatedAt: mills,
					UpdatedAt: mills,
				}
				addDataV2 = append(addDataV2, tmp2)
				ids2 = append(ids2, id2)
				if noticesetting.New().DelMany(ids2) == nil {
					err = noticesetting.New().AddMany(addDataV2)
					fmt.Println("err=====", err)
				}
			}
		}
	}
	fmt.Println("saveData=====success")
}

func InitNoticeSettingV2(ctx *gin.Context) {
	list := notice.GetSetting(ctx)
	err := notice.AddSettingV2(ctx, list)
	if err != nil {
		return
	}
	return
	filename := app.Config().GetPublicConfigDir() + "notice_setting.json"
	filePtr, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
		return
	}
	defer filePtr.Close()
	var settings []noticesettingv2.Setting
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&settings)

	//languages := []string{"ar", "cs", "tl", "de", "da", "el", "es", "sw", "fi", "fr", "hu", "id",
	//	"is", "it", "ja", "ko", "lt", "ms", "no", "pl", "pt", "ro", "ru", "sk", "sl", "sv", "th", "uk", "vi"}
	//languages := []string{"ar", "zh-TW", "cs", "id"}
	languages := []string{"id", "zh-TW"}
	mills := funcs.GetMillis()
	dao := noticesettingv2.New()
	seq, _ := dao.GetMaxSequence()
	fmt.Println("seq:", seq)
	for _, v := range settings {
		id := funcs.Md516(v.Action + v.Lang)
		seq += 1
		tmp := noticesettingv2.Setting{
			ID:        id,
			Lang:      v.Lang,
			Action:    v.Action,
			Content:   v.Content,
			Type:      20,
			Sequence:  seq,
			CreatedAt: mills,
			UpdatedAt: mills,
		}
		fmt.Println("tmp:", tmp)
		err = dao.Upsert(tmp)
		if err != nil {
			return
		}
		if v.Lang == "en" {
			var roleInfo noticesettingv2.Content
			c, _ := json.Marshal(v.Content)
			json.Unmarshal(c, &roleInfo)
			fmt.Println("roleInfo:", roleInfo)
			for _, v1 := range languages {
				content := make(map[string]interface{}, 0)
				seq += 1
				if roleInfo.Operator != "" {
					enStr := replaceNoticeStr(roleInfo.Operator)
					text, _ := google.TranslV2(google.TransData{Text: enStr, TargetLang: v1})
					content["operator"] = returnNoticeStr(text)
				}

				if roleInfo.Target != "" {
					enStr := replaceNoticeStr(roleInfo.Target)
					text, _ := google.TranslV2(google.TransData{Text: enStr, TargetLang: v1})
					content["target"] = returnNoticeStr(text)
				}

				if roleInfo.Other != "" {
					enStr := replaceNoticeStr(roleInfo.Other)
					text, _ := google.TranslV2(google.TransData{Text: enStr, TargetLang: v1})
					content["other"] = returnNoticeStr(text)
				}

				lang := v1
				if v1 == "zh-TW" {
					lang = "zh-hant"
				} else if v1 == "tl" {
					lang = "fil"
				}
				id2 := funcs.Md516(v.Action + lang)
				tmp2 := noticesettingv2.Setting{
					ID:        id2,
					Lang:      lang,
					Action:    v.Action,
					Content:   content,
					Type:      20,
					Sequence:  seq,
					CreatedAt: mills,
					UpdatedAt: mills,
				}
				fmt.Println("tmp2:", tmp2)
				dao.Upsert(tmp2)
			}
		}
	}
	fmt.Println("saveData=====success")
}
