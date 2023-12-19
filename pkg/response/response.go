package response

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/pkg/utils/crypt"
	walletutil "imsdk/internal/client/pkg/utils/wallet-util"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/config"

	//"imsdk/internal/bucket/model/copywriting"
	"imsdk/pkg/encrypt"
	"imsdk/pkg/errno"
	"net/http"
)

// Response HTTP返回数据结构体, 可使用这个, 也可以自定义
type Response struct {
	Code int         `json:"code"` // 状态码,这个状态码是与前端和APP约定的状态码,非HTTP状态码
	Data interface{} `json:"data"` // 返回数据
	Msg  string      `json:"msg"`  // 自定义返回的消息内容
}

func RespErr(ctx *gin.Context, err error) {
	var errInfo errno.Errno
	errStr := err.Error()
	err = json.Unmarshal([]byte(errStr), &errInfo)
	var resp Response
	enData := map[string]interface{}{
		"err_code": errInfo.Code,
		"err_msg":  errInfo.Msg,
	}
	fmt.Println("RespErr----", enData)
	byteData, _ := json.Marshal(enData)
	pubKey := ctx.Value(base.HeaderFieldPubKey).(string)
	priKey, _ := config.GetConfigSk()
	client, _ := walletutil.New(priKey)
	key, _ := client.GetSharedSecret(pubKey)

	res, _ := crypt.En(key, byteData)
	resp = Response{
		Code: errno.OK,
		Msg:  errInfo.Msg,
		Data: res,
	}
	response(ctx, resp)
	return
}

func RespSuc(ctx *gin.Context) {
	resp := Response{
		Code: errno.OK,
		Msg:  "",
		Data: nil,
	}
	response(ctx, resp)
	return
}

func ResPubErr(ctx *gin.Context, err error) {
	var errInfo errno.Errno
	errStr := err.Error()
	err = json.Unmarshal([]byte(errStr), &errInfo)
	resp := Response{
		Code: errInfo.Code,
		Msg:  errInfo.Msg,
		Data: nil,
	}
	response(ctx, resp)
	return
}

func ResData(ctx *gin.Context, data interface{}) {
	resp := Response{
		Code: errno.OK,
		Msg:  "",
		Data: map[string]interface{}{
			"items":    data,
			"err_code": 0,
		},
	}
	response(ctx, resp)
	return
}

func RespListData(ctx *gin.Context, data interface{}) {
	enData := map[string]interface{}{
		"items": data,
	}
	a, _ := json.Marshal(enData)
	fmt.Println("RespListData----", string(a))
	byteData, _ := json.Marshal(enData)
	pubKey := ctx.Value(base.HeaderFieldPubKey).(string)
	priKey, _ := config.GetConfigSk()
	client, _ := walletutil.New(priKey)
	key, _ := client.GetSharedSecret(pubKey)
	res, _ := crypt.En(key, byteData)
	fmt.Println("res:", res)
	resp := Response{
		Code: errno.OK,
		Msg:  "",
		Data: res,
	}
	response(ctx, resp)
	return
}

func RespDataWithNoEnc(ctx *gin.Context, data interface{}) {
	resp := Response{
		Code: errno.OK,
		Msg:  "",
		Data: data,
	}
	response(ctx, resp)
	return
}

func RespData(ctx *gin.Context, data interface{}) {
	fmt.Println("RespData---", data)
	byteData, _ := json.Marshal(data)
	pubKey := ctx.Value("X-Pub-Key").(string)
	priKey, _ := config.GetConfigSk()
	client, _ := walletutil.New(priKey)
	key, _ := client.GetSharedSecret(pubKey)
	res, _ := crypt.En(key, byteData)
	resp := Response{
		Code: errno.OK,
		Msg:  "",
		Data: res,
	}
	response(ctx, resp)
	return
}

func response(ctx *gin.Context, resp Response) {
	res, _ := json.Marshal(resp)
	_, isEnc := ctx.Get("is_enc")
	encKey, ok := ctx.Get("enc_key")
	if isEnc && ok {
		encResStr := encrypt.AesCbcEncrypt(res, encKey.(string))
		ctx.String(http.StatusOK, encResStr)
		return
	}
	ctx.String(http.StatusOK, string(res))
	return
}

func RespOriginData(ctx *gin.Context, data interface{}) {
	resp := Response{
		Code: errno.OK,
		Msg:  "",
		Data: data,
	}
	response(ctx, resp)
	return
}
