package google

import v2 "imsdk/pkg/translate/google/v2"

type TransData struct {
	Text       string `json:"text"`
	TargetLang string `json:"target_lang"`
}

func TranslV2(data TransData) (string, error) {
	return v2.TranslateText(data.Text, data.TargetLang)
}
