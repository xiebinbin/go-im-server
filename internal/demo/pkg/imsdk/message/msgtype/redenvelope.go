package msgtype

type RedEnvelope struct {
	Operator   string   `json:"operator"`
	Target     []string `json:"target"`
	RecordId   string   `json:"recordId"`
	OutTradeNo string   `json:"outTradeNo"`
	TemId      string   `json:"temId"`
}

func (r *RedEnvelope) RedEnvelopeContent() *RedEnvelope {
	return r
}
