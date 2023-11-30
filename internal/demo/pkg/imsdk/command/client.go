package command

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"imsdk/internal/demo/pkg/imsdk"
	"imsdk/internal/demo/pkg/imsdk/request"
	"imsdk/internal/demo/pkg/imsdk/signature"
	"imsdk/internal/demo/pkg/imsdk/utils"
)

type Options struct {
	Model       imsdk.ModelType
	Credentials *imsdk.Credentials
}

type Client struct {
	options      *Options
	SocketParams *SocketParams
}

func NewClient(options *Options) *Client {
	client := &Client{
		options:      options,
		SocketParams: &SocketParams{},
	}
	return client
}

type SocketParams struct {
	Cmd         string      `json:"cmd,omitempty"`
	Items       interface{} `json:"items,omitempty"`
	ReceiveTers []string    `json:"receive_ters,omitempty"`
}

func (c *Client) SetSocketParams(cmd string, items interface{}) *Client {
	c.SocketParams = &SocketParams{
		Cmd:   cmd,
		Items: items,
	}
	return c
}

func (c *Client) GetSocketParams() *SocketParams {
	return c.SocketParams
}

type SendCmdParams struct {
	AUIds         []string      `json:"auids"`
	Data          string        `json:"data"`
	OfflineData   OfflineParams `json:"offline"` // {"title":"","body":""}
	NoPushOffline bool          `json:"no_push_offline"`
	AppCtx        string        `json:"app_ctx"`
}

type OfflineParams struct {
	Title       string `json:"tile"`
	Body        string `json:"body"`
	CollapseID  string `json:"collapse_id"`
	ClickAction string `json:"click_action"`
}

func (c *Client) SendCmd(ctx context.Context, auids []string, offline OfflineParams, noPushOffline bool, appCtx string) {
	nonce := utils.GetRandString(16)
	timestamp := utils.GetTimeSecs()
	dataByte, _ := json.Marshal(c.GetSocketParams())
	cmdDataByte, _ := json.Marshal(SendCmdParams{
		AUIds:         auids,
		Data:          string(dataByte),
		OfflineData:   offline,
		NoPushOffline: noPushOffline,
		AppCtx:        appCtx,
	})
	c.RequestIMServer(ctx, c.options.Credentials, RequestParams{
		Action:    "sendCmd",
		Nonce:     nonce,
		Timestamp: timestamp,
		Ver:       2,
		Data:      string(cmdDataByte),
	})
}

type RequestParams struct {
	Action    string `json:"action"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Ver       int8   `json:"ver"`
	Data      string `json:"data"`
}

func (c *Client) RequestIMServer(ctx context.Context, credentials *imsdk.Credentials, params RequestParams) (CurlResponse, error) {
	AK := credentials.GetAK()
	SK := credentials.GetSK()
	sign, _ := signature.CreateSign(SK, signature.SignParams{
		AK:        AK,
		Nonce:     params.Nonce,
		Data:      params.Data,
		Timestamp: params.Timestamp,
		Ver:       2,
	})
	bodyData := map[string]interface{}{
		"ak":        AK,
		"data":      params.Data,
		"nonce":     params.Nonce,
		"timestamp": params.Timestamp,
		"signature": sign,
		"ver":       2,
	}
	bodyByte, _ := json.Marshal(bodyData)
	re, _ := request.New().SetHost(request.GetImsdkServerHost(ctx, c.options.Model) + params.Action).SetMethod(request.POST).
		SetContent(bodyByte).Post()
	var res CurlResponse
	err := json.Unmarshal(re, &res)
	fmt.Printf("Curl res ------- %+v, err %+v", res, err)
	return res, nil
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
