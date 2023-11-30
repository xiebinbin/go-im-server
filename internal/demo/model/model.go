package model

import (
	"imsdk/internal/common/pkg/config"
	"imsdk/internal/demo/model/signature"
	"imsdk/internal/demo/pkg"
	"imsdk/pkg/funcs"
)

func RequestIMServer(uri, data string) (pkg.CurlResponse, error) {
	nonce := funcs.GetRandString(16)
	timestamp := funcs.GetTimeSecs()
	ak, _ := config.GetConfigAk()
	sk, _ := config.GetConfigSk()
	sign, _ := signature.CreateSign(sk, signature.SignParams{
		AK:        ak,
		Nonce:     nonce,
		Data:      data,
		Timestamp: timestamp,
		Ver:       2,
	})
	params := map[string]interface{}{
		"ak":        ak,
		"data":      data,
		"nonce":     nonce,
		"timestamp": timestamp,
		"signature": sign,
		"ver":       2,
	}
	return pkg.Curl(uri, params)
}
