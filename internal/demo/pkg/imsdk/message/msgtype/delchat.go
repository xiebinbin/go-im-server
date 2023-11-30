package msgtype

type DelChat struct {
	ChatId string `json:"chatid"`
}

func (d *DelChat) DelChatContent() *DelChat {
	return d
}
