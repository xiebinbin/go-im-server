package eccsign

import (
	"context"
	crypto2 "crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"math/big"
)

type Wallets = hdwallet.Wallet

type Account = accounts.Account

//seed rule: HMAC_SHA256(uid + password, password)
func GenerateSeed(uid, password string) (string, error) {
	return funcs.HmacSha256(uid+password, password), nil
}

func GetPrivateKey(seed string) (string, error) {
	wallet, err := hdwallet.NewFromSeed([]byte(seed))
	if err != nil {
		return "", err
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", err
	}
	if priKeyByte, err := wallet.PrivateKeyBytes(account); err != nil {
		return "", err
	} else {
		return base58.Encode(priKeyByte), nil
	}
}

func GetPrivateKeyByString(key string) (*ecdsa.PrivateKey, error) {
	return crypto.ToECDSA(base58.Decode(key))
}

func SignByPrivateKeyStr(src []byte, pri string) ([]byte, error) {
	priKey, err := crypto.ToECDSA(base58.Decode(pri))
	if err != nil {
		return nil, err
	}
	return Sign(src, priKey)
}

func Sign(src []byte, priKey *ecdsa.PrivateKey) ([]byte, error) {
	var ops crypto2.SignerOpts
	return priKey.Sign(rand.Reader, src, ops)
}

func GetPublicKeyStr(seed string) (string, error) {
	wallet, err := hdwallet.NewFromSeed([]byte(seed))
	if err != nil {
		return "", err
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", err
	}
	publicKeyByte, err := wallet.PublicKeyBytes(account)
	if err != nil {
		return "", err
	}
	return base58.Encode(publicKeyByte), nil
}

func GetPublicKeyByStr(pubKeyStr string) (*ecdsa.PublicKey, error) {
	pubKeyByte := base58.Decode(pubKeyStr)
	return crypto.UnmarshalPubkey(pubKeyByte)
}

func VerifySign(message, signature []byte, pubKeyStr string) bool {
	//return true
	pubKey, err := GetPublicKeyByStr(pubKeyStr)
	ctx := log.WithFields(context.Background(), map[string]string{"action": "VerifySign"})
	if err != nil {
		log.Logger().Error(ctx, "failed to get public key, err: ", err)
		return false
	}
	var esig struct {
		R, S *big.Int
	}
	if _, err = asn1.Unmarshal(signature, &esig); err != nil {
		log.Logger().Error(ctx, "failed to get unmarshal signature, err: ", err)
		return false
	}
	return ecdsa.Verify(pubKey, message, esig.R, esig.S)
}

func VerifyPassword(pubKeyStr string, adr string) bool {
	//return true
	address, err := GetAddressByPubKeyStr(pubKeyStr)
	if err != nil {
		return false
	}
	return adr == address
}

func GetAddressByPubKeyStr(pubKeyStr string) (string, error) {
	pubKey, err := GetPublicKeyByStr(pubKeyStr)
	if err != nil {
		return "", err
	}
	//adrByte := crypto.PubkeyToAddress(*pubKey).Hex()
	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}
