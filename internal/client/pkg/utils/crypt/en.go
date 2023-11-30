package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"imsdk/pkg/funcs"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func En(key string, data []byte) (string, error) {

	var rel string
	var err error
	if len(key) == 0 {
		err = errors.New("key 不能为空")
	}
	if err != nil {
		return rel, err
	}

	k := []byte(funcs.StrSha256(key)[:16])
	block, er := aes.NewCipher(k)
	if er != nil {
		return rel, er
	}
	blockSize := block.BlockSize()
	data = pkcs7Padding(data, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k)
	cryted := make([]byte, len(data))
	blockMode.CryptBlocks(cryted, data)
	rel = hexutil.Encode(cryted)
	return rel, nil
}
