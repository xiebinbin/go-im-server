package signature

import (
	"github.com/btcsuite/btcutil/base58"
	"imsdk/pkg/eccsign"
	"imsdk/pkg/funcs"
	"strconv"
)

type SignParams struct {
	AK        string `json:"ak"`
	AUId      string `json:"auid"`
	Data      string `json:"data"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Ver       int64  `json:"ver"`
}

func hash256(ak, auid, nonce string, timestamp, ver int64) []byte {
	var contentToBeSigned string
	contentToBeSigned = "ak=" + ak
	contentToBeSigned += "&auid=" + auid
	contentToBeSigned += "&nonce=" + nonce
	contentToBeSigned += "&timestamp=" + strconv.Itoa(int(timestamp))
	contentToBeSigned += "&ver=" + strconv.Itoa(int(ver))
	return funcs.Hash256([]byte(contentToBeSigned))
}

func hash256V2(ak, data, nonce string, timestamp, ver int64) []byte {
	var contentToBeSigned string
	contentToBeSigned = "ak=" + ak
	contentToBeSigned += "&data=" + data
	contentToBeSigned += "&nonce=" + nonce
	contentToBeSigned += "&timestamp=" + strconv.Itoa(int(timestamp))
	contentToBeSigned += "&ver=" + strconv.Itoa(int(ver))
	return funcs.Hash256([]byte(contentToBeSigned))
}

func CreateSign(sk string, params SignParams) (string, error) {
	hashByte := hash256(params.AK, params.AUId, params.Nonce, params.Timestamp, params.Ver)
	if params.Ver == 2 {
		hashByte = hash256V2(params.AK, params.Data, params.Nonce, params.Timestamp, params.Ver)
	}
	sign, _ := eccsign.SignByPrivateKeyStr(hashByte, sk)
	return base58.Encode(sign), nil
}
