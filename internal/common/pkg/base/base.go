package base

import (
	"imsdk/pkg/app"
	"imsdk/pkg/funcs"
	"strings"
)

type JsonString = string
type JsonMap = map[string]interface{}

const (
	HeaderFieldDeviceId  = "device-id"
	HeaderFieldReqId     = "req-id"
	HeaderFieldLang      = "lang"
	HeaderFieldOs        = "os"
	HeaderFieldVersion   = "version"
	HeaderFieldUserAgent = "user-agent"
	HeaderFieldPubKey    = "X-Pub-Key"
	HeaderFieldUID       = "X-UID"
	HeaderFieldSign      = "X-Sign"
	HeaderFieldTime      = "X-Time"
	HeaderFieldDataHash  = "X-Data-Hash"

	TypeUser            = "im:login:user:"
	AKSingapore         = "68oni7jrg31qcsaijtg76qln"
	TypeUserSDK         = "im:login:"
	TypeUserAdmin       = "im:login:user:admin:"
	RedisUserLatestConn = "im:conn:uid:"
	RedisUserLatestVer  = "im:user:latest:ver"
	TypeUserApp         = "im:login:user:app:"
	TypeUserPc          = "im:login:user:pc:"
	ClientTypePC        = "pc"
	ClientTypeAdmin     = "admin"
	ClientTypeApp       = "app"
	ClientTypeWeb       = "web"
	OsUnknown           = "unknown"
	OsAndroid           = "android"
	OsIos               = "ios"
	OsWindows           = "win"
	OsMac               = "mac"
	OsWeb               = "web"
	OsAdmin             = "admin"
	OsThirdServer       = "third-server"
	OsUnknownCode       = 0
	OsIosCode           = 1
	OsAndroidCode       = 2
	OsWindowsCode       = 3
	OsMacCode           = 4

	MsgNumberCountersKey = "246d8360c39d736b" // funcs.Md516("messagedetail")
	PageRowLimit         = 1000
	MsgPageRowLimit      = 10000
)
const (
	OfficialPaymentChatId = "g_pay"
	OfficialNoticeChatId  = "g_official"
)

func GetTerType(os string) string {
	os = strings.ToLower(os)
	res := ""
	switch os {
	case OsWindows:
		res = ClientTypePC
	case OsMac:
		res = ClientTypePC
	case OsAndroid:
		res = ClientTypeApp
	case OsIos:
		res = ClientTypeApp
	case OsAdmin:
		res = TypeUserAdmin
	}
	return res
}

func GetOsCode(osString string) int {
	var osCode int
	osStr := strings.ToLower(osString)
	switch osStr {
	case OsIos:
		osCode = OsIosCode
		break
	case OsAndroid:
		osCode = OsAndroidCode
		break
	case OsWindows:
		osCode = OsWindowsCode
		break
	case OsMac:
		osCode = OsMacCode
		break
	default:
		osCode = OsAndroidCode
		break
	}
	return osCode
}

func GetOs(os int) string {
	osString := "android"
	switch os {
	case OsIosCode:
		osString = OsIos
		break
	case OsAndroidCode:
		osString = OsAndroid
		break
	case OsWindowsCode:
		osString = OsWindows
		break
	case OsMacCode:
		osString = OsMac
		break
	default:
		osString = OsAndroid
		break
	}
	return osString
}

func CreateMessageId(amid, ak string) string {
	return funcs.Md5Str(amid + ak + "createMessageId")
}

func CreateUId(auid, ak string) string {
	if ak == app.OfficialAK {
		return auid
	}
	return funcs.Md5Str(auid + ak + "createUId")
}

func CreateChatId(achatId, ak string) string {
	return funcs.Md5Str(achatId + ak + "createChatId")
}
