package typestruct

type DelMessage struct {
	TemId    string   `json:"temId"`
	Operator string   `json:"operator"`
	Target   []string `json:"target"`
	MIds     []string `json:"mids"`
}
