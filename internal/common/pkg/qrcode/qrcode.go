package qrcode

import (
	"bytes"
	"context"
	"errors"
	"imsdk/pkg/app"
	"imsdk/pkg/errno"
	"strings"
)

const (
	CodeInvalid      = errno.QrCodeInvalid
	PrimeNumber uint = 251
	VersionOne       = 1
)

var ModeMap = map[int]uint{
	0: 1, 1: 2, 2: 4, 3: 8, 4: 16, 5: 32, 6: 64, 7: 128, 8: 5, 9: 10,
	10: 20, 11: 40, 12: 80, 13: 160, 14: 69, 15: 138, 16: 25, 17: 50, 18: 100, 19: 200,
	20: 149, 21: 47, 22: 94, 23: 188, 24: 125, 25: 250, 26: 249, 27: 247, 28: 243, 29: 235,
	30: 219, 31: 187, 32: 123, 33: 246, 34: 241, 35: 231, 36: 211, 37: 171, 38: 91, 39: 182,
	40: 113, 41: 226, 42: 201, 43: 151, 44: 51, 45: 102, 46: 204, 47: 157, 48: 63, 49: 126,
	50: 1, 51: 2, 52: 4, 53: 8, 54: 16, 55: 32, 56: 64, 57: 128, 58: 5, 59: 10,
	60: 20, 61: 40, 62: 80, 63: 160, 64: 69, 65: 138, 66: 25, 67: 50, 68: 100, 69: 200,
	70: 149, 71: 47, 72: 94, 73: 188, 74: 125, 75: 250, 76: 249, 77: 247, 78: 243, 79: 235,
	80: 219, 81: 187, 82: 123, 83: 246, 84: 241, 85: 231, 86: 211, 87: 171, 88: 91, 89: 182,
	90: 113, 91: 226, 92: 201, 93: 151, 94: 51, 95: 102, 96: 204, 97: 157, 98: 63, 99: 126,
}

const (
	TypePcLogin  = "pl"
	TypeUserInfo = "ui"
	TypeGroup    = "g"
	TypeRecMoney = "rm"
)

type Resp struct {
	UrlPre     string `json:"url_pre"`
	Value      []uint `json:"value"`
	Expire     int64  `json:"expire"`
	CreateTime int64  `json:"create_time"`
}

var (
	Host = ""
)

func GetUrlPre(ctx context.Context, qrType string) string {
	if Host == "" {
		host, _ := app.Config().GetChildConf("global", "hosts", "qrcode_host")
		Host = host.(string)
	}
	return strings.TrimRight(Host, "/") + "/" + qrType + "/"
}

func Generate(code []byte) ([]uint, error) {
	length := len(code) + 2
	if length%3 != 0 {
		return nil, errors.New("code format is wrong")
	}
	bytesArr := [][]byte{{VersionOne}, code}
	codeSlice := bytes.Join(bytesArr, []byte(""))
	res := make([]uint, length)
	var p uint = 0
	for k, v := range codeSlice {
		vv := uint(v)
		res[k] = vv
		p = (p + vv*ModeMap[k]) % PrimeNumber
	}
	res[length-1] = p
	return res, nil
}

func Verify(code []byte) ([]byte, error) {
	ver := code[0]
	if ver != VersionOne {
		return nil, errno.Add("version is not right", errno.ParamsErr)
	}
	var p uint = 0
	length := len(code)
	for k, v := range code[0 : length-1] {
		vv := uint(v)
		p = (p + vv*ModeMap[k]) % PrimeNumber
	}
	if code[length-1] != uint8(p) {
		return nil, errno.Add("code is invalid ", CodeInvalid)
	}
	return code[1 : length-1], nil
}
