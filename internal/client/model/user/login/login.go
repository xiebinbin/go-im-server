package login

import (
	"context"
	"fmt"
	"imsdk/internal/common/dao"
	"imsdk/internal/common/dao/user"
	"imsdk/internal/common/dao/user/terminalinfo"
	"imsdk/internal/common/model/user/token"

	//"imsdk/internal/common/model/notify"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"math/rand"
	"strconv"
	"time"
)

type ByLoginRequest struct {
	Prefix  string `json:"prefix" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Captcha string `json:"captcha"`
	Env     string `json:"env"`
	Ip      string `json:"ip"`
}

type ByPasswordRequest struct {
	Prefix   string `json:"prefix" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password"`
	Env      string `json:"env"`
	Ip       string `json:"ip"`
}

type ByThirdPartRequest struct {
	Type   int    `json:"type" binding:"required"`
	Token  string `json:"token" binding:"required"`
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Env    string `json:"env"`
	Ip     string `json:"ip"`
}

type SucResponse struct {
	ID           string `json:"id"`
	Token        string `json:"token"`
	Code         int8   `json:"code"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Language     string `json:"language"`
	IsNew        int8   `json:"is_new"`
}

type RegManyRequest struct {
	Prefix   string `json:"prefix" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password"`
}

type Avatar struct {
	BucketId string `json:"bucketId"`
	ObjectId string `json:"objectId"`
	Text     string `json:"text"`
	FileType string `json:"file_type"`
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	Size     int64  `json:"size"`
	IsOrigin int8   `json:"isOrigin"`
}

const (
	UserCodeNormal       = 0
	UserCodeForbidden    = 1
	UserCodeMute         = 2
	ThirdPartGoogle      = 1
	ThirdPartApple       = 2
	TypeCaptcha          = 0
	TypePassword         = 1
	StatusRegisterFailed = 100012
	ErrUnSupportLoginWay = 300001
	AvatarBucketId       = "7e4216ac1941c61c"
)

func login(ctx context.Context, uId string, isNew int8) (SucResponse, error) {
	var emptyUserInfo SucResponse
	terminalSource := terminalinfo.SourceLogin
	if isNew == 1 {
		terminalSource = terminalinfo.SourceReg
	} else {
		uInfo, err := user.New().GetInfoById(uId)
		if err != nil && !dao.IsNoDocumentErr(err) {
			return emptyUserInfo, errno.Add("sys error", errno.SysErr)
		}
		if uInfo.ID != "" {
			return emptyUserInfo, errno.Add("user-unavailable", errno.UserUnavailable)
		}
	}
	deviceId := ""
	getToken, err := token.GetToken(ctx, token.GetTokenParams{
		Uid:      uId,
		DeviceId: deviceId,
	})
	if err != nil {
		return emptyUserInfo, err
	}
	data := SucResponse{
		ID:           uId,
		Code:         UserCodeNormal,
		Token:        getToken.Token,
		RefreshToken: getToken.RefreshToken,
		ExpiresIn:    getToken.Expire,
		Language:     "",
		IsNew:        isNew,
	}
	logData := LogRequest{
		UId:         uId,
		DeviceName:  "",
		ReceiveTers: []string{base.ClientTypePC, base.ClientTypeApp},
	}
	AddLoginLog(logData)
	fmt.Println("terminalSource:", terminalSource)
	//terminalData := terminal.AddTerminalInfoRequest{
	//	UId:    uId,
	//	Ip:     ctx.Value("ip").(string),
	//	Source: terminalSource,
	//}
	//terminal.AddTerminalInfo(ctx, terminalData)
	return data, nil
}

func ByScan(ctx context.Context, uid string) (SucResponse, error) {
	var emptyUserInfo SucResponse
	// CHeck if the user exists
	userDao := user.New()
	uInfo, err := userDao.GetByID(uid)
	if err != nil && !dao.IsNoDocumentErr(err) {
		return emptyUserInfo, errno.Add("sys-err", errno.DefErr)
	}
	logData := LogRequest{
		UId: uid,
	}
	AddLoginLog(logData)
	getToken, er := token.GetToken(ctx, token.GetTokenParams{
		Uid: uInfo.ID,
	})
	if er != nil {
		return emptyUserInfo, err
	}
	data := SucResponse{
		ID:           uInfo.ID,
		Code:         UserCodeNormal,
		Token:        getToken.Token,
		RefreshToken: getToken.RefreshToken,
		ExpiresIn:    getToken.Expire,
		Language:     "",
	}
	return data, nil
}

func createDefaultOldAvatar(randNum int) string {
	return "member/avatar/default/" + strconv.Itoa(randNum) + ".jpg"
}

func GetDefaultName() string {
	str := "0123456789"
	b := []byte(str)
	result := make([]byte, 6)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		result[i] = b[r.Intn(len(b))]
	}
	return "bbchat" + string(result)
}

func GetDefaultAvatar(randNum int) (avatar map[string]interface{}) {
	avatar = map[string]interface{}{
		"text":      "avatar/default/" + strconv.Itoa(randNum),
		"objectId":  "avatar/default/" + strconv.Itoa(randNum),
		"height":    240,
		"width":     240,
		"file_type": "png",
		"bucketId":  "local",
	}
	return
}

func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}
