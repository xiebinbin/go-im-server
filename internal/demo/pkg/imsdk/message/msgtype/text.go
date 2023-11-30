package msgtype

type Text struct {
	Text string `json:"text"`
}

func NewText(text string) *Text {
	return &Text{
		Text: text,
	}
}
