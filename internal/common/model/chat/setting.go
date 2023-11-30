package chat

import (
	"context"
	"errors"
	json "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao/chat/setting"
	"imsdk/internal/common/model/forward"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"strconv"
)

type HideRequest struct {
	ChatId       string `json:"chat_id" binding:"required"`
	HideSequence int64  `json:"hide_sequence"`
	LocalTime    int64  `json:"local_time"`
}

type IsTopRequest struct {
	ChatId    string `json:"chat_id" binding:"required"`
	IsTop     uint8  `json:"is_top" binding:"oneof=0 1"`
	LocalTime int64  `json:"local_time"`
}

func GetSetting(uid, chatId string) (setting.Setting, error) {
	dao := setting.New()
	data, err := dao.GetSetting(uid, chatId)
	if errors.Is(err, mongo.ErrNoDocuments) { // setting not exist, add default setting
		//data := setting.DefSetting
		data.ID = setting.GetId(uid, chatId)
		data.UID = uid
		data.ChatID = chatId
		_, err = dao.Add(data)
		if err != nil {
			return setting.Setting{}, err
		}
		return data, nil
	} else if err == nil {
		return data, nil
	}
	return data, errno.Add("fail", errno.DefErr)
}

func GetSettings(uid string, chatId []string) ([]setting.Setting, error) {
	dao := setting.New()
	data := dao.GetSettings(uid, chatId)
	uIds := make([]string, 0)
	uIds = append(uIds, uid)
	var existsChatIds []string
	var notExistsChatIds []string
	if len(data) == 0 {
		data, _ = SetDefaultInfoByChatIds(uid, chatId)
	} else if len(data) < len(chatId) {
		for _, v := range data {
			existsChatIds = append(existsChatIds, v.ChatID)
		}
		for _, v := range chatId {
			if !funcs.In(v, existsChatIds) {
				notExistsChatIds = append(notExistsChatIds, v)
			}
		}
		notExistsSetting, _ := SetDefaultInfoByChatIds(uid, notExistsChatIds)
		for _, v := range notExistsSetting {
			data = append(data, v)
		}
	}
	return data, nil
}

func SetDefaultInfoByUIds(uIds []string, chatId string) (interface{}, error) {
	dao := setting.New()
	dataList := make([]setting.Setting, 0)
	for _, v := range uIds {
		data := setting.DefSetting
		data.ID = setting.GetId(v, chatId)
		data.UID = v
		data.ChatID = chatId
		dataList = append(dataList, data)
	}
	res, _ := dao.AddMany(dataList)
	return res, nil
}

func SetDefaultInfoByChatIds(uid string, chatIds []string) ([]setting.Setting, error) {
	dao := setting.New()
	ctx := context.Background()
	dataList := make([]setting.Setting, 0)
	if len(chatIds) > 0 {
		for _, v := range chatIds {
			data := setting.DefSetting
			data.ID = setting.GetId(uid, v)
			data.UID = uid
			data.ChatID = v
			dataList = append(dataList, data)
		}
		insertIds, err := dao.AddMany(dataList)
		if err != nil {
			logCtx := log.WithFields(ctx, map[string]string{"action": "AddMany", "uid": uid})
			log.Logger().Error(logCtx, "AddMany fail", insertIds, err)
			return nil, err
		}
	}
	return dataList, nil
}

func UpdateIsTop(ctx context.Context, uid string, request IsTopRequest) bool {
	chatId, isTop := request.ChatId, request.IsTop
	dao := setting.New()
	settingInfo, err := dao.GetSetting(uid, chatId)
	id := setting.GetId(uid, chatId)
	logCtx := log.WithFields(ctx, map[string]string{"action": "chatUpdateIsTop", "uid": uid, "chatId": chatId, "isTop": strconv.Itoa(int(request.IsTop))})
	topSequence := 0
	if err != nil && err == mongo.ErrNoDocuments {
		_, er := addDefaultSetting(ctx, uid, chatId)
		if er != nil {
			log.Logger().Error(logCtx, "failed to update ,err : ", er)
			return false
		}
	} else {
		topSequence = int(settingInfo.TopSequence)
	}
	topTime := 0
	if isTop == 1 {
		topTime = int(funcs.GetMillis())
	}

	err = dao.UpdateByMap(id, map[string]interface{}{
		"is_top":       isTop,
		"top_time":     topTime,
		"top_sequence": topSequence + 1,
	})
	if err == nil {
		sysMsgData := map[string]interface{}{
			"chat_id":      chatId,
			"is_top":       isTop,
			"top_sequence": topSequence + 1,
			"device_id":    ctx.Value(base.HeaderFieldDeviceId),
		}
		SendSyncChatSettingMsg(ctx, forward.CmdSyncSettingIsTop, uid, sysMsgData)
	}
	return true
}

