package login

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao"
	user2 "imsdk/internal/common/dao/user"
	"imsdk/internal/common/dao/user/qrcodelogin"
	"imsdk/internal/common/pkg/qrcode"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"time"
)

const (
	CodeCacheTag      = "pc_scan_login_uuid:"
	CodeExpiredSecond = 120
	EmptyString       = ""
	StatusWait        = 0
	StatusScan        = 1
	StatusLogin       = 2

	ErrNoNotScan = 1
)

type IdRequest struct {
	Id   string `json:"id" binding:"required"`
	Code []byte `json:"code"`
	SK   string `json:"sk"`
}

type CodeInfo struct {
	Status int
	UID    string
}

func GetUniqueScanCode() (map[string]string, error) {
	res := make(map[string]string)
	codeInfo := map[string]interface{}{
		"status": StatusWait,
		"uid":    EmptyString,
	}
	code := funcs.UniqueId32()
	cacheTag := CodeCacheTag + code
	bytesData, _ := json.Marshal(codeInfo)
	err := redis.Client.Set(cacheTag, bytesData, time.Second*CodeExpiredSecond)
	if err != nil {
		res["code"] = code
		return res, nil
	}
	return res, errno.Add("fail", errno.DefErr)
}

func UpdateCodeStatus(code string, status int, uid string) error {
	codeInfo := map[string]interface{}{
		"status": status,
		"uid":    uid,
	}
	cacheTag := CodeCacheTag + code
	bytesData, _ := json.Marshal(codeInfo)
	_, err := GetCodeDetail(code)
	if err != nil {
		return err
	}
	err = redis.Client.Set(cacheTag, bytesData, time.Second*CodeExpiredSecond).Err()
	if err != nil {
		return errno.Add("code-expired", errno.Expired)
	}
	return nil
}

func GetCodeDetail(code string) (map[string]interface{}, error) {
	cacheTag := CodeCacheTag + code
	str := redis.Client.Get(cacheTag).Val()
	res := make(map[string]interface{}, 0)
	if str != "" {
		err := json.Unmarshal([]byte(str), &res)
		if err != nil {
			return res, errno.Add("fail", errno.DefErr)
		}
		return res, nil
	}
	return res, errno.Add("code-expired", errno.Expired)
}

func AppScanQrCode(ctx context.Context, uid, id string) error {
	// get code info
	qrInfo, _ := qrcodelogin.New().GetByID(id)
	if qrInfo.Uid != "" && qrInfo.Uid != uid {
		return errno.Add("code is used", qrcodelogin.ErrorUsed)
	}
	if qrInfo.Expire < funcs.GetMillis() {
		return errno.Add("code-expired", errno.Expired)
	}
	//get user info
	uInfo, err := user2.New().GetInfoById(uid)
	if err != nil {
		if dao.IsNoDocumentErr(err) {
			return errno.Add("user does not exist", errno.UserNotExist)
		} else {
			return errno.Add("sys error", errno.SysErr)
		}
	} else if uInfo.Status != 1 {
		return errno.Add("user is unusable", errno.UserUnavailable)
	}
	data := qrcodelogin.QrCodeLogin{
		ID:        id,
		Uid:       uid,
		Avatar:    uInfo.Avatar,
		Name:      uInfo.Name,
		Expire:    funcs.GetMillis() + 60000,
		Status:    qrcodelogin.StatusScan,
		UpdatedAt: funcs.GetMillis(),
	}
	if err = qrcodelogin.New().Save(data); err != nil && !mongo.IsDuplicateKeyError(err) {
		return errno.Add("fail", errno.DefErr)
	}
	return nil
}

func AppScanQrCodeV2(ctx context.Context, uid string, code []byte) error {
	if codeByte, err := qrcode.Verify(code); err != nil {
		return err
	} else {
		return AppScanQrCode(ctx, uid, string(codeByte))
	}
}

func ConfirmLogin(ctx context.Context, uid string, request IdRequest) error {
	codeLoginDao := qrcodelogin.New()
	loginInfo, err := codeLoginDao.GetByID(request.Id)
	if err != nil {
		return errno.Add("please scan first", ErrNoNotScan)
	} else if loginInfo.Uid != uid {
		return errno.Add("forbidden", errno.FORBIDDEN)
	} else if loginInfo.Expire < funcs.GetMillis() {
		return errno.Add("expired", errno.Expired)
	}
	uData := qrcodelogin.QrCodeLogin{
		Status: qrcodelogin.StatusLoginConfirmed,
		Sk:     request.SK,
	}
	if err = codeLoginDao.UpByID(request.Id, uData); err != nil {
		return errno.Add("fail", errno.DefErr)
	}
	return nil
}

func ConfirmLoginV2(ctx context.Context, uid string, request IdRequest) error {
	codeByte, er := qrcode.Verify(request.Code)
	if er != nil {
		return er
	}
	return ConfirmLogin(ctx, uid, IdRequest{
		Id: string(codeByte),
		SK: request.SK,
	})
}

func GetCodeRes(id string) (qrcodelogin.QrCodeLogin, error) {
	var emptyCodeRes qrcodelogin.QrCodeLogin
	data, err := qrcodelogin.New().GetByID(id)
	if err != nil {
		return emptyCodeRes, nil
	} else if data.Expire < funcs.GetMillis() {
		return emptyCodeRes, errno.Add("expired", errno.Expired)
	}
	return data, nil
}

func GetCodeResV2(code []byte) (qrcodelogin.QrCodeLogin, error) {
	codeByte, er := qrcode.Verify(code)
	if er != nil {
		return qrcodelogin.QrCodeLogin{}, er
	}
	return GetCodeRes(string(codeByte))
}
