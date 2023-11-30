package msgtype

type Transfer struct {
	ID          string `json:"id"`
	Amount      int64  `json:"amount"`
	CoinId      string `json:"coin_id"`
	CoinName    string `json:"coin_name"`
	CoinDecimal int64  `json:"coinDecimal"`
	OutTradeNo  string `json:"out_trade_no"`
}

func NewTransfer() *Transfer {
	return &Transfer{

	}
}
