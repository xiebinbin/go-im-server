package socket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"imsdk/internal/common/pkg/config"
	"imsdk/internal/sdkserver/model/message"
)

func SendCmd(ctx *gin.Context) error {
	params := map[string]interface{}{
		"cmd":       "_application_app",
		"auid":      "2x2qlr88wdcz",
		"target_id": "973ed426b316",
		"is_self":   1,
		"data": map[string]interface{}{
			"cmd":   "apply-add-friend",
			"items": "",
		},
	}
	dataByte, _ := json.Marshal(params)
	//_, err := model.RequestIMServer("sendCmd", string(dataByte))
	//if err != nil {
	//	return err
	//}
	fmt.Println("string(dataByte)----", string(dataByte))
	ak, _ := config.GetConfigAk()
	ctx.Set("ak", ak)
	ctx.Set("data", string(dataByte))
	err := message.SendCmd(ctx)
	if err != nil {
		return err
	}
	return nil
}
