package resource

type Video struct {
	*Resource
	Poster   *Image `json:"poster"`
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	Duration int64  `json:"duration"`
}

//func NewVideo(resource *Resource, poster *Image, width, height, duration int64) *Video {
func NewVideo(byte []byte) *Video {
	resource := NewResource(byte)
	poster := NewImage(byte, IsOriginYes)
	width := int64(100)
	height := int64(100)
	duration := int64(100)
	return &Video{
		resource,
		poster,
		width,
		height,
		duration,
	}
}
