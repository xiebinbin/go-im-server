package typestruct

type Image struct {
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	IsOrigin int8   `json:"isOrigin"`
	Size     int64  `json:"size"`
	BucketId string `json:"bucketId"`
	FileType string `json:"fileType"`
	ObjectId string `json:"objectId"`
}
