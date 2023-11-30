package msgtype

type Location struct {
	Addr string `json:"addr"`
	Desc string `json:"desc"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
	Zoom string `json:"zoom"`
}

func (l *Location) LocationContent() *Location {
	return l
}
