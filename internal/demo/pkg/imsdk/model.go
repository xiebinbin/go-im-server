package imsdk

type ModelType = string

const (
	ModelDebug   ModelType = "debug"
	ModelRelease           = "release"

	OssHostRelease string = ""
	OssHostDebug          = ""

	IMHostRelease = ""
	IMHostDebug   = ""
)
