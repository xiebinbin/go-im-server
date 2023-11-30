package typestruct

type Meeting struct {
	MeetingId   string `json:"mtid"`
	MeetingType int8   `json:"mt"`
	Creator     string `json:"creator"`
	StartTime   int64  `json:"stime"`
	EndTime     int64  `json:"etime"`
	CreateTime  int64  `json:"ctime"`
	ResStatus   int8   `json:"res"`
}
