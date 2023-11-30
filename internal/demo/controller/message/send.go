package message

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/demo/model/message"
	"imsdk/pkg/errno"
	"imsdk/pkg/log"
	"imsdk/pkg/response"
)

func SendCardAndTempMessage(ctx *gin.Context) {
	var params message.SendCardAndTempMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := message.SendCardAndTempMessage(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

// Post sendMessage
func SendTextMessage(ctx *gin.Context) {
	var params message.SendMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := message.SendTextMessage(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func SendImageMessage(ctx *gin.Context) {
	var params message.SendMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := message.SendImageMessage(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func SendAttachmentMessage(ctx *gin.Context) {
	var params message.SendMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := message.SendTextMessage(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func SendCardMessage(ctx *gin.Context) {
	var params message.SendMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := message.SendCardMessage(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func SendVerticalCardMessage(ctx *gin.Context) {
	var params message.SendMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := message.SendVerticalCardMessage(ctx, params)
	//err := message.SendVerticalCardMessageV2(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func SendMiddleMessage(ctx *gin.Context) {
	var params message.SendMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := message.SendMiddleMessage(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func SendCustomizeMessage(ctx *gin.Context) {
	var params message.SendMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	ctxLog := log.WithFields(ctx, map[string]string{"action": "SendCustomizeMessage"})
	log.Logger().Info(ctxLog, params, ctx.Request.Header)
	err := message.SendCustomizeMessage(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func SendNotificationMessage(ctx *gin.Context) {
	var params message.SendMessageParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := message.SendNotificationMessage(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
