package singapore

import (
	"context"
	"encoding/json"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"strconv"
)

type OfflineReqParams struct {
	Title      interface{} `json:"title"`
	Body       interface{} `json:"body"`
	Extra      interface{} `json:"extra"`
	Badge      int64       `json:"badge"`
	Devices    []string    `json:"devices"`
	AMId       string      `json:"amid"`
	SenderId   string      `json:"senderId"`
	ReceiveIds []string    `json:"receive_ids"`
	AChatId    string      `json:"achat_id"`
	Type       int8        `json:"type"`
	CollapseID string      `json:"collapse_id"`
}

type KeyValMsgItem struct {
	Type string      `json:"t"`
	Val  interface{} `json:"v"`
}

type OfflineReqExtra struct {
	Cmd   string      `json:"cmd"`
	Items interface{} `json:"items"`
}

type OfflineRequest struct {
	Title       string `json:"tile"`
	Body        string `json:"body"`
	CollapseID  string `json:"collapse_id"`
	ClickAction string `json:"click_action"`
}

func Offline(ctx context.Context, callbackUrl string, params OfflineReqParams) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "ProcessOfflineMsg"})
	singaporeTitle, _ := json.Marshal(params.Title)
	singaporeBody, _ := json.Marshal(params.Body)
	singaporeReceiveIds, _ := json.Marshal(params.ReceiveIds)
	log.Logger().Info(logCtx, map[string]interface{}{"params.Devices": params.Devices})
	var deviceIds []string
	for _, device := range params.Devices {
		if device == base.WebDeviceId || device == base.OsWeb {
			continue
		}
		deviceIds = append(deviceIds, device)
	}
	if len(deviceIds) == 0 {
		return
	}
	callbackParams := CallBackOfflineParams{
		AMID:       params.AMId,
		Type:       strconv.Itoa(int(params.Type)),
		AChatId:    params.AChatId,
		Title:      string(singaporeTitle),
		Body:       string(singaporeBody),
		SenderId:   params.SenderId,
		ReceiveIds: string(singaporeReceiveIds),
		Extra:      params.Extra,
		DeviceIds:  deviceIds,
	}
	msg, _ := json.Marshal(callbackParams)
	log.Logger().Info(logCtx, map[string]interface{}{
		"msg":       string(msg),
		"startTime": funcs.GetMillis(),
		"amid":      params.AMId,
	})
	formData := make(map[string]interface{}, 0)
	err := json.Unmarshal(msg, &formData)
	if err != nil {
		log.Logger().Error(logCtx, map[string]interface{}{
			"err": err,
		})
		return
	}
	headers := map[string]string{
		"Api-Token": GetToken(),
	}
	res, er := PostForm(callbackUrl, formData, headers)
	if er != nil {
		log.Logger().Error(logCtx, map[string]interface{}{
			"res": res,
			"err": er,
		})
		return
	}
	// var resData CallbackResp
	var resData map[string]interface{}
	err = json.Unmarshal(res, &resData)
	log.Logger().Info(logCtx, map[string]interface{}{
		"endTime":        funcs.GetMillis(),
		"amid":           params.AMId,
		"params.Devices": params.Devices,
	})
	return
}
