package typestruct

type Notice struct {
	BId  string       `json:"bid"`
	Data []NoticeData `json:"data"`
}

type NoticeData struct {
	Lang  string `json:"lang"`
	Title string `json:"title"`
	TXT   string `json:"txt"`
}


