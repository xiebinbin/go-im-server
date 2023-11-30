package message

import (
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/model/message"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/response"
)

type DelSelfMessageRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

type DelSelfByChatRequest struct {
	ChatIds    []string `json:"chat_ids" binding:"required"`
	ExceptMIds []string `json:"except_mids"`
}

func DelSelf(ctx *gin.Context) {
	var params DelSelfMessageRequest
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	defErr := errno.Add("params-err", errno.ParamsErr)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, defErr)
		return
	}
	if message.DeleteSelfByIds(uid, params.IDs) {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, defErr)
	return
}

func DelSelfByChatIds(ctx *gin.Context) {
	var params DelSelfByChatRequest
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	defErr := errno.Add("params-err", errno.ParamsErr)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, defErr)
		return
	}
	if message.DeleteSelfByChatIds(uid, params.ChatIds, params.ExceptMIds) {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, defErr)
	return
}

func DelSelfAll(ctx *gin.Context) {
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	defErr := errno.Add("params-err", errno.ParamsErr)
	if message.DeleteSelfAll(uid) {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, defErr)
	return
}

func RevokeByChatIds(ctx *gin.Context) {
	var params message.RevokeByChatIdsParams
	userId, _ := ctx.Get("uid")
	defErr := errno.Add("params-err", errno.ParamsErr)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, defErr)
		return
	}
	//os := funcs.GetHeaders(ctx)["Device-Id"][0]
	deviceId := ctx.Value(base.HeaderFieldDeviceId).(string)
	_, err := message.RevokeByChatIds(ctx, message.RevokeByChatIdsParams{
		ChatIds:    params.ChatIds,
		ExceptMIds: params.ExceptMIds,
		UID:        userId.(string),
		Os:         deviceId,
	})
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}

func Revoke(ctx *gin.Context) {
	var params message.RevokeParams
	userId, _ := ctx.Get("uid")
	defErr := errno.Add("params-err", errno.ParamsErr)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, defErr)
		return
	}
	os := funcs.GetHeaders(ctx)["Device-Id"][0]
	if _, res := message.Revoke(ctx, message.RevokeParams{
		UID:    userId.(string),
		ChatId: params.ChatId,
		MIds:   params.MIds,
		Os:     os,
	}); res {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, defErr)
	return
}

func ClearByChatIds(ctx *gin.Context) {
	var params message.RevokeByChatIdsParams
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	defErr := errno.Add("params-err", errno.ParamsErr)
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, defErr)
		return
	}
	os := funcs.GetHeaders(ctx)["Device-Id"][0]
	_, err := message.ClearByChatIds(ctx, message.RevokeByChatIdsParams{
		ChatIds:    params.ChatIds,
		ExceptMIds: params.ExceptMIds,
		UID:        uid,
		Os:         os,
	})
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}
