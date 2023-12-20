package message

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/model/message"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/log"
	"imsdk/pkg/response"
	"time"
)

func SendMessage(ctx *gin.Context) {
	// auto build test
	logCtx := log.WithFields(ctx, map[string]string{"action": "SendMessage"})
	t1 := time.Now()
	var params message.SendMessageParams
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	uid, _ := ctx.Get(base.HeaderFieldUID)
	params.SenderId = uid.(string)
	resp, err := message.Send(ctx, params)
	log.Logger().Info(logCtx, "send message duration: time :  ", time.Now().Sub(t1))
	if err != nil {
		log.Logger().Error(logCtx, "SendMessage error -err : ", resp, err)
		response.RespErr(ctx, err)
		return
	}
	response.RespData(ctx, resp)
	return
}

func HollowManSendMessage(ctx *gin.Context) {
	// auto build test
	// do thing
	os := ctx.GetHeader("os")
	data := message.SendMessageParams{
		Type: base.MsgContentTypeHollow,
	}
	if err := message.HollowManSendMsg(ctx, data, os); err == nil {
		response.RespSuc(ctx)
	} else {
		response.RespErr(ctx, err)
	}
	return
}
