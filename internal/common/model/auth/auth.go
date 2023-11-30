package auth

import (
	"context"
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
	"imsdk/internal/common/model/errors"
	"imsdk/pkg/eccsign"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/redis"
	"strconv"
	"time"
)

type VerifyParams struct {
	AK        string `json:"ak"`
	AUId      string `json:"auid"`
	Data      string `json:"data"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Ver       int64  `json:"ver"`
	Signature string `json:"signature"`
}

func Verify(ctx context.Context, pk string, request VerifyParams) error {
	logCtx := log.WithFields(ctx, map[string]string{"action": "GetAuth Verify"})
	if err := verifyTimestamp(ctx, request.Timestamp); err != nil {
		log.Logger().Error(logCtx, "Verify Timestamp 1:", request.Timestamp, err)
		return err
	}

	if err := verifyNonce(request.Nonce); err != nil {
		log.Logger().Error(logCtx, "Verify Nonce 2:", request.Nonce, err)
		return err
	}
	if err := verifySign(pk, request); err != nil {
		log.Logger().Error(logCtx, "Verify Sign 3:", pk, request, err)
		return err
	}

	return nil
}

func verifyTimestamp(ctx context.Context, timestamp int64) error {
	logCtx := log.WithFields(ctx, map[string]string{"action": "verifyTimestamp"})
	now := time.Now().Unix()
	if now > timestamp+180 {
		log.Logger().Info(logCtx, now, timestamp)
		return errors.ErrSdkTimeOut
	}
	return nil
}

func verifyNonce(nonce string) error {
	logCtx := log.WithFields(context.Background(), map[string]string{"action": "verifyNonce"})
	cache, er := redis.Client.Get(nonce).Result()
	if er != nil && er != redis.NilErr {
		log.Logger().Error(logCtx, "verifyNonce redis er-:", er)
		return errors.ErrSdkDefErr
	}
	if cache != "" {
		return errors.ErrSdkRequestRepeat
	}
	return nil
}

/**
ak appKey
pk publicKey
*/
func verifySign(pk string, params VerifyParams) error {
	logCtx := log.WithFields(context.Background(), map[string]string{"action": "verifySign"})
	message := hash256(params.AK, params.AUId, params.Nonce, params.Timestamp, params.Ver)
	if params.Ver == 2 {
		message = hash256V2(params.AK, params.Data, params.Nonce, params.Timestamp, params.Ver)
	}
	log.Logger().Info(logCtx, "verifySign:", "message: ", message, "- string: ", hex.EncodeToString(message))
	signature := base58.Decode(params.Signature)
	isSuc := eccsign.VerifySign(message, signature, pk)
	if !isSuc {
		log.Logger().Error(logCtx, "verifySign:", pk, isSuc)
		return errors.ErrSdkSignVerifyFail
	}
	return nil
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

func hash256V2(ak, data, nonce string, timestamp, ver int64) []byte {
	var contentToBeSigned string
	contentToBeSigned = "ak=" + ak
	contentToBeSigned += "&data=" + data
	contentToBeSigned += "&nonce=" + nonce
	contentToBeSigned += "&timestamp=" + strconv.Itoa(int(timestamp))
	contentToBeSigned += "&ver=" + strconv.Itoa(int(ver))
	return funcs.Hash256([]byte(contentToBeSigned))
}

func CreateSign(sk string, params VerifyParams) (string, error) {
	hashByte := hash256(params.AK, params.AUId, params.Nonce, params.Timestamp, params.Ver)
	if params.Ver == 2 {
		hashByte = hash256V2(params.AK, params.Data, params.Nonce, params.Timestamp, params.Ver)
	}
	sign, _ := eccsign.SignByPrivateKeyStr(hashByte, sk)
	return base58.Encode(sign), nil
}
