package walletutil

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	priKey *ecdsa.PrivateKey
	pubKey *ecdsa.PublicKey
}

func New(hexPriKey string) (*Wallet, error) {
	var val *Wallet
	priKey, err := crypto.HexToECDSA(hexPriKey)
	if err == nil {
		pubKey, ok := priKey.Public().(*ecdsa.PublicKey)
		if !ok {
			err = errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}
		if err == nil {
			val = &Wallet{
				priKey: priKey,
				pubKey: pubKey,
			}
		}
	}
	return val, err
}
