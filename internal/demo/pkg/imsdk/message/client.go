package message

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"imsdk/internal/demo/pkg/imsdk"
	_ "imsdk/internal/demo/pkg/imsdk"
	"imsdk/internal/demo/pkg/imsdk/message/msgtype"
	"imsdk/internal/demo/pkg/imsdk/message/options"
	"imsdk/internal/demo/pkg/imsdk/request"
	"imsdk/internal/demo/pkg/imsdk/resource"
	"imsdk/internal/demo/pkg/imsdk/signature"
	"imsdk/internal/demo/pkg/imsdk/utils"
	"sync"
)

type Client struct {
	options *options.Options
}

var (
	once   sync.Once
	client *Client
)

func NewClient(options *options.Options) *Client {
	once.Do(func() {
		client = &Client{
			options: options,
		}
	})
	return client
}

func (c *Client) SendCustomizeMessage(ctx context.Context, aChatId, aMId string, customizeContent *msgtype.Customize, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(customizeContent)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeText, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendTextMessage(ctx context.Context, aChatId, aMId string, text *msgtype.Text, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(text)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeText, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendImageMessage(ctx context.Context, aChatId, aMId string, image *resource.Image, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(image)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeImage, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendVoiceMessage(ctx context.Context, aChatId, aMId string, image *resource.Attachment, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(image)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeVoice, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendVideoMessage(ctx context.Context, aChatId, aMId string, image *resource.Video, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(image)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeVideo, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendAttachmentMessage(ctx context.Context, aChatId, aMId string, attachment *resource.Attachment, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(attachment)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeAttachment, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendNoticeMessage(ctx context.Context, aMId, aChatId string, notice *msgtype.Notice, extra map[string]interface{}, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(notice)
	if err != nil {
		return 0, err
	}
	extraByte, _ := json.Marshal(extra)
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), MiddleMsgSendUID, msgtype.TypeNotice, receiveUid...),
		Extra:    string(extraByte),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendVerticalCardMessage(ctx context.Context, aChatId, aMId string, verticalCard *msgtype.VerticalCard, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(verticalCard)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeVerticalCard, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendRedEnvelopeMessage(ctx context.Context, aChatId, aMId string, text *msgtype.RedEnvelope, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(text)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeRedPacket, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendPaymentMessage(ctx context.Context, aChatId, aMId string, notice *msgtype.Payment, extra map[string]interface{}, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(notice)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase:  NewSendBase(aChatId, aMId, string(contentStr), MiddleMsgSendUID, msgtype.TypePayment, receiveUid...),
		ExtraInfo: extra,
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *Client) SendMeetingMessage(ctx context.Context, aChatId, aMId string, meeting *msgtype.Meeting, senderId string, receiveUid ...string) (uint64, error) {
	contentStr, err := json.Marshal(meeting)
	if err != nil {
		return 0, err
	}
	msg := SendParams{
		SendBase: NewSendBase(aChatId, aMId, string(contentStr), senderId, msgtype.TypeMeeting, receiveUid...),
	}
	_, err = c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

type DisableParams struct {
	Amid      string   `json:"amid"`
	AUIds     []string `json:"auids"`
	ButtonIds []string `json:"button_ids"`
}

func (c *Client) sendMessage(ctx context.Context, msg SendParams) (CurlResponse, error) {
	nonce, timestamp := utils.GetRandString(16), utils.GetTimeSecs()
	createDataByte, _ := json.Marshal(msg)
	res, err := c.RequestIMServer(ctx, c.options.Credentials, RequestParams{
		Action:    SendMessage,
		Nonce:     nonce,
		Timestamp: timestamp,
		Ver:       2,
		Data:      string(createDataByte),
	})
	if err != nil {
		return CurlResponse{}, err
	}
	return res, nil
}

func (c *Client) SetDisable(ctx context.Context, disable DisableParams) (uint64, error) {
	nonce := utils.GetRandString(16)
	timestamp := utils.GetTimeSecs()
	createDataByte, _ := json.Marshal(disable)
	_, err := c.RequestIMServer(ctx, c.options.Credentials, RequestParams{
		Action:    SetMessageDisableMultiUId,
		Nonce:     nonce,
		Timestamp: timestamp,
		Ver:       2,
		Data:      string(createDataByte),
	})
	if err != nil {
		return 0, err
	}
	return 0, nil
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
	params.Ver = 2
	sign, _ := signature.CreateSign(SK, signature.SignParams{
		AK:        AK,
		Nonce:     params.Nonce,
		Data:      params.Data,
		Timestamp: params.Timestamp,
		Ver:       params.Ver,
	})

	bodyData := map[string]interface{}{
		"ak":        AK,
		"data":      params.Data,
		"nonce":     params.Nonce,
		"timestamp": params.Timestamp,
		"signature": sign,
		"ver":       params.Ver,
	}
	fmt.Println("message bodyData:", request.GetImsdkServerHost(ctx, c.options.Model)+params.Action, bodyData)
	bodyByte, _ := json.Marshal(bodyData)
	re, _ := request.New().SetHost(request.GetImsdkServerHost(ctx, c.options.Model) + params.Action).SetMethod(request.POST).
		SetContent(bodyByte).Post()
	var res CurlResponse
	err := json.Unmarshal(re, &res)
	fmt.Printf("Curl res ------- %+v, err %+v", res, err)
	return res, nil
}

type CurlResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
