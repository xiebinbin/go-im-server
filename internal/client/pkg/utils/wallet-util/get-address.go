package walletutil

import "github.com/ethereum/go-ethereum/crypto"

func (w *Wallet) GetAddress() string {
	return crypto.PubkeyToAddress(*w.pubKey).Hex()
}
