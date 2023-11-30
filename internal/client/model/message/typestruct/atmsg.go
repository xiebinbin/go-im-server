package typestruct

type AtMsg struct {
	Items []AtMsgItems `json:"items"`
}
type itemType string

type AtMsgItems struct {
	Type itemType `json:"t"`
	Val  string   `json:"v"`
}

const (
	textT itemType = "t"
	aUidT          = "u"
	allV           = "all"
)

func NewAtMsg(items ...AtMsgItems) *AtMsg {
	return &AtMsg{
		items,
	}
}

func At(aUId string) AtMsgItems {
	return AtMsgItems{
		Type: aUidT,
		Val:  aUId,
	}
}

func AtAll() AtMsgItems {
	return AtMsgItems{
		Type: aUidT,
		Val:  allV,
	}
}

func AtText(txt string) AtMsgItems {
	return AtMsgItems{
		Type: textT,
		Val:  txt,
	}
}
