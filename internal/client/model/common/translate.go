package common

import (
	"imsdk/internal/common/dao/setting/copywriting"
)

func GetTransByTag(tag, lang string) string {
	res := copywriting.New().GetText(tag, copywriting.TypePush, lang)
	return res
}
