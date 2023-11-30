package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"imsdk/pkg/funcs"
)

func De(key string, data string) ([]byte, error) {
	var rel []byte
	k := []byte(funcs.StrSha256(key)[:16])
	block, err := aes.NewCipher(k)
	if err != nil {
		return rel, err
	}
	if data[0:2] != "0x" {
		data = "0x" + data
	}
	dataBytes, er1 := hexutil.Decode(data)
	if er1 != nil {
		return rel, er1
	}

	blockMode := cipher.NewCBCDecrypter(block, k)
	blockMode.CryptBlocks(dataBytes, dataBytes)
	dataBytes = pkcs7UnPadding(dataBytes) //和前端代码对应:  padding: CryptoJS.pad.Pkcs7
	return dataBytes, nil
}
