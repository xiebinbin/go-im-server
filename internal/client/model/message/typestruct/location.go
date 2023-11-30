package typestruct

type Location struct {
	Addr string `json:"addr"`
	Desc string `json:"desc"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
	Zoom string `json:"zoom"`
}
