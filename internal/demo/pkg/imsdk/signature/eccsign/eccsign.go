package eccsign

import (
	crypto2 "crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"math/big"
)

type Wallets = hdwallet.Wallet

type Account = accounts.Account

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

func VerifySign(message, signature []byte, pubKeyStr string) error {
	//return true
	pubKey, err := GetPublicKeyByStr(pubKeyStr)
	if err != nil {
		return err
	}
	var esig struct {
		R, S *big.Int
	}
	if _, err = asn1.Unmarshal(signature, &esig); err != nil {
		return err
	}
	//message = funcs.Hash256(message)
	if ecdsa.Verify(pubKey, message, esig.R, esig.S) {
		return nil
	}
	return errors.New("failed")
}