func UpdateIsTopV2(ctx context.Context, uid string, request IsTopRequest) (map[string]interface{}, error) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "UpdateIsTopV2"})
	chatId, isTop := request.ChatId, request.IsTop
	dao := setting.New()
	settingInfo, err := dao.GetSetting(uid, chatId)
	id := setting.GetId(uid, chatId)
	topSequence := 0
	res := make(map[string]interface{}, 0)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		_, er := addDefaultSetting(ctx, uid, chatId)
		if er != nil {
			log.Logger().Error(logCtx, "failed to update ,err : ", er)
			return res, errno.Add("add err:", errno.DefErr)
		}
	}
	topSequence = int(settingInfo.TopSequence)
	topTime := 0
	if isTop == 1 {
		topTime = int(funcs.GetMillis())
	}

	err = dao.UpdateByMap(id, map[string]interface{}{
		"is_top":       isTop,
		"top_time":     topTime,
		"top_sequence": topSequence + 1,
	})
	if err == nil {
		res = map[string]interface{}{
			"chat_id":      chatId,
			"is_top":       isTop,
			"top_sequence": topSequence + 1,
			"device_id":    ctx.Value(base.HeaderFieldDeviceId),
			"top_time":     topTime,
		}
		SendSyncChatSettingMsg(ctx, forward.CmdSyncSettingIsTop, uid, res)
	}
	return res, nil
}

func SendSyncChatSettingMsg(ctx context.Context, cmd, uid string, msgData map[string]interface{}) bool {
	contentByte, _ := json.Marshal(msgData)
	pushParams := forward.PushMessageParams{
		Cmd: cmd,
		Uid: uid,
		Data: forward.MsgItem{
			ID:        funcs.CreateMsgId(uid),
			Type:      base.MsgTypeCmd,
			Content:   string(contentByte),
			CreatedAt: funcs.GetMillis(),
		},
		NoPushOffline: true,
	}
	if err := forward.PushMessageToUserSocketDirectly(ctx, pushParams, false); err != nil {
		return false
	}
	return true
}

func HideChat(ctx context.Context, uid string, request HideRequest) bool {
	dao := setting.New()
	data, err := dao.GetSetting(uid, request.ChatId)
	id := setting.GetId(uid, request.ChatId)
	//logCtx1 := log.WithFields(context.Background(), map[string]string{"action": "DeleteChat"})
	//log.Logger().Info(logCtx1, "--4: ", uid, request)
	logCtx := log.WithFields(ctx, map[string]string{"action": "HideChat"})
	if err == mongo.ErrNoDocuments { // setting not exist, add default setting
		data = setting.Setting{
			ID:           id,
			UID:          uid,
			ChatID:       request.ChatId,
			HideSequence: request.HideSequence,
			IsTop:        setting.DefSetting.IsTop,
			IsMuteNotify: setting.DefSetting.IsMuteNotify,
			IsShowName:   setting.DefSetting.IsShowName,
			MuteTime:     0,
			Background:   setting.DefSetting.Background,
		}
		if _, err = dao.Add(data); err == nil {
			er := SendSyncHidePushMsg(ctx, uid, request.ChatId, request.HideSequence)
			if er != nil {
				log.Logger().Error(logCtx, "SendSyncHidePushMsg fail", er)
			}
			return true
		} else {
			log.Logger().Error(logCtx, "add configs fail", err)
		}
	} else if err == nil && data.ID != "" {
		if data.HideSequence > request.HideSequence {
			return true
		}
		uData := map[string]interface{}{
			"hide_sequence": request.HideSequence,
		}
		err = dao.UpdateByMap(id, uData)
		if err == nil {
			er := SendSyncHidePushMsg(ctx, uid, request.ChatId, request.HideSequence)
			if er != nil {
				log.Logger().Error(logCtx, "SendSyncHidePushMsg fail", er)
			}
			return true
		}
		return false
	}
	return false
}

func UpdateIsMuteNotify(ctx context.Context, uid, chatId string, isMuteNotify uint8) bool {
	logCtx := log.WithFields(ctx, map[string]string{"action": "UpdateIsMuteNotify"})
	dao := setting.New()
	data, err := dao.GetSetting(uid, chatId)
	id := setting.GetId(uid, chatId)
	muteSequence := data.MuteSequence + 1
	sysMsgData := map[string]interface{}{
		"chat_ids":      []string{chatId},
		"device_id":     ctx.Value(base.HeaderFieldDeviceId),
		"mute_sequence": muteSequence,
	}
	if errors.Is(err, mongo.ErrNoDocuments) { // setting not exist, add default setting
		data = setting.Setting{
			ID:           id,
			UID:          uid,
			ChatID:       chatId,
			IsMuteNotify: isMuteNotify,
			MuteTime:     0,
			IsTop:        setting.DefSetting.IsTop,
			IsShowName:   setting.DefSetting.IsShowName,
			Background:   setting.DefSetting.Background,
			HideSequence: setting.DefSetting.HideSequence,
			MuteSequence: 0,
		}
		if data.IsMuteNotify == 1 {
			data.MuteTime = funcs.GetMillis()
		}
		if _, err = dao.Add(data); err == nil {
			err = SendSyncMutePushMsg(ctx, uid, sysMsgData)
			if err != nil {
				return false
			}
			return true
		} else {
			log.Logger().Error(logCtx, "add configs fail", err)
		}
	} else if err == nil && data.ID != "" {
		muteTime := 0
		if isMuteNotify == 1 {
			muteTime = int(funcs.GetMillis())
		}

		uData := map[string]interface{}{
			"is_mute_notify": isMuteNotify,
			"mute_time":      muteTime,
			"mute_sequence":  muteSequence,
		}
		if dao.UpdateByMap(id, uData) == nil {
			SendSyncMutePushMsg(ctx, uid, sysMsgData)
			return true
		}
		return false
	}
	return false
}

