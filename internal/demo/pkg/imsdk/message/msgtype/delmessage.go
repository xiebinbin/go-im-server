package msgtype

type DelMessage struct {
	TemId    string   `json:"temId"`
	Operator string   `json:"operator"`
	Target   []string `json:"target"`
	MIds     []string `json:"mids"`
}

func (d *DelMessage) DelMessageContent() *DelMessage {
	return d
}
