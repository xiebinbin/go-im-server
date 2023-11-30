package initialization

import (
	_ "github.com/go-playground/locales/de"
	"strings"
)

func replaceNoticeStr(str string) string {
	replace := map[string]string{
		"@{@}":        "{0}",
		"@{operator}": "{1}",
		"@{target}":   "{2}",
		"@{duration}": "{3}",
		"@{luckUid}":  "{4}",
		"@{gift}":     "{5}",
		"@{gift1}":    "{6}",
		"@{gift2}":    "{7}",
		"@{gift3}":    "{8}",
		"@{gift4}":    "{9}",
	}
	for k, v := range replace {
		str = strings.Replace(str, k, v, -1)
	}
	return str
}

func returnNoticeStr(str string) string {
	replace := map[string]string{
		"{1}": "@{operator}",
		"{2}": "@{target}",
		"{3}": "@{duration}",
		"{4}": "@{luckUid}",
		"{5}": "@{gift}",
		"{6}": "@{gift1}",
		"{7}": "@{gift2}",
		"{8}": "@{gift3}",
		"{9}": "@{gift4}",
	}
	for k, v := range replace {
		str = strings.Replace(str, k, v, -1)
	}
	return str
}
