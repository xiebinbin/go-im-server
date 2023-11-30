package walletutil

import (
	"crypto/ecdsa"
	"errors"

	"github.com/dustinxie/ecc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func PubKeyToAddress(hexPubKey string) (string, error) {
	var address string
	//if strutil.Substring(hexPubKey, 0, 2) != "0x" {
	//	hexPubKey = "0x" + hexPubKey
	//}
	pubKeyBytes, err := hexutil.Decode(hexPubKey)
	if err == nil {
		x, y := secp256k1.DecompressPubkey(pubKeyBytes)
		if x != nil && y != nil {
			pubKey := &ecdsa.PublicKey{
				Curve: ecc.P256k1(),
				X:     x,
				Y:     y,
			}
			address = crypto.PubkeyToAddress(*pubKey).Hex()
		} else {
			err = errors.New("pub key error")
		}
	}
	return address, err
}
