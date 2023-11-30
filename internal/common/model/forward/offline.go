package forward

import (
	"context"
	"encoding/json"
	"errors"
	"imsdk/internal/common/dao"
	"imsdk/internal/common/dao/chat"
	chatSetting "imsdk/internal/common/dao/chat/setting"
	"imsdk/internal/common/dao/common"
	"imsdk/internal/common/dao/conversation/usergroupsetting"
	"imsdk/internal/common/dao/conversation/usersinglesetting"
	"imsdk/internal/common/dao/message/usermessage"
	"imsdk/internal/common/dao/message/usermessage/unreadincrement"
	"imsdk/internal/common/dao/user/setting"
	"imsdk/internal/common/model/forward/singapore"
	"imsdk/internal/common/model/translate"
	"imsdk/internal/common/model/user/token"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/config"
	"imsdk/internal/common/pkg/req/request"
	"imsdk/pkg/app"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"net/http"
	"strconv"
	"strings"
)

type MsgRequest struct {
	PushType    int64  `json:"push_type,omitempty"`
	MsgType     uint16 `json:"msg_type,omitempty"`
	ChatId      string `json:"chat_id,omitempty"`
	MsgId       string `json:"mid,omitempty"`
	PushUId     string `json:"push_uid,omitempty"`
	SenderId    string `json:"sender_id,omitempty"`
	SenderOs    string `json:"os,omitempty"`
	Title       string `json:"title,omitempty"`
	PushContent string `json:"push_content,omitempty"`
	MsgContent  string `json:"content,omitempty"`
}

type MsgContentBase struct {
	Text string `json:"text"`
}

type MsgContentType13 struct {
	Addr string `json:"addr"`
	Desc string `json:"desc"`
	Lat  string `json:"lat"`
	Lon  string `json:"lon"`
}

type MsgContentType14 struct {
	MtTd       string `json:"mtid"`
	MtType     uint8  `json:"mt"`
	Creator    string `json:"creator"`
	StartTime  int64  `json:"stime"`
	EndTime    int64  `json:"etime"`
	CreateTime int64  `json:"ctime"`
	Res        int8   `json:"res"`
}

type MsgContentType15 struct {
	Type string `json:"t"`
	Val  string `json:"v"`
}

type Type15Items struct {
	Items []MsgContentType15 `json:"items"`
}

type MsgContentType81 struct {
	Operator string   `json:"operator"`
	TemId    string   `json:"temId"`
	Target   []string `json:"target"`
	MIds     []string `json:"mids"`
}

type ExtraInfoForPush struct {
	PushTitle   string `json:"push_title"`
	PushBodyTag string `json:"push_body_tag"`
}

type OfflineRequest struct {
	Title       string `json:"tile"`
	Body        string `json:"body"`
	CollapseID  string `json:"collapse_id"`
	ClickAction string `json:"click_action"`
}

type OfflineMsgDetail struct {
	ID         string                 `json:"mid"`
	Amid       string                 `json:"amid"`
	Type       uint16                 `json:"type"`
	ChatId     string                 `json:"chat_id"`
	SenderId   string                 `json:"sender_id"`
	Content    string                 `json:"content"`
	Extra      string                 `json:"extra"`
	Action     map[string]interface{} `json:"action,omitempty"`
	SendAt     int64                  `json:"send_time"`
	ReceiveIds []string               `json:"receive_ids,omitempty"`
	ExtraInfo  map[string]interface{} `json:"extra_info"`
	CreatedAt  int64                  `json:"create_time"`
	Sequence   int64                  `json:"sequence"`
}

const (
	KeyValMsgItemTypeText = "t"
	KeyValMsgItemTypeUser = "u"
)

type KeyValMsgItem struct {
	Type string      `json:"t"`
	Val  interface{} `json:"v"`
}

type OfflineReqExtra struct {
	Cmd   string      `json:"cmd"`
	Items interface{} `json:"items"`
}

