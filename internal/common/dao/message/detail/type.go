package detail

import "encoding/json"

type Type1Struct struct {
	Text string `json:"text"`
}

type Type2Struct struct {
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	IsOrigin int8   `json:"isOrigin"`
	Size     int64  `json:"size"`
	BucketId string `json:"bucketId"`
	FileType string `json:"fileType"`
	ObjectId string `json:"objectId"`
}

type Type5Struct struct {
	FileType string `json:"fileType"`
	BucketId string `json:"bucketId"`
	ObjectId string `json:"objectId"`
	Size     int64  `json:"size"`
	Name     string `json:"width"`
}

type Type15Struct struct {
	Items []MsgContentType15 `json:"items"`
}
type MsgContentType15 struct {
	Type string `json:"t"`
	Val  string `json:"v"`
}

type Type19Struct struct {
	Icon    TypeIcon        `json:"icon"`
	Text    []Type19Text    `json:"text"`
	Buttons []Type19Buttons `json:"buttons"`
}

type TypeIcon struct {
	BucketId string `json:"bucketId"`
	FileType string `json:"file_type"`
	Height   int64  `json:"height"`
	Width    int64  `json:"width"`
	Text     string `json:"text"`
}

type Type19Text struct {
	T     string `json:"t"`
	Color string `json:"color"`
	Value string `json:"value"`
}

type Type19Buttons struct {
	TXT          string `json:"txt"`
	EnableColor  string `json:"enableColor"`
	DisableColor string `json:"disableColor"`
	ButtonId     string `json:"buttonId"`
}

type Type20Struct struct {
	TemId    string   `json:"temId"`
	Operator string   `json:"operator"`
	Target   []string `json:"target"`
	Duration int64    `json:"duration"`
	Number   int64    `json:"number"`
}

type Type21Struct struct {
	BId  string       `json:"bid"`
	Data []Type21Data `json:"data"`
}

type Type21Data struct {
	Lang  string `json:"lang"`
	Title string `json:"title"`
	TXT   string `json:"txt"`
}

type Type22Struct struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func ParseContentType1(msgContent string) Type1Struct {
	var content Type1Struct
	err := json.Unmarshal([]byte(msgContent), &content)
	if err != nil {
		return Type1Struct{}
	}
	return content
}

func ParseContentType15(msgContent string) Type15Struct {
	var content Type15Struct
	err := json.Unmarshal([]byte(msgContent), &content)
	if err != nil {
		return Type15Struct{}
	}
	return content
}

func ParseContentType19(msgContent string) Type19Struct {
	var content Type19Struct
	err := json.Unmarshal([]byte(msgContent), &content)
	if err != nil {
		return Type19Struct{}
	}
	return content
}

func ParseContentType20(msgContent string) Type20Struct {
	var content Type20Struct
	err := json.Unmarshal([]byte(msgContent), &content)
	if err != nil {
		return Type20Struct{}
	}
	return content
}

func ParseContentType21(msgContent string) Type21Struct {
	var content Type21Struct
	err := json.Unmarshal([]byte(msgContent), &content)
	if err != nil {
		return Type21Struct{}
	}
	return content
}

