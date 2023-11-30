package setting

import (
	"imsdk/internal/common/dao/setting/copywriting"
	"imsdk/pkg/funcs"
	"regexp"
	"strconv"
	"strings"
)

const DefaultLang = "en"

func ReplaceMsgByLang(msg, lang string) string {
	alreadyLang := GetAlreadyLang()
	if !funcs.In(lang, alreadyLang) {
		lang = DefaultLang
	}
	tmpMsg := strings.ToLower(msg)
	reg := regexp.MustCompile(`\[.*?\]`)
	vars := reg.FindAllString(tmpMsg, -1)
	if len(vars) > 0 {
		for k, v := range vars {
			newStr := "[" + strconv.Itoa(k) + "]"
			tmpMsg = strings.Replace(tmpMsg, v, newStr, 1)
		}
	}
	data := copywriting.New().GetMsgListByLang(msg, lang)
	listData := make(map[string]interface{}, 0)
	newMsg := ""
	if len(data) > 0 {
		for _, v := range data {
			listData[strings.ToLower(v.Action)] = v.Text
		}
		newMsg = listData[tmpMsg].(string)
	}
	if newMsg != "" {
		for k, v := range vars {
			newStr1 := "[" + strconv.Itoa(k) + "]"
			newMsg = strings.Replace(newMsg, newStr1, v, 1)
		}
		newMsg = strings.Replace(newMsg, "[", "", -1)
		newMsg = strings.Replace(newMsg, "]", "", -1)
		return newMsg
	}
	return msg
}
