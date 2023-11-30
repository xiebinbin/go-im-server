package setting

import (
	"imsdk/internal/common/dao/user/setting"
	"imsdk/pkg/errno"
)

func UpdateLanguage(uid, lang string) (bool, error) {
	err := setting.New().UpdateValue(uid, "Language", lang)
	if err != nil {
		return false, errno.Add("fail", errno.DefErr)
	}
	return true, nil
}
