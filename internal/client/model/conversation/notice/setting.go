package notice

import (
	"context"
	"fmt"
	"imsdk/internal/common/dao/conversation/noticesetting"
	noticesettingv2 "imsdk/internal/common/dao/conversation/noticesetting/v2"
	"imsdk/internal/common/dao/user/setting"
	"imsdk/pkg/funcs"
	"strings"
)

type GetSeqRequest struct {
	Sequence int64 `json:"sequence"`
}

func GetSetting(ctx context.Context) []map[string]interface{} {
	resTmp := make(map[string]map[string]interface{}, 0)
	data, _ := noticesetting.New().GetSettings()
	for _, v := range data {
		id := funcs.Md516(v.Action + v.Lang)
		content := make(map[string]interface{}, 0)
		if resTmp[id] == nil {
			content[v.Role] = v.Content
			resTmp[id] = map[string]interface{}{
				"id":          v.ID,
				"action":      v.Action,
				"lang":        v.Lang,
				"type":        v.Type,
				"content":     content,
				"create_time": v.CreatedAt,
				"update_time": v.UpdatedAt,
			}
		} else {
			tmp := resTmp[id]["content"].(map[string]interface{})
			tmp[v.Role] = v.Content
			resTmp[id]["content"] = tmp
		}
	}
	res := make([]map[string]interface{}, 0)
	for _, v := range resTmp {
		res = append(res, v)
	}
	return res
}

func GetSettingV2(ctx context.Context, params GetSeqRequest) []noticesettingv2.Setting {
	data, _ := noticesettingv2.New().GetSettingsBySeq(params.Sequence)
	return data
}

func AddSettingV2(ctx context.Context, data []map[string]interface{}) error {
	dao := noticesettingv2.New()
	seq, _ := dao.GetMaxSequence()
	for _, datum := range data {
		seq = seq + 1
		addData := noticesettingv2.Setting{
			ID:        datum["id"].(string),
			Sequence:  seq,
			Content:   datum["content"].(map[string]interface{}),
			Lang:      datum["lang"].(string),
			Action:    datum["action"].(string),
			Type:      datum["type"].(int8),
			CreatedAt: datum["create_time"].(int64),
			UpdatedAt: datum["update_time"].(int64),
		}
		er := dao.Upsert(addData)
		if er != nil {
			return er
		}
		fmt.Println("seq----:", seq, er)
	}
	fmt.Println("seq:", seq)
	return nil
}

func GetUserTempByRoleAndLang(uid, action, role string) noticesetting.Setting {
	lang := setting.New().GetUserLang(uid)
	noticeTmpInfo, _ := noticesetting.New().GetSettingByRoleAndLang(action, role, lang)
	return noticeTmpInfo
}

func ReplaceTmp(operationStr string, targetArr []string, tempTxt string) string {
	if tempTxt != "" {
		targetStr := strings.Join(targetArr, ",")
		tempTxt = strings.Replace(tempTxt, "@{operator}", operationStr, -1)
		tempTxt = strings.Replace(tempTxt, "@{target}", targetStr, -1)
		return tempTxt
	}
	return ""
}