func SendSyncHidePushMsg(ctx context.Context, uid, chatId string, sequence int64) error {
	sysMsgData := map[string]interface{}{
		"chat_id":  chatId,
		"sequence": sequence,
	}
	contentByte, _ := json.Marshal(sysMsgData)
	pushParams := forward.PushMessageParams{
		Cmd: forward.CmdSyncHideSetting,
		Uid: uid,
		Data: forward.MsgItem{
			ID:        funcs.CreateMsgId(uid),
			Type:      base.MsgTypeCmd,
			Content:   string(contentByte),
			CreatedAt: funcs.GetMillis(),
		},
		NoPushOffline: true,
	}
	return forward.PushMessageToUserSocketDirectly(ctx, pushParams, false)
}

func UpdateBackground(ctx context.Context, uid, chatId string, background map[string]interface{}) bool {
	logCtx := log.WithFields(ctx, map[string]string{"action": "UpdateBackground"})
	dao := setting.New()
	data, err := dao.GetSetting(uid, chatId)
	id := setting.GetId(uid, chatId)
	if err == mongo.ErrNoDocuments { // setting not exist, add default setting
		data = setting.Setting{
			ID:           id,
			UID:          uid,
			ChatID:       chatId,
			IsTop:        setting.DefSetting.IsTop,
			IsMuteNotify: setting.DefSetting.IsMuteNotify,
			IsShowName:   setting.DefSetting.IsShowName,
			HideSequence: setting.DefSetting.HideSequence,
			MuteTime:     0,
			Background:   background,
		}
		if data.IsMuteNotify == 1 {
			data.MuteTime = funcs.GetMillis()
		}
		if _, err := dao.Add(data); err == nil {
			return true
		} else {
			log.Logger().Error(logCtx, "add configs fail", err)
		}
	} else if err == nil && data.ID != "" {
		uData := setting.Setting{
			Background: background,
		}
		return dao.Update(id, uData)
	}
	return false
}

func UpdateIsShowName(ctx context.Context, uid, chatId string, isShowName uint8) bool {
	dao := setting.New()
	data, err := dao.GetSetting(uid, chatId)
	logCtx := log.WithFields(ctx, map[string]string{"action": "UpdateIsShowName"})
	log.Logger().Info(logCtx, data, err)
	id := setting.GetId(uid, chatId)
	if err == mongo.ErrNoDocuments { // setting not exist, add default setting
		data = setting.Setting{
			ID:           id,
			UID:          uid,
			ChatID:       chatId,
			IsTop:        setting.DefSetting.IsTop,
			IsMuteNotify: setting.DefSetting.IsMuteNotify,
			HideSequence: setting.DefSetting.HideSequence,
			IsShowName:   isShowName,
			MuteTime:     0,
			Background:   setting.DefSetting.Background,
		}
		if data.IsMuteNotify == 1 {
			data.MuteTime = funcs.GetMillis()
		}
		if _, err := dao.Add(data); err == nil {
			return true
		} else {
			log.Logger().Error(logCtx, "add configs fail", err)
		}
	} else if err == nil && data.ID != "" {
		return dao.UpdateValue(id, "is_show_name", isShowName)
	}
	return false
}

func addDefaultSetting(ctx context.Context, uid, chatId string) (bool, error) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "addDefaultSetting"})
	id := setting.GetId(uid, chatId)
	data := setting.DefSetting
	data = setting.Setting{
		ID:           id,
		UID:          uid,
		ChatID:       chatId,
		HideSequence: setting.DefSetting.HideSequence,
		MuteTime:     0,
		TopTime:      0,
	}
	_, err := setting.New().Add(data)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		log.Logger().Error(logCtx, "add default setting fail:", err)
		return false, err
	}
	return true, nil
}

func SendSyncMutePushMsg(ctx context.Context, uid string, sysMsgData map[string]interface{}) error {
	contentByte, _ := json.Marshal(sysMsgData)
	pushParams := forward.PushMessageParams{
		Cmd: forward.CmdSyncMuteSetting,
		Uid: uid,
		Data: forward.MsgItem{
			ID:        funcs.CreateMsgId(uid),
			Type:      base.MsgTypeCmd,
			Content:   string(contentByte),
			CreatedAt: funcs.GetMillis(),
		},
		NoPushOffline: true,
	}
	return forward.PushMessageToUserSocketDirectly(ctx, pushParams, false)
}
