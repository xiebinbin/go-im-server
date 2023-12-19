package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"imsdk/internal/client/model/user/login"
	"imsdk/internal/common/dao/user/qrcodelogin"
	"imsdk/internal/common/model/qrcode"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/response"
)

type PcLoginRes struct {
	ID           string      `json:"id"`
	Avatar       interface{} `json:"avatar"`
	PhonePrefix  string      `json:"phone_prefix"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Name         string      `json:"name"`
	Status       int         `json:"status"`
	Token        string      `json:"token"`
	SK           string      `json:"sk"`
	RefreshToken string      `json:"refresh_token"`
	ExpireIn     int64       `json:"expire_in"`
}

type CodeRequest struct {
	Code string `json:"code" binding:"required"`
}

func GenerateLoginCode(ctx *gin.Context) {
	data, err := qrcode.GeneratePcLoginQrCode(ctx)
	if err != nil {
		response.RespErr(ctx, err)
	} else {
		response.RespData(ctx, data)
	}
	return
}

func ScanQrCodeRes(ctx *gin.Context) {
	var params struct {
		Code []byte `json:"code" binding:"required"`
	}
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	//if err := ctx.ShouldBindJSON(&params); err != nil {
	//	response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
	//	return
	//}
	re, err := login.GetCodeResV2(params.Code)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	if re.ID == "" { // didn't scan qr code
		response.RespData(ctx, map[string]int{"status": qrcodelogin.StatusWaitScan})
		return
	}
	res := PcLoginRes{}
	if re.Status == qrcodelogin.StatusLoginConfirmed {
		_, er := login.ByScan(ctx, ctx.Value(base.HeaderFieldUID).(string))
		if er != nil {
			response.RespErr(ctx, er)
			return
		}
		res = PcLoginRes{
			ID:     re.Uid,
			Avatar: re.Avatar,
			Name:   re.Name,
			Status: qrcodelogin.StatusLoginConfirmed,
			SK:     re.Sk,
		}
		//rs := notify.SendCmdNotice(ctx, notify.SendCmdNoticeRequest{
		//	CMD: notify.CmdPCLoginIn,
		//	UID: data.Uid,
		//})
		////
		//if rs != nil {
		//	fmt.Println("notice err:", err)
		//}
	} else if re.Status == qrcodelogin.StatusScan {
		res = PcLoginRes{
			ID:     re.Uid,
			Avatar: re.Avatar,
			Name:   re.Name,
			Status: qrcodelogin.StatusScan,
		}
	}

	response.RespData(ctx, res)
	return
}

func AppScanLoginQrCode(ctx *gin.Context) {
	uid, _ := ctx.Get("uid")
	var params login.IdRequest
	//if err := ctx.ShouldBindJSON(&params); err != nil {
	//	response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
	//	return
	//}
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	fmt.Println("AppScanLoginQrCode******")
	if err := login.AppScanQrCode(ctx, uid.(string), params.Id); err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func AppConfirmLogin(ctx *gin.Context) {
	uid := ctx.Value(base.HeaderFieldUID)
	var params login.IdRequest
	//if err := ctx.ShouldBindJSON(&params); err != nil {
	//	response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
	//	return
	//}
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	if err := login.ConfirmLoginV2(ctx, uid.(string), params); err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