type OfflineReqParams struct {
	Title       []KeyValMsgItem `json:"title"`
	Body        []KeyValMsgItem `json:"body"`
	Badge       int64           `json:"badge"`
	Devices     []string        `json:"devices"`
	Extra       OfflineReqExtra `json:"extra"`
	OfflineData OfflineRequest  `json:"offline"`
	AppCtx      string          `json:"app_ctx"`
	SdkCtx      string          `json:"sdk_ctx"`
	AUid        string          `json:"auid"`
	CollapseID  string          `json:"collapse_id"`
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

func ProcessOfflineMsg(ctx context.Context, val []byte) error {
	var pushParams PushMessageParams
	err := json.Unmarshal(val, &pushParams)
	logCtx := log.WithFields(ctx, map[string]string{"action": "ProcessOfflineMsg"})
	if err != nil {
		log.Logger().Error(logCtx, "failed to decode pushParams , err: ", err)
		return err
	}
	msgItem := pushParams.Data
	log.Logger().Info(logCtx, "ProcessOfflineMsg pushParams:", pushParams)
	log.Logger().Info(logCtx, "ProcessOfflineMsg msgItem:", msgItem, msgItem.Type, "receiver: ", pushParams.Uid)

	if msgItem.SenderId == pushParams.Uid && msgItem.Type != base.MsgTypeMeeting {
		return nil
	}
	collapseID := msgItem.ID
	var revokeIds []string
	isRange := false
	if funcs.NumIn(int(msgItem.Type), base.NotAddUnreadCountType) {
		if int(msgItem.Type) == base.MsgContentTypeHollowDel {
			var content MsgContentType81
			err = json.Unmarshal([]byte(msgItem.Content), &content)
			if (err != nil && len(content.MIds) == 0) || pushParams.Uid == content.Operator {
				return nil
			}
			revokeIds, _ = usermessage.New().GetUnreadIds(pushParams.Uid, content.MIds)
			if len(revokeIds) == 0 {
				return nil
			}
			isRange = true
		}
	}
	// get auid by uid
	noPushOffline, offlineTitle, offlineBody := PackageOfflineData(ctx, pushParams)
	if noPushOffline {
		return nil
	}

	if !pushParams.NoIncrUnread && !funcs.NumIn(int(msgItem.Type), base.NotAddUnreadCountType) {
		err = unreadincrement.New().Add(pushParams.Uid, msgItem.Sequence)
		if !dao.DataIsSaveSuccessfully(err) {
			log.Logger().Warn(logCtx, "insert unread increment unsuccessfully, err: ", err)
		}
	}

	extra := OfflineReqExtra{
		Cmd:   pushParams.Cmd,
		Items: msgItem,
	}
	if pushParams.Cmd == CmdApplicationCmd {
		extra.Items = msgItem.Content
	}
	if pushParams.Cmd != CmdApplicationCmd {
		extraData, _ := json.Marshal(extra)
		extra = OfflineReqExtra{
			Cmd:   CmdSdkContent,
			Items: string(extraData),
		}
	}
	var offline OfflineRequest
	offlineByte, _ := json.Marshal(msgItem.Offline)
	err = json.Unmarshal(offlineByte, &offline)
	if err != nil {
		return err
	}
	sdkCtx := map[string]interface{}{
		"type": msgItem.Type,
	}
	sdkCtxByte, _ := json.Marshal(sdkCtx)
	reqParams := OfflineReqParams{
		Title:       offlineTitle,
		Body:        offlineBody,
		CollapseID:  collapseID,
		Badge:       int64(common.GetAllBadge(pushParams.Uid)),
		Extra:       extra,
		Devices:     pushParams.Devices,
		AppCtx:      pushParams.AppCtx,
		SdkCtx:      string(sdkCtxByte),
		OfflineData: offline,
		//Uid:     pushParams.Uid,
	}
	callBackUrl, _ := config.GetOfflineCallBackUrl()
	if pushParams.Ak == base.AKSingapore {
		receiveIds := []string{pushParams.Uid}
		senderAId := ""
		singapore.Offline(ctx, callBackUrl, singapore.OfflineReqParams{
			Title:      offlineTitle,
			Body:       offlineBody,
			AChatId:    pushParams.ChatInfo.GId,
			SenderId:   senderAId,
			ReceiveIds: receiveIds,
			Type:       msgItem.Type,
			AMId:       msgItem.Amid,
			Extra:      extra.Items,
			Devices:    pushParams.Devices,
		})
		return nil
	}
	var resp CurlResponse
	var response []byte
	var er error
	if len(revokeIds) > 0 && isRange {
		for _, id := range revokeIds {
			reqParams.CollapseID = id
			response, er = request.InnerReq(ctx, callBackUrl, reqParams)
		}
	} else {
		response, er = request.InnerReq(ctx, callBackUrl, reqParams)
	}
	//response, er := request.InnerReq(ctx, callBackUrl, reqParams)
	if er != nil {
		printInfo, _ := json.Marshal(reqParams)
		log.Logger().Info(logCtx, "failed to request offline forward ,data: ", string(printInfo))
		log.Logger().Error(logCtx, "failed to request offline forward , err: ", er)
		return er
	}
	er = json.Unmarshal(response, &resp)
	if er != nil {
		log.Logger().Error(logCtx, "failed to decode offline forward response , err: ", er)
		return er
	}
	if resp.Code != http.StatusOK || resp.Data.ErrCode != 0 {
		log.Logger().Error(logCtx, "failed to forward , err code: ", resp.Data.ErrCode, " , err msg : ", resp.Data.ErrMsg)
		return errors.New("fail")
	}
	return nil
}

func PackageOfflineData(ctx context.Context, data PushMessageParams) (noPushOffline bool, offlineTitle []KeyValMsgItem, offlineBody []KeyValMsgItem) {
	//msgItem := data.Data.(MsgItem)
	msgItem := data.Data
	msgType := msgItem.Type
	chatInfo := data.ChatInfo
	noPushOffline = data.NoPushOffline
	if data.OfflineData.Title != "" && data.OfflineData.Body != "" {
		offlineTitle = []KeyValMsgItem{{Type: KeyValMsgItemTypeText, Val: data.OfflineData.Title}}
		offlineBody = []KeyValMsgItem{{Type: KeyValMsgItemTypeText, Val: data.OfflineData.Body}}
		return
	}
	if data.Cmd != CmdNoticeNewMsg {
		return
	}
	muteSet := chatSetting.New().GetUserChatMute(msgItem.ChatId, data.Uid)
	// todo old-new version start ========
	logCtx := log.WithFields(ctx, map[string]string{"action": "PackageOfflineData:"})
	ak := data.Ak
	log.Logger().Info(logCtx, "IsLatestVersion", data.Uid, token.IsLatestVersion(data.Uid, base.AppDeviceId))
	muteSetOld := uint8(0)
	if ak == app.OfficialAK && !token.IsLatestVersion(data.Uid, base.AppDeviceId) {
		convInfo := strings.Split(msgItem.ChatId, "_")
		muteSettingInfo := make(map[string]uint8)
		if len(convInfo) > 0 && convInfo[0] == "s" {
			muteSettingInfo = usersinglesetting.New().GetMuteSettingsForChatId(msgItem.ChatId)
		}
		if len(convInfo) > 0 && convInfo[0] == "g" {
			muteSettingInfo = usergroupsetting.New().GetMuteSettingsForChatId(msgItem.ChatId, []string{data.Uid})
		}
		if u, ok := muteSettingInfo[data.Uid]; ok {
			muteSet = u
		}
	}
	// todo old-new version end ========
	muteSetOr := muteSet | muteSetOld
	isMute := muteSet == 1
	log.Logger().Info(logCtx, "muteSetOr:", muteSetOr, muteSet, muteSetOld, isMute, data.Uid)
	// get offline forward title
	isOneToOne := chatInfo.Type == chat.TypeSingle
	if int(msgItem.Type) == base.MsgContentTypeHollowDel {
		var content MsgContentType81
		err := json.Unmarshal([]byte(msgItem.Content), &content)
		if err == nil {

		}
	}
	senderAUId := ""
	if isOneToOne {
		offlineTitle = []KeyValMsgItem{{Type: KeyValMsgItemTypeUser, Val: senderAUId}}
	} else {
		offlineTitle = append(offlineTitle, KeyValMsgItem{Type: KeyValMsgItemTypeText, Val: chatInfo.Name})
	}

	var resType15 []MsgContentType15
	var remindUIds []string
	isTransAll := false
	if msgType == base.MsgTypeRemind {
		resType15 = parseContentType15(msgItem.Content)
		for _, v := range resType15 {
			if v.Type == KeyValMsgItemTypeUser && v.Val == "@_all" {
				isTransAll = true
			}
			if v.Type == KeyValMsgItemTypeUser && v.Val != "@_all" {
				remindUIds = append(remindUIds, v.Val)
				if len(remindUIds) > 5 {
					continue
				}
			}
		}
	}
	if !isTransAll && (isMute && !funcs.In(data.Uid, remindUIds)) {
		noPushOffline = true
		return
	}

	// query user's language
	lang := setting.New().GetUserLang(data.Uid)
	msgItem.SenderId = senderAUId
	pushTitle, pushBody := GetPushContentByParseMsgData(data.Uid, lang, msgItem, isOneToOne)
	if pushTitle != nil {
		offlineTitle = pushTitle
	}
	if pushBody != nil {
		offlineBody = pushBody
	}

	if msgType == base.MsgTypeRemind {
		offlineBody = GetPushContentByParseMsgType15(resType15, isTransAll, senderAUId, data.Uid, lang)
	}

	return
}

func GetPushContentByParseMsgData(targetUId, targetLang string, msgItem MsgItem, isOneToOne bool) (pushTitle []KeyValMsgItem, pushBody []KeyValMsgItem) {
	extraInfo := parseExtra(msgItem.Offline)
	var transTag translate.PushTagType
	// the default value of extraInfo.PushBodyTag is empty string
	if extraInfo.PushTitle != "" {
		pushTitle = []KeyValMsgItem{{Type: KeyValMsgItemTypeText, Val: extraInfo.PushTitle}}
	}
	if !isOneToOne {
		pushBody = []KeyValMsgItem{{Type: KeyValMsgItemTypeUser, Val: msgItem.SenderId}}
		pushBody = append(pushBody, KeyValMsgItem{Type: KeyValMsgItemTypeText, Val: ":"})
	}
	transTag = translate.PushTagType(extraInfo.PushBodyTag)
	pushContent := ""
	switch msgItem.Type {
	case base.MsgTypeImage:
		transTag = translate.ImPushImage
		break
	case base.MsgTypeVoice:
		transTag = translate.ImPushVoice
		break
	case base.MsgTypeVideo:
		transTag = translate.ImPushVideo
		break
	case base.MsgTypeAttachment:
		transTag = translate.ImPushAttachment
		break
	case base.MsgTypeCalls:
		transTag = translate.ImPushConversation
		break
	case base.MsgTypeRedPacket:
		transTag = translate.PaymentRedPacket
		break
	case base.MsgTypeTransfer:
		transTag = translate.PaymentTransfer
		break
	case base.MsgTypeMoments:
		transTag = translate.MomentsBase
		break
	case base.MsgContentTypeHollowDel:
		transTag = translate.ImRevokeMessage
		break
	case base.MsgTypeLocation:
		transTag = translate.ImPushLocation
		pushContent = translate.GetTransByTag(transTag, targetLang)
		resType13 := parseContentType13(msgItem.Content)
		pushContent += resType13.Desc
		pushBody = append(pushBody, KeyValMsgItem{Type: KeyValMsgItemTypeText, Val: pushContent})
		transTag = ""
		break
	case base.MsgTypeMeeting:
		pushBody = append(pushBody, GetTransTagType14(targetUId, msgItem.Content)...)
		transTag = ""
		break
	default:
		defaultContent := parseContentBase(msgItem.Content)
		pushBody = append(pushBody, KeyValMsgItem{Type: KeyValMsgItemTypeText, Val: defaultContent.Text})
		transTag = ""
		break
	}
	if transTag != "" {
		pushContent = translate.GetTransByTag(transTag, targetLang)
		pushBody = append(pushBody, KeyValMsgItem{Type: KeyValMsgItemTypeText, Val: pushContent})
	}
	return
}

func GetPushContentByParseMsgType15(resType15 []MsgContentType15, isTransAll bool, senderId, targetId, targetLang string) []KeyValMsgItem {
	allTransStr := ""
	if isTransAll {
		allTransStr = translate.GetTransByTag(translate.ImPushRemindAll, targetLang)
	}
	pushStr := ""
	//res := make([]KeyValMsgItem, 0)
	res := []KeyValMsgItem{{Type: KeyValMsgItemTypeUser, Val: senderId}}
	res = append(res, KeyValMsgItem{Type: KeyValMsgItemTypeText, Val: ":"})
	for _, v1 := range resType15 {
		if v1.Type == KeyValMsgItemTypeUser && v1.Val == "@_all" {
			res = append(res, KeyValMsgItem{Type: KeyValMsgItemTypeText, Val: "@" + allTransStr + " "})
		}
		if v1.Type == KeyValMsgItemTypeText {
			pushStr += v1.Val
			res = append(res, KeyValMsgItem{Type: KeyValMsgItemTypeText, Val: v1.Val})
		}
		if v1.Type == KeyValMsgItemTypeUser && v1.Val != "@_all" {
			res = append(res, KeyValMsgItem{Type: KeyValMsgItemTypeText, Val: "@"})
			res = append(res, KeyValMsgItem{Type: KeyValMsgItemTypeUser, Val: v1.Val})
		}
	}
	return res
}

func GetTransTagType14(targetUId string, msgContent string) []KeyValMsgItem {
	content := parseContentType14(msgContent)
	return []KeyValMsgItem{{Type: KeyValMsgItemTypeText, Val: content.Res}}
}

func parseContentType13(msgContent string) MsgContentType13 {
	var content MsgContentType13
	err := json.Unmarshal([]byte(msgContent), &content)
	if err != nil {
		return MsgContentType13{}
	}
	return content
}

func parseContentType14(msgContent string) MsgContentType14 {
	var content MsgContentType14
	err := json.Unmarshal([]byte(msgContent), &content)
	if err != nil {
		return MsgContentType14{}
	}
	return content
}

func parseContentType15(msgContent string) []MsgContentType15 {
	var content Type15Items
	err := json.Unmarshal([]byte(msgContent), &content)
	if err != nil {
		return []MsgContentType15{}
	}
	return content.Items
}

func parseContentBase(msgContent string) MsgContentBase {
	var content MsgContentBase
	err := json.Unmarshal([]byte(msgContent), &content)
	if err != nil {
		return MsgContentBase{}
	}
	return content
}

func parseExtra(extra map[string]interface{}) ExtraInfoForPush {
	var content ExtraInfoForPush
	byteStr, _ := json.Marshal(extra)
	err := json.Unmarshal(byteStr, &content)
	if err != nil {
		return ExtraInfoForPush{}
	}
	return content
}

func ProcessKafkaOfflineMsg(ctx context.Context, val []byte) error {
	var params PushMessageParams
	err := json.Unmarshal(val, &params)
	if err != nil {
		return err
	}
	callbackUrl, _ := config.GetOfflineCallBackUrl()
	CallBackThirdProviderOffline(ctx, callbackUrl, params)
	return nil
}

func CallBackThirdProviderOffline(ctx context.Context, callbackUrl string, params PushMessageParams) {
	chatInfo, _ := chat.New().GetByID(params.Data.ChatId, "init_count,achat_id")
	token := singapore.GetToken()
	auid := ""
	UIds := make([]string, 0)
	for _, id := range UIds {
		if id == params.Data.SenderId {
			continue
		}
		_, offlineTitle, offlineBody := PackageOfflineData(ctx, PushMessageParams{
			Cmd:      CmdNoticeNewMsg,
			ChatInfo: chatInfo,
			Uid:      id,
			Data: MsgItem{
				ID:      params.Data.ID,
				Content: params.Data.Content,
				Type:    params.Data.Type,
			},
		})
		// ========singapore data=========
		singaporeTitle, _ := json.Marshal(offlineTitle)
		singaporeBody, _ := json.Marshal(offlineBody)
		singaporeReceiveIds, _ := json.Marshal([]string{auid})
		callbackParams := singapore.CallBackOfflineParams{
			AMID:       params.Data.Amid,
			Type:       strconv.Itoa(int(params.Data.Type)),
			AChatId:    chatInfo.GId,
			Title:      string(singaporeTitle),
			Body:       string(singaporeBody),
			ReceiveIds: string(singaporeReceiveIds),
		}
		msg, _ := json.Marshal(callbackParams)
		log.Logger().Info(ctx,
			map[string]interface{}{
				"amid":      params.Data.Amid,
				"msg":       string(msg),
				"startTime": funcs.GetMillis(),
			})
		formData := make(map[string]interface{}, 0)
		err := json.Unmarshal(msg, &formData)
		if err != nil {
			//fmt.Println("CallBackThirdProviderOffline formData:", callbackUrl, params, err)
			return
		}
		headers := map[string]string{
			"Api-Token": token,
		}
		res, err := singapore.PostForm(callbackUrl, formData, headers)
		//var resData CallbackResp
		var resData map[string]interface{}
		json.Unmarshal(res, &resData)
		log.Logger().Info(ctx,
			map[string]interface{}{
				"amid":    params.Data.Amid,
				"resData": resData,
				"endTime": funcs.GetMillis(),
			})
		if err != nil {
			//fmt.Println("CallBackThirdProviderOffline formData:", callbackUrl, params, err)
			return
		}

	}
	return
}
