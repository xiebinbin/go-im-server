package login

import (
	"imsdk/internal/common/dao/user/loginlog"
	"imsdk/pkg/funcs"
)

type LogRequest struct {
	UId         string   `json:"uid"`
	Os          string   `json:"os"`
	DeviceName  string   `json:"device_name"`
	DeviceId    string   `json:"device_id"`
	ReceiveTers []string `json:"receive_ters"`
}

func AddLoginLog(request LogRequest) {
	t := funcs.GetMillis()
	addData := loginlog.LoginLog{
		ID:         funcs.UniqueId16(),
		Os:         request.Os,
		UId:        request.UId,
		DeviceName: request.DeviceName,
		DeviceId:   request.DeviceId,
		CreatedAt:  t,
		UpdatedAt:  t,
	}
	_, err := loginlog.New().AddLogin(addData)
	if err != nil {
		return
	}
	//closeTag := notify.CloseTagAnotherDevice
	//if request.Os == base.OsAndroid || request.Os == base.OsIos {
	//	closeTag = notify.CloseTagAnotherAppDevice
	//	offline(request.UId, request.Os, request.DeviceName, request.ReceiveTers, int8(closeTag))
	//}
}
