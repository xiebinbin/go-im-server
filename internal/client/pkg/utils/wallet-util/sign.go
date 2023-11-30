package walletutil

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func (w *Wallet) Sign(content string) (string, error) {
	var sign string
	content = fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(content), content)
	hash := crypto.Keccak256Hash([]byte(content))
	signature, err := crypto.Sign(hash.Bytes(), w.priKey)
	if err == nil {
		if signature[64] != 27 && signature[64] != 28 {
			signature[64] += 27
		}
		sign = hexutil.Encode(signature)
	}
	return sign, err
}
