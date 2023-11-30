package walletutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySign(content string, signatureHex string, address string) (bool, error) {
	if signatureHex[0:2] != "0x" {
		signatureHex = "0x" + signatureHex
	}
	if address[0:2] != "0x" {
		address = "0x" + address
	}
	var rel bool
	var err error
	content = fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(content), content)
	data := []byte(content)
	hash := crypto.Keccak256Hash(data)
	signature := hexutil.MustDecode(signatureHex)

	if signature[64] != 27 && signature[64] != 28 {
		err = errors.New("invalid Ethereum signature (V is not 27 or 28)")
	}
	if err == nil {
		signature[64] -= 27
		pk, er := crypto.SigToPub(hash.Bytes(), signature)
		if er == nil {
			rel = strings.ToLower(crypto.PubkeyToAddress(*pk).Hex()) == address
		}

	}
	return rel, err
}
