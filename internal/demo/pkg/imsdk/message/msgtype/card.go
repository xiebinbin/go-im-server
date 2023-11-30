package msgtype

import "imsdk/internal/demo/pkg/imsdk/resource"

type Card struct {
	Icon    *resource.Image `json:"icon"`
	Text    []CardText      `json:"text"`
	Buttons []CardButtons   `json:"buttons"`
}

//type CardIcon struct {
//	BucketId string `json:"bucketId"`
//	FileType string `json:"file_type"`
//	Height   int64  `json:"height"`
//	Width    int64  `json:"width"`
//	Text     string `json:"text"`
//}

type CardText struct {
	Text  string `json:"t"`
	Color string `json:"color"`
	Value string `json:"value"`
}

type CardButtons struct {
	TXT          string `json:"txt"`
	EnableColor  string `json:"enableColor"`
	DisableColor string `json:"disableColor"`
	ButtonId     string `json:"buttonId"`
}
