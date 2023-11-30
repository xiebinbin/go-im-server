package msgtype

type TradeBase struct {
	ID     string `json:"id"`
	Act    int16  `json:"act,omitempty"`
	Amount uint64 `json:"amount"`
}

type coin struct {
	decimal uint8
	amount  uint64 `json:a, flatten`
	name    string
}

type Payment struct {
	ID           string      `json:"id"`
	OutTradeNo   string      `json:"out_trade_no"`
	FromId       string      `json:"from_id,omitempty"`
	ToId         string      `json:"to_id,omitempty"`
	Act          int16       `json:"act,omitempty"`
	Description  string      `json:"description"`
	Amount       uint64      `json:"amount"`
	TotalAmount  uint64      `json:"total_amount"`
	CoinId       string      `json:"coin_id"`
	CoinName     string      `json:"coin_name"`
	Decimal      uint8       `json:"decimal"`
	CreatedAt    int64       `json:"create_time"`
	Associated   []TradeBase `json:"associated,omitempty"`
	BizType      string      `json:"biz_type"`
	WithdrawAddr string      `json:"withdraw_addr"`
}

type SendPayMessage struct {
	TargetId  string                 `json:"target_id"`
	MsgId     string                 `json:"mid"`
	Content   string                 `json:"content"`
	ExtraInfo map[string]interface{} `json:"extra_info"`
}
