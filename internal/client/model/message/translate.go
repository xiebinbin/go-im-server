package message

import (
	"imsdk/internal/common/dao/message/detail"
	"imsdk/pkg/app"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/translate/google"
	"regexp"
	"strings"
)

type TranslateParams struct {
	Text string `json:"text"`
	MId  string `json:"mid"`
	Lang string `json:"lang"`
}

const (
	ErrorCodeTrans  = 200000
	ErrorCodeSource = 200001
	ErrorCodeTarget = 200002
)

func DealTransText(text string) string {
	reg := regexp.MustCompile(`\[.*?]`)
	vars := reg.FindAllString(text, -1)
	if len(vars) == 0 {
		return text
	}
	for _, v := range vars {
		text = strings.Replace(text, v, "{0}", -1)
	}
	return text
}

func Translate(params TranslateParams, target string) (string, error) {
	if params.Lang != "" {
		target = params.Lang
	}
	changeTarget := map[string]string{
		"cn":      "zh-CN",
		"zh-hant": "zh-TW",
		"fil":     "tl",
	}
	if changeTarget[target] != "" {
		target = changeTarget[target]
	}
	var configLanguage struct {
		UsableTarget []string `toml:"usable" json:"usable"`
	}
	err := app.Config().Bind("translation", "language", &configLanguage)
	if err != nil {
		return "", err
	}
	if !funcs.In(target, configLanguage.UsableTarget) {
		return "", errno.Add("Target language is not supported", ErrorCodeTarget)
	}
	text := DealTransText(params.Text)
	if strings.Trim(text, " ") == "" {
		return "", nil
	}
	transData := google.TransData{
		Text:       text,
		TargetLang: target,
	}
	res, err := google.TranslV2(transData)
	if res == "" {
		return "", errno.Add("Source language is not supported", ErrorCodeTarget)
	}
	if err != nil {
		return "", errno.Add("Translation failed", ErrorCodeTrans)
	}
	return res, err
}

func isTrans(str string) bool {
	rs := false
	if str != "" {
		reg := regexp.MustCompile(`\{([0-9]|@)}|[《》？：“”【】、；‘'，。、]+|[~!@#$%^&*()_\-+=<>?:"{}\s]+`)
		vars := reg.FindAllString(str, -1)
		if len(vars) > 0 {
			for _, v := range vars {
				str = strings.Replace(str, v, " ", -1)
			}
			if strings.Trim(str, "") != "" {
				rs = true
			}
		}
		rs = true
	}
	return rs
}

func returnStrType1(content, transStr string) string {
	//contentStruct := detail.ParseContentType1(content)
	reg := regexp.MustCompile(`\[.*?]`)
	vars := reg.FindAllString(content, -1)
	if len(vars) == 0 {
		return transStr
	}
	for _, v := range vars {
		//oldStr := "{" + strconv.Itoa(k) + "}"
		transStr = strings.Replace(transStr, "{0}", v, 1)
	}
	return transStr
}

func returnStrType15(content, transStr string) []detail.MsgContentType15 {
	var data []detail.MsgContentType15
	contentStruct := detail.ParseContentType15(content)
	transArr := strings.Split(transStr, "{@}")
	var returnArr []string
	for _, v := range transArr {
		if v != "" {
			returnArr = append(returnArr, v)
		}
	}
	var i = 0
	for _, v := range contentStruct.Items {
		val := v.Val
		if v.Type == "t" && val != "" {
			if len(returnArr) > 0 && i < len(returnArr) {
				val = returnStrType1(val, returnArr[i])
				i++
			}
		}
		tmp := detail.MsgContentType15{
			Type: v.Type,
			Val:  val,
		}
		data = append(data, tmp)
	}
	return data
}
