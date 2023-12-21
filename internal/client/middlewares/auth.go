package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/client/pkg/utils/crypt"
	walletutil "imsdk/internal/client/pkg/utils/wallet-util"
	"imsdk/internal/common/dao/user"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/config"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
	"net/http"
)

type VerifyParams struct {
	Data string `json:"data"`
}

func Auth(ctx *gin.Context) {
	dataHash := ctx.Request.Header.Get(base.HeaderFieldDataHash)
	if ctx.Request.Method == http.MethodGet {
		//计算共享密钥
	}
	var data string
	if params, er := getCommonParams(ctx); er != nil {
		ctx.Abort()
		response.ResPubErr(ctx, errno.Add("body-params-err", http.StatusBadRequest))
		return
	} else {
		data = params
	}
	time := ctx.Value(base.HeaderFieldTime).(string)
	sign := ctx.Value(base.HeaderFieldSign).(string)
	address := ctx.Value(base.HeaderFieldUID).(string)
	fmt.Println("---", time, sign, dataHash)
	if ctx.Value(base.HeaderIsEnc).(string) != "false" {
		rel, err := walletutil.VerifySign(dataHash+":"+time, sign, address)
		if rel == false || err != nil {
			fmt.Println("sign-auth-err", rel, err)
			ctx.Abort()
			response.ResPubErr(ctx, errno.Add("sign-auth-err", http.StatusBadRequest))
			return
		}
	}
	pubKey := ctx.Value(base.HeaderFieldPubKey).(string)
	deCrypto(ctx, data, pubKey)
	uInfo, err := user.New().GetByID(address)
	fmt.Println("err:", err, address)
	if uInfo.ID == "" || errors.Is(err, mongo.ErrNoDocuments) {
		ctx.Abort()
		response.ResPubErr(ctx, errno.Add("request-err-x-uid", http.StatusBadRequest))
		return
	}
	ctx.Set("uid", uInfo.ID)
}

func deCrypto(ctx *gin.Context, data, pubKey string) {
	priKey, _ := config.GetConfigSk()
	client, _ := walletutil.New(priKey)
	key, _ := client.GetSharedSecret(pubKey)
	//-------- DEMO Start ---------
	//dataDemo := map[string]interface{}{
	//	"chat_id":    "7f020ec3ee54fe6a79225e1f6bd29bba",
	//	"content":    "加密的hello world",
	//	"to_address": "0x0b84b2d122cb1c058b988d9f0291a6e25364c6f8d",
	//}
	//dataByte, _ := json.Marshal(dataDemo)
	//
	//data, _ = crypt.En(key, dataByte)
	//========= DEMO End =========
	fmt.Println("ctx.Value(base.HeaderIsEnc)", ctx.Value(base.HeaderIsEnc))
	if ctx.Value(base.HeaderIsEnc).(string) == "false" {
		ctx.Set("data", data)
		return
	}

	res, _ := crypt.De(key, data)
	fmt.Println("req data:", string(res))
	ctx.Set("data", string(res))
	/*res, _ := json.Marshal(data)*/
}

func getCommonParams(ctx *gin.Context) (string, error) {
	var params VerifyParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		return params.Data, errno.Add("params-err", errno.ParamsErr)
	}
	return params.Data, nil
}

func GetUId(ctx *gin.Context) string {
	uid, exists := ctx.Get("uid")
	if !exists || uid == nil {
		return ""
	}
	return uid.(string)
}
