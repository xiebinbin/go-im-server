package signature

import (
	"github.com/btcsuite/btcutil/base58"
	"imsdk/internal/demo/pkg/imsdk"
	"imsdk/internal/demo/pkg/imsdk/signature/eccsign"
	"imsdk/internal/demo/pkg/imsdk/utils"
	"strconv"
)

type SignParams struct {
	AK        string `json:"ak"`
	AUId      string `json:"auid"`
	Data      string `json:"data"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Ver       int8   `json:"ver"`
}

func Sign(data []byte, credentials imsdk.Credentials) ([]byte, error) {
	return eccsign.SignByPrivateKeyStr(data, credentials.Sk)
}

func hash256(ak, auid, nonce string, timestamp int64, ver int8) []byte {
	var contentToBeSigned string
	contentToBeSigned = "ak=" + ak
	contentToBeSigned += "&auid=" + auid
	contentToBeSigned += "&nonce=" + nonce
	contentToBeSigned += "&timestamp=" + strconv.Itoa(int(timestamp))
	contentToBeSigned += "&ver=" + strconv.Itoa(int(ver))
	return utils.Hash256([]byte(contentToBeSigned))
}

func hash256V2(ak, data, nonce string, timestamp int64, ver int8) []byte {
	var contentToBeSigned string
	contentToBeSigned = "ak=" + ak
	contentToBeSigned += "&data=" + data
	contentToBeSigned += "&nonce=" + nonce
	contentToBeSigned += "&timestamp=" + strconv.Itoa(int(timestamp))
	contentToBeSigned += "&ver=" + strconv.Itoa(int(ver))
	return utils.Hash256([]byte(contentToBeSigned))
}

func CreateSign(sk string, params SignParams) (string, error) {
	hashByte := hash256(params.AK, params.AUId, params.Nonce, params.Timestamp, params.Ver)
	if params.Ver == 2 {
		hashByte = hash256V2(params.AK, params.Data, params.Nonce, params.Timestamp, params.Ver)
	}
	sign, _ := eccsign.SignByPrivateKeyStr(hashByte, sk)
	return base58.Encode(sign), nil
}

//func Sign(data, nonce string, timestamp int64, credentials imsdk.Credentials) ([]byte, error) {
//	ver := 2
//	tmp := "ak=%s&data=%s&nonce=%s&timestamp=%d&ver=%d"
//	tmp = fmt.Sprintf(tmp, credentials.Ak, data, nonce, timestamp, ver)
//	return eccsign.SignByPrivateKeyStr([]byte(tmp), credentials.Sk)
//}
