package message

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"imsdk/internal/common/dao/chat/members"
	user2 "imsdk/internal/common/dao/user"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/config"
	"imsdk/internal/demo/model"
	"imsdk/internal/demo/pkg"
	"imsdk/internal/demo/pkg/imsdk"
	imsdkMessage "imsdk/internal/demo/pkg/imsdk/message"
	"imsdk/internal/demo/pkg/imsdk/message/msgtype"
	"imsdk/internal/demo/pkg/imsdk/message/options"
	"imsdk/internal/demo/pkg/imsdk/resource"
	"imsdk/pkg/funcs"
	"time"
)

type SendMessageParams struct {
	AChatId    string      `json:"achat_id"`
	AMid       string      `json:"amid"`
	SendTime   int64       `json:"send_time"`
	SenderId   string      `json:"sender_id"`
	Text       string      `json:"text"`
	ReceiveIds []string    `json:"receive_ids"`
	Content    interface{} `json:"content"`
	Extra      string      `json:"extra"`
	Type       int8        `json:"type"`
}

type SendCardAndTempMessageParams struct {
	AChatId  string `json:"achat_id"`
	AMid     string `json:"amid"`
	SenderId string `json:"sender_id"`
	SendTime int64  `json:"send_time"`
}

type GetMessageInfoParams struct {
	AMIds []string `json:"amids"`
}

type GetMessageListParams struct {
	Number int64 `json:"number"`
}

type SetDisableParams struct {
	AMId      string   `json:"amid"`
	AUId      string   `json:"auid"`
	AChatId   string   `json:"achat_id"`
	ButtonIds []string `json:"button_ids"`
}

func NewMessageClient() *imsdkMessage.Client {
	conf, _ := config.GetIMSdkKey()
	messageClient := imsdkMessage.NewClient(&options.Options{
		Credentials: imsdk.NewStaticCredentials(conf.AK, conf.SK),
		Model:       funcs.GetEnv(),
	})
	return messageClient
}

func getAk() string {
	ak, _ := config.GetConfigAk()
	return ak
}

func SendVerticalCardMessageV2(ctx *gin.Context, params SendMessageParams) error {
	demoCon := CreateType23Content()
	con, _ := json.Marshal(demoCon)
	params.AMid = funcs.GetRandString(6)
	params = SendMessageParams{
		AChatId:    params.AChatId,
		AMid:       params.AMid,
		Content:    string(con),
		SendTime:   funcs.GetMillis(),
		SenderId:   params.SenderId,
		ReceiveIds: params.ReceiveIds,
	}
	fmt.Println("tempInfo----", params.AMid)
	buttons := make([]msgtype.VerticalCardButtonsItems, 0)
	if len(demoCon.Buttons) > 0 {
		for _, button := range demoCon.Buttons {
			buttons = append(buttons, msgtype.VerticalCardButtonsItems{
				TXT:          button.TXT,
				EnableColor:  button.EnableColor,
				DisableColor: button.DisableColor,
				ButtonId:     button.ButtonId,
			})
		}
	}
	content := msgtype.NewCardVertical(msgtype.VerticalCardTitle{
		Type:  demoCon.Title.Type,
		Color: demoCon.Title.Color,
		Value: demoCon.Title.Value,
	}, msgtype.VerticalCardText{
		Type:  demoCon.Text.Type,
		Color: demoCon.Text.Color,
		Value: demoCon.Text.Value,
	},
		resource.NewOriginImage(map[string]interface{}{}),
		msgtype.NewCardVerticalButtons(buttons...),
	)
	_, err := NewMessageClient().SendVerticalCardMessage(ctx, params.AChatId, params.AMid, content, params.SenderId, params.ReceiveIds...)
	if err != nil {
		return err
	}
	return nil
}

func SendMessage(ctx context.Context, params SendMessageParams) {
	//sequence, err := NewMessageClient().SendMessage(ctx, imsdkMessage.SendParams{
	//	SendBase: imsdkMessage.NewSendBase(params.AChatId, params.AMid, params.Content.(string), params.SenderId, msgtype.MessageType(params.Type), params.ReceiveIds...),
	//SendUid:    params.SenderId,
	//Amid:       params.AMid,
	//Type:       int8(params.Type),
	//AChatId:    params.ChatId,
	//SendTime:   params.SendTime,
	//ReceiveIds: params.ReceiveIds,
	//Content:    params.Content,
	//Action:     params.Action,
	//Extra: params.Extra,
	//})
}
func SendCardAndTempMessage(ctx *gin.Context, params SendCardAndTempMessageParams) error {
	sendTime := params.SendTime
	if params.SendTime == 0 {
		sendTime = funcs.GetMillis()
	}
	if params.AMid == "" {
		params.AMid = funcs.GetRandString(12)
	}
	err := SendMiddleMessage(ctx, SendMessageParams{
		AChatId:    params.AChatId,
		AMid:       params.AMid,
		SendTime:   sendTime,
		SenderId:   params.SenderId,
		ReceiveIds: []string{params.SenderId},
	})

	otherUId := GetChatMembers(params.AChatId, params.SenderId)
	otherAUId := GetAUId(otherUId)
	if err == nil {
		SendCardMessage(ctx, SendMessageParams{
			AChatId:    params.AChatId,
			SenderId:   params.SenderId,
			AMid:       funcs.GetRandString(16),
			SendTime:   sendTime,
			ReceiveIds: []string{otherAUId},
		})
		fmt.Println("SendCardMessage-----", otherAUId)
	}
	fmt.Println("achatId-----", params.AChatId)
	fmt.Println("sender-auid-----", params.SenderId)
	fmt.Println("senderId-----", base.CreateUId(params.SenderId, getAk()))
	fmt.Println("otherUId-----", otherUId)
	fmt.Println("otherAUId-----", otherAUId)
	return nil
}

