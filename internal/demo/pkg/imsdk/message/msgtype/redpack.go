package msgtype

type RedPack struct {
	Bless      string `json:"bless"`
	CoinId     string `json:"coin_id"`
	CoinName   string `json:"coin_name"`
	Amount     int64  `json:"amount"`
	FromUid    string `json:"from_uid"`
	OutTradeNo string `json:"out_trade_no"`
}

func NewRedPack() *RedPack {
	return &RedPack{

	}
}
