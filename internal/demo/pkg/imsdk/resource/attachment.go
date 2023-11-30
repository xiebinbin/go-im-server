package resource

type Attachment struct {
	*Resource
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func NewAttachment(byte []byte) *Attachment {
	resource := NewResource(byte)
	name := ""
	size := 0
	return &Attachment{
		Resource: resource,
		Name:     name,
		Size:     int64(size),
	}
}
