package msgtype

type Attachment struct {
	FileType string `json:"fileType"`
	BucketId string `json:"bucketId"`
	ObjectId string `json:"objectId"`
	Size     int64  `json:"size"`
	Name     string `json:"width"`
}
