package typestruct

type RedEnvelope struct {
	Bless      string `json:"bless"`
	CoinId     string `json:"coin_id"`
	CoinName   string `json:"coin_name"`
	FromUid    string `json:"from_uid"`
	OutTradeNo string `json:"out_trade_no"`
	Amount     int64  `json:"amount"`
}
