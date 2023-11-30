package walletutil

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func (w *Wallet) GetPubKey() string {
	pubKey := w.priKey.Public().(*ecdsa.PublicKey)
	return hex.EncodeToString(secp256k1.CompressPubkey(pubKey.X, pubKey.Y))
}
