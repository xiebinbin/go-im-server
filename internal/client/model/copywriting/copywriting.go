package copywriting

import (
	"imsdk/internal/common/dao/conversation/noticesetting"
	"imsdk/internal/common/dao/setting/copywriting"
	"imsdk/pkg/funcs"
)

type SaveParams = copywriting.CopyWriting

func List() []copywriting.CopyWriting {
	return copywriting.New().GetList()
}

func Save(data []SaveParams) error {
	return copywriting.New().AddMany(data)
}

func DelByAction(act string) error {
	return copywriting.New().DelByAction(act)
}

func GetSetting() []map[string]interface{} {
	resTmp := make(map[string]map[string]interface{}, 0)
	data, _ := noticesetting.New().GetSettings()
	for _, v := range data {
		id := funcs.Md516(v.Action + v.Lang)
		content := make(map[string]interface{}, 0)
		if resTmp[id] == nil {
			content[v.Role] = v.Content
			resTmp[id] = map[string]interface{}{
				"id":      v.ID,
				"action":  v.Action,
				"lang":    v.Lang,
				"type":    v.Type,
				"content": content,
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
