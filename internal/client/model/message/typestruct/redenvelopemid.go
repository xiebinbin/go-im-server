package typestruct

type RedEnvelopeMid struct {
	Operator   string   `json:"operator"`
	Target     []string `json:"target"`
	RecordId   string   `json:"recordId"`
	OutTradeNo string   `json:"outTradeNo"`
	TemId      string   `json:"temId"`
}
