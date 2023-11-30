package walletutil

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateKey() (string, error) {
	var hexPriKey string
	priKey, err := crypto.GenerateKey()
	if err == nil {
		hexPriKey = hexutil.Encode(crypto.FromECDSA(priKey))[2:]
	}
	return hexPriKey, err
}
