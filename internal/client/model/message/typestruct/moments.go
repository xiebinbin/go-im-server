package typestruct

type Moments struct {
	UID         string       `json:"uid"`
	MomentsId   string       `json:"mid"`
	MomentsText string       `json:"text"`
	Type        int8         `json:"type"`
	ContentType int8         `json:"contentType"`
	Image       ImageContent `json:"image"`
}

type playButtonValue int64

const (
	Hide playButtonValue = 1
	Show playButtonValue = 2
)

type ImageContent struct {
	*Image
	PlayButton playButtonValue `json:"mediaType"`
}
