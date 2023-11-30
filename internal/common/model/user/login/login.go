package login

import (
	//"encoding/json"
	//"imsdk/internal/common/dao/common/devicetoken"
)

type ByLoginRequest struct {
	Prefix  string `json:"prefix" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Captcha string `json:"captcha" binding:"required"`
	Ip      string `json:"ip"`
}

type SucResponse struct {
	ID           string `json:"id"`
	PhonePrefix  string `json:"phone_prefix"`
	Token        string `json:"token"`
	Code         int8   `json:"code"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Language     string `json:"language"`
}

type RegManyRequest struct {
	Prefix string `json:"prefix" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
}

const (
	UserCodeNormal    = 0
	UserCodeForbidden = 1
	UserCodeMute      = 2

	StatusRegisterFailed = 100012
)

//func offline(uid, os, newDevice string, receiverTers []string, closeTag int8) bool {
//	msgContent := map[string]interface{}{
//		"device_name": newDevice,
//	}
//	bytesData, _ := json.Marshal(msgContent)
//	data := notify.CloseSocketConnectionParams{
//		Os:          os,
//		CloseTag:    closeTag,
//		MsgContent:  bytesData,
//		ReceiveTers: receiverTers,
//	}
//	devicetoken.EmptyTokenInfo(uid)
//	notify.CloseSocketConnection(uid, data)
//	return true
//}