func SendTextMessage(ctx *gin.Context, params SendMessageParams) error {
	content := params.Content
	content = CreateType1Content(params.Text)
	con, _ := json.Marshal(content)
	if params.AMid == "" {
		params.AMid = funcs.GetRandString(16)
	}
	params = SendMessageParams{
		AChatId:    params.AChatId,
		AMid:       params.AMid,
		Content:    string(con),
		SendTime:   funcs.GetMillis(),
		SenderId:   params.SenderId,
		ReceiveIds: params.ReceiveIds,
	}
	dataByte, _ := json.Marshal(params)
	fmt.Println("string(dataByte)-----", string(dataByte))
	_, err := model.RequestIMServer("sendTextMessage", string(dataByte))
	//ctx.Set("ak", model.AK)
	//ctx.Set("data", string(dataByte))
	//err := message2.SendTextMessage(ctx)
	if err != nil {
		return err
	}
	return nil
}
func SendImageMessage(ctx *gin.Context, params SendMessageParams) error {
	content := CreateType2Content()
	if params.AMid == "" {
		params.AMid = funcs.GetRandString(16)
	}
	params = SendMessageParams{
		AChatId:    params.AChatId,
		AMid:       params.AMid,
		Content:    content,
		SendTime:   funcs.GetMillis(),
		SenderId:   params.SenderId,
		ReceiveIds: params.ReceiveIds,
	}
	dataByte, _ := json.Marshal(params)
	fmt.Println("string(dataByte)-----", string(dataByte))
	_, err := model.RequestIMServer("sendImageMessage", string(dataByte))
	//ctx.Set("ak", getAk)
	//ctx.Set("data", string(dataByte))
	//err := message2.SendCardMessage(ctx)
	if err != nil {
		return err
	}
	return nil
}

func SendCardMessage(ctx *gin.Context, params SendMessageParams) error {
	content := CreateType19Content()
	con, _ := json.Marshal(content)
	//params.AMid = funcs.GetRandString(6)
	params = SendMessageParams{
		AChatId:    params.AChatId,
		AMid:       params.AMid,
		Content:    string(con),
		SendTime:   funcs.GetMillis(),
		SenderId:   params.SenderId,
		ReceiveIds: params.ReceiveIds,
	}
	fmt.Println("tempInfo----", params.AMid)
	dataByte, _ := json.Marshal(params)
	fmt.Println("string(dataByte)-----", string(dataByte))
	_, err := model.RequestIMServer("sendCardMessage", string(dataByte))
	//ctx.Set("ak", model.AK)
	//ctx.Set("data", string(dataByte))
	//err := message2.SendCardMessage(ctx)
	if err != nil {
		return err
	}
	return nil
}

func SendCardMessageV2(ctx *gin.Context, params SendMessageParams) error {
	content := CreateType19Content()
	con, _ := json.Marshal(content)
	//params.AMid = funcs.GetRandString(6)
	params = SendMessageParams{
		AChatId:    params.AChatId,
		AMid:       params.AMid,
		Content:    string(con),
		SendTime:   funcs.GetMillis(),
		SenderId:   params.SenderId,
		ReceiveIds: params.ReceiveIds,
	}
	fmt.Println("tempInfo----", params.AMid)
	dataByte, _ := json.Marshal(params)
	fmt.Println("string(dataByte)-----", string(dataByte))
	_, err := model.RequestIMServer("sendCardMessage", string(dataByte))
	//ctx.Set("ak", model.AK)
	//ctx.Set("data", string(dataByte))
	//err := message2.SendCardMessage(ctx)
	if err != nil {
		return err
	}
	return nil
}

func SendVerticalCardMessage(ctx *gin.Context, params SendMessageParams) error {
	content := CreateType23Content()
	con, _ := json.Marshal(content)
	params.AMid = funcs.GetRandString(6)
	params = SendMessageParams{
		AChatId:    params.AChatId,
		AMid:       params.AMid,
		Content:    string(con),
		SendTime:   funcs.GetMillis(),
		SenderId:   params.SenderId,
		ReceiveIds: params.ReceiveIds,
	}
	fmt.Println("tempInfo----", params.AMid)
	dataByte, _ := json.Marshal(params)
	fmt.Println("string(dataByte)-----", string(dataByte))
	_, err := model.RequestIMServer("sendVerticalCardMessage", string(dataByte))
	//ctx.Set("ak", model.AK)
	//ctx.Set("data", string(dataByte))
	//err := message2.SendCardMessage(ctx)
	if err != nil {
		return err
	}
	return nil
}

