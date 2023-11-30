package resource

type Voice struct {
	*Resource
	Duration int64 `json:"duration"`
	Size     int64 `json:"size"`
}

func NewVoice(byte []byte) *Voice {
	resource := NewResource(byte)
	duration := int64(1000)
	size := int64(1000)
	return &Voice{
		Resource: resource,
		Duration: duration,
		Size:     size,
	}
}
