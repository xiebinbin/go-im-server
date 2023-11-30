package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"imsdk/internal/demo/pkg/imsdk/request"
	"imsdk/internal/demo/pkg/imsdk/signature"
	"imsdk/internal/demo/pkg/imsdk/utils"
)

type CurlRequest struct {
	Action actionType `json:"action"`
	Data   string     `json:"data"`
	Ver    int8       `json:"ver"`
}

type CurlResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ErrCode int         `json:"err_code"`
		ErrMsg  string      `json:"err_msg"`
		Items   interface{} `json:"items"`
	} `json:"data"`
}

func NewRequest(action actionType, data string) *CurlRequest {
	return &CurlRequest{
		Action: action,
		Data:   data,
		Ver:    SignatureVersion,
	}
}

func RequestIMServer(ctx context.Context, options *Options, curlRequest *CurlRequest) (CurlResponse, error) {
	ak, sk := options.Credentials.GetAK(), options.Credentials.GetSK()
	nonce, timestamp := utils.GetRandString(16), utils.GetTimeSecs()
	sign, _ := signature.CreateSign(sk, signature.SignParams{
		AK:        ak,
		Nonce:     nonce,
		Timestamp: timestamp,
		Data:      curlRequest.Data,
		Ver:       curlRequest.Ver,
	})

	bodyData := map[string]interface{}{
		"ak":        ak,
		"nonce":     nonce,
		"timestamp": timestamp,
		"signature": sign,
		"data":      curlRequest.Data,
		"ver":       curlRequest.Ver,
	}
	bodyByte, _ := json.Marshal(bodyData)
	re, _ := request.New().SetHost(request.GetImsdkServerHost(ctx, options.Model) + curlRequest.Action).SetMethod(request.POST).
		SetContent(bodyByte).Post()
	var res CurlResponse
	err := json.Unmarshal(re, &res)
	fmt.Printf("mode:%v, Curl res ------- %+v, err %+v", options.Model, res, err)
	return res, nil
}
