package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/pkg/utils/config"
	crypto "imsdk/internal/client/pkg/utils/crypt"
	walletutil "imsdk/internal/client/pkg/utils/wallet-util"
	"io/ioutil"
	"net/http"
)

func Ecc(ctx *gin.Context) {
	pubKey := ctx.Request.Header.Get("PubKey")
	if pubKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "参数不正确"})
		ctx.Abort()
		return
	}
	data, err := ctx.GetRawData()
	if ctx.Request.Method == http.MethodGet {
		data = []byte(ctx.Request.Header.Get("Data"))
	}
	if err == nil {
		sign := ctx.Request.Header.Get("Sign")
		address := ctx.Request.Header.Get("Address")
		//dataLength, err := convertor.ToInt(ctx.Request.Header.Get("DataLength"))
		dataLength := len(ctx.Request.Header.Get("DataLength"))
		if sign == "" || address == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "参数不正确"})
			ctx.Abort()
			return
		}
		rel, err := walletutil.VerifySign(string(data), sign, address)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
			ctx.Abort()
			return
		}
		if rel == false {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "签名验证失败"})
			ctx.Abort()
			return
		}

		serverWallet := config.GetServerWallet()
		sharedSecret, err := serverWallet.GetSharedSecret(pubKey)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
			ctx.Abort()
			return
		}
		hexEnData := string(data)
		if hexEnData[:2] != "0x" {
			hexEnData = "0x" + hexEnData
		}
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
			ctx.Abort()
			return
		}
		deData, err := crypto.De(sharedSecret, hexEnData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
			ctx.Abort()
			return
		}
		var deRel map[string]interface{}
		err = json.Unmarshal(deData[:dataLength], &deRel)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
			ctx.Abort()
			return
		}
		//signUnixTime, err := convertor.ToInt(deRel["time"])
		//if err != nil {
		//	ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
		//	ctx.Abort()
		//	return
		//}
		//now := datetime.AddMinute(time.Now(), -1)
		//if signUnixTime < now.Unix() {
		//	if err != nil {
		//		ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "签名过期"})
		//		ctx.Abort()
		//		return
		//	}
		//}
		//dataBytes, err := convertor.ToBytes(deRel["data"])
		dataBytes, err := json.Marshal(deRel["data"])
		//if err != nil {
		//	ctx.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
		//	ctx.Abort()
		//	return
		//}
		if ctx.Request.Method != http.MethodGet {
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(dataBytes))
		}
		ctx.Next()
	}
}
