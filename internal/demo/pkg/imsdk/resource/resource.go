package resource

type Resource struct {
	BucketId string `json:"bucketId"`
	//ObjectId string `json:"objectId"`
	//FileType string `json:"fileType"`
	ObjectId string `json:"text"`
	FileType string `json:"file_type"`
}

type IsOriginType uint8

const (
	IsOriginYes IsOriginType = 1
	IsOriginNo               = 0
)

func NewResource(resource []byte) *Resource {
	return &Resource{
		BucketId: "7e4216ac1941c61c",
		ObjectId: "img/e2/e2110bcb2904e50787f4100d371004c8cee71dae",
		FileType: "jpg",
	}
}
