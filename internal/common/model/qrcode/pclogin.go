package qrcode

import (
	"context"
	"imsdk/internal/common/dao/user/qrcodelogin"
	"imsdk/internal/common/pkg/ip"
	"imsdk/internal/common/pkg/qrcode"
	"imsdk/pkg/app"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/unique"
	"strconv"
	"strings"
	"time"
)

const (
	TargetOpenWebsite    = "h5"
	TargetUserProfile    = "ui"
	TargetJoinGroup      = "g"
	TargetMomentDetail   = "mi"
	TargetPcLogin        = "pl"
	TargetReceiveMoney   = "rm"
	TargetUserTempQrcode = "tq"
)

type BaseResp struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	Expire     int64  `json:"expire"`
	CreateTime int64  `json:"create_time"`
}

func formatQrCodeUrl(qrType, indexId string) string {
	host, _ := app.Config().GetChildConf("global", "hosts", "qrcode_host")
	return strings.TrimRight(host.(string), "/") + "/qr/" + qrType + "/" + indexId
}

func GeneratePcLoginQrCode(ctx context.Context) (qrcode.Resp, error) {
	// timeout's unit : second
	timeout, err := app.Config().GetChildConf("global", "system", "pc_login_qrcode_timeout")
	if err != nil {
		return qrcode.Resp{}, errno.Add("lack of config", errno.Exception)
	}
	t := funcs.GetMillis()
	id, ipStr := funcs.Md516(unique.Id12()), ctx.Value("ip").(string)
	expire := t + int64(timeout.(float64))*1000
	version := ctx.Value("version").(string)
	ipInfo, _ := ip.GetIpInfo(ipStr)
	//deviceName := ctx.Value(base.FieldDeviceName).(string)
	//bindInfo := terminal.GetBandInfo(deviceName)
	date := time.Now().Format("20060102")
	dateInt, _ := strconv.Atoi(date)
	osType := ctx.Value("os").(string)
	data := qrcodelogin.QrCodeLogin{
		ID:        id,
		Expire:    expire + 60*1000,
		Status:    qrcodelogin.StatusWaitScan,
		Country:   ipInfo.Country,
		Region:    ipInfo.RegionName,
		City:      ipInfo.City,
		Timezone:  ipInfo.Timezone,
		Date:      dateInt,
		DateIdx:   date[4:] + date[0:4],
		Version:   version,
		Os:        osType,
		Ip:        ipStr,
		CreateAt:  t,
		UpdatedAt: t,
	}
	if qrcodelogin.New().Save(data) != nil {
		return qrcode.Resp{}, errno.Add("fail", errno.DefErr)
	}

	code, err := qrcode.Generate([]byte(data.ID))
	if err != nil {
		return qrcode.Resp{}, errno.AddSysErr("fail:2", errno.SysErr)
	}

	return qrcode.Resp{
		UrlPre:     qrcode.GetUrlPre(ctx, qrcode.TypePcLogin),
		Value:      code,
		Expire:     expire,
		CreateTime: t,
	}, nil
}
