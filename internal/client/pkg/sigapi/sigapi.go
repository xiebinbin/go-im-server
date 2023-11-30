package sigapi

import (
	"bytes"
	crypto2 "crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	_ "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	"imsdk/pkg/funcs"
	"strconv"
)

type UserSign struct {
	AK        string `json:"ak,omitempty"`
	AUId      string `json:"auid,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Sign      string `json:"sign,omitempty"`
	Ver       int64  `json:"ver,omitempty"`
}

var sk = "5hNMDvExyAsbMv8PPDwCzEVwx62JSZ2NjpxmTw9pgtig"
var ver int64 = 1

func GenUserSig(ak, auid, nonce string, timestamp int64) (string, error) {
	hashByte := hash256(ak, auid, nonce, timestamp, ver)
	sign, _ := SignByPrivateKeyStr(hashByte, sk)
	return base58.Encode(sign), nil
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

func hash256(ak, auid, nonce string, timestamp, ver int64) []byte {
	var contentToBeSigned string
	contentToBeSigned = "ak=" + ak
	contentToBeSigned += "&auid=" + auid
	contentToBeSigned += "&nonce=" + nonce
	contentToBeSigned += "&timestamp=" + strconv.Itoa(int(timestamp))
	contentToBeSigned += "&ver=" + strconv.Itoa(int(ver))
	return funcs.Hash256([]byte(contentToBeSigned))
}

func (u UserSign) sign(key string) string {
	var sb bytes.Buffer
	sb.WriteString("ak:")
	sb.WriteString(u.AK)
	sb.WriteString("\n")
	sb.WriteString("auid:")
	sb.WriteString(u.AUId)
	sb.WriteString("\n")
	sb.WriteString("nonce:")
	sb.WriteString(u.Nonce)
	sb.WriteString("\n")
	sb.WriteString("timestamp:")
	sb.WriteString(strconv.FormatInt(u.Timestamp, 10))
	sb.WriteString("\n")
	sb.WriteString("ver:")
	sb.WriteString(strconv.FormatInt(u.Ver, 10))
	sb.WriteString("\n")
	h := hmac.New(sha256.New, []byte(key))
	if key == "" {
		h = sha256.New()
	}
	h.Write(sb.Bytes())
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
