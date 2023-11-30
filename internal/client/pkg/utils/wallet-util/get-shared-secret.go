package walletutil

import (
	"crypto/ecdsa"
	"fmt"

	//"github.com/duke-git/lancet/v2/strutil"
	"github.com/dustinxie/ecc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func (w *Wallet) GetSharedSecret(hexPubKey string) (string, error) {
	var hexSharedSecret string
	if hexPubKey[0:2] != "0x" {
		hexPubKey = "0x" + hexPubKey
	}
	pubKeyBytes, err := hexutil.Decode(hexPubKey)
	if err == nil {
		fmt.Println("hexPubKey", hexPubKey)
		x, y := secp256k1.DecompressPubkey(pubKeyBytes)
		fmt.Println("x", x)
		pubKey := &ecdsa.PublicKey{
			Curve: ecc.P256k1(),
			X:     x,
			Y:     y,
		}
		fmt.Println("pubKey", pubKey)
		SharedSecretX, SharedSecretY := ecc.P256k1().ScalarMult(pubKey.X, pubKey.Y, w.priKey.D.Bytes())
		fmt.Println("SharedSecretX", SharedSecretX)
		hexSharedSecret = hexutil.Encode(SharedSecretX.Bytes())[2:] + hexutil.Encode(SharedSecretY.Bytes())[2:]
		fmt.Println("hexSharedSecret", hexSharedSecret)
	}
	return hexSharedSecret, err
}