func SendMiddleMessage(ctx *gin.Context, params SendMessageParams) error {
	tempInfo := CreateType20TempInfo(params.SenderId, params.ReceiveIds)
	con, _ := json.Marshal(tempInfo)
	//params.SenderId, _ = config.GetHollowUId()
	params.AMid = funcs.GetRandString(6)
	params = SendMessageParams{
		AChatId:  params.AChatId,
		AMid:     params.AMid,
		Content:  string(con),
		SendTime: funcs.GetMillis(),
		//SenderId:   params.SenderId,
		ReceiveIds: params.ReceiveIds,
	}
	dataByte, _ := json.Marshal(params)
	fmt.Println("string(dataByte)-----", string(dataByte))
	_, err := model.RequestIMServer("sendMiddleMessage", string(dataByte))
	//ctx.Set("ak", getAk)
	//ctx.Set("data", string(dataByte))
	//err := message2.SendMiddleMessage(ctx)
	if err != nil {
		return err
	}
	return nil

}

func SendCustomizeMessage(ctx *gin.Context, params SendMessageParams) error {
	tempInfo := CreateType22TempInfo()
	con, _ := json.Marshal(tempInfo)
	params.SendTime = funcs.GetTimeSecs()
	t := time.Unix(params.SendTime, 0)
	timeUnix, _ := time.Parse("2006-01-02 15:04:05", t.Format("2006-01-02 15:04:05"))
	params = SendMessageParams{
		AChatId:  params.AChatId,
		AMid:     params.AMid,
		Content:  string(con),
		SendTime: timeUnix.UnixNano() / 1e6,
		//SenderId:   params.SenderId,
		ReceiveIds: params.ReceiveIds,
	}
	dataByte, _ := json.Marshal(params)
	fmt.Println("string(dataByte)-----", string(dataByte))
	_, err := model.RequestIMServer("sendCustomizeMessage", string(dataByte))
	//ctx.Set("ak", getAk)
	//ctx.Set("data", string(dataByte))
	//err := message2.SendCustomizeMessage(ctx)
	if err != nil {
		return err
	}
	return nil

}

func SendNotificationMessage(ctx *gin.Context, params SendMessageParams) error {
	return nil
	notification := CreateType21Notification()
	//params.ReceiveIds = []string{"397ec4dbc65a",
	//	"332a5680d7dd", "bb33385154069f3c", "4e7feecdecdaac2a"}
	notificationByte, _ := json.Marshal(notification)
	aChatId := params.AChatId
	if aChatId == "" {
		aChatId = "notice-01"
	}
	params = SendMessageParams{
		AChatId:    aChatId,
		AMid:       params.AMid,
		Content:    string(notificationByte),
		SendTime:   funcs.GetMillis(),
		ReceiveIds: params.ReceiveIds,
	}
	dataByte, _ := json.Marshal(params)
	_, err := model.RequestIMServer("sendNotificationMessage", string(dataByte))
	//ctx.Set("ak", model.AK)
	//ctx.Set("data", string(dataByte))
	//err := message2.SendNotificationMessage(ctx)
	if err != nil {
		return err
	}
	return nil

}

func GetMessageInfo(ctx *gin.Context, params GetMessageInfoParams) (pkg.CurlResponse, error) {
	params = GetMessageInfoParams{
		AMIds: params.AMIds,
	}
	dataByte, _ := json.Marshal(params)
	return model.RequestIMServer("getMessageInfo", string(dataByte))
}

func GetMessageList(ctx *gin.Context, params GetMessageListParams) (pkg.CurlResponse, error) {
	params = GetMessageListParams{
		Number: params.Number,
	}
	dataByte, _ := json.Marshal(params)
	return model.RequestIMServer("getMessageList", string(dataByte))
}

func SetMessageDisable(ctx *gin.Context, params SetDisableParams) (pkg.CurlResponse, error) {
	params = SetDisableParams{
		AMId:      params.AMId,
		AUId:      params.AUId,
		ButtonIds: params.ButtonIds,
	}
	dataByte, _ := json.Marshal(params)
	return model.RequestIMServer("setMessageDisable", string(dataByte))
	//ctx.Set("ak", getAk)
	//ctx.Set("data", string(dataByte))
	//err := message2.SetMessageDisable(ctx)
	//return pkg.CurlResponse{}, err
}

func GetChatMembers(aChatId, senderId string) string {
	chatId := base.CreateChatId(aChatId, getAk())
	senderId = base.CreateUId(senderId, getAk())
	info, _ := members.New().GetChatOtherMemberId(chatId, senderId)
	return info.UID
}

func GetAUId(uid string) string {
	info, _ := user2.New().GetInfoById(uid)
	return info.ID
}
