package msgtype

type ReadReport struct {
	MIds []string `json:"mids"`
}

func (r *ReadReport) ReadReportContent() *ReadReport {
	return r
}
