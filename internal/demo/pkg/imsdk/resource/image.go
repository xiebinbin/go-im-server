package resource

import "encoding/json"

type Image struct {
	*Resource
	Width    int64        `json:"width"`
	Height   int64        `json:"height"`
	IsOrigin IsOriginType `json:"isOrigin"`
}

func NewImage(byte []byte, originType IsOriginType) *Image {
	resource := NewResource(byte)
	width := 100
	height := 100
	return &Image{
		Resource: resource,
		Width:    int64(width),
		Height:   int64(height),
		IsOrigin: originType,
	}
}

func NewOriginImage(data map[string]interface{}) *Image {
	var res struct {
		BucketId string       `json:"bucketId"`
		ObjectId string       `json:"text"`
		FileType string       `json:"file_type"`
		Width    int64        `json:"width"`
		Height   int64        `json:"height"`
		IsOrigin IsOriginType `json:"isOrigin"`
	}
	dataByte, _ := json.Marshal(data)
	json.Unmarshal(dataByte, &res)
	resource := &Resource{
		BucketId: res.BucketId,
		ObjectId: res.ObjectId,
		FileType: res.FileType,
	}
	return &Image{
		Resource: resource,
		Width:    res.Width,
		Height:   res.Height,
		IsOrigin: res.IsOrigin,
	}
}
