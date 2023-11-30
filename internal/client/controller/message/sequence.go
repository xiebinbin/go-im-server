package message

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/message"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/response"
)

type UserMessageIdRequest = message.UserMessageIdsParams
type UserMessageIdResponse struct {
	ID       string `json:"_id" json:"id"` // base id
	Sequence int64  `json:"sequence" json:"sequence"`
}

func GetMessageIds(ctx *gin.Context) {
	sT := funcs.GetMillis()
	fmt.Println("getMessageIds request startTime:", sT)
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	var params UserMessageIdRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	if params.Sequence < 0 {
		response.RespErr(ctx, errno.Add("params-err-", errno.ParamsErr))
		return
	}
	fmt.Println("getMessageIds request params.Sequence:", params.Sequence, ":uid:", uid)
	data := message.GetMsgIds(ctx, uid, params.Sequence)
	eT := funcs.GetMillis()
	fmt.Println("getMessageIds request endTime:", eT, ":-diffTime-:", eT-sT, "data len :", len(data))
	response.ResData(ctx, data)
	return
}

func GetMessageIdsV2(ctx *gin.Context) {
	sT := funcs.GetMillis()
	fmt.Println("getMessageIds request startTime:", sT)
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	var params UserMessageIdRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	if params.Sequence < 0 {
		response.RespErr(ctx, errno.Add("params-err-", errno.ParamsErr))
		return
	}
	data := message.GetMsgIdsV2(ctx, uid, params.Sequence)
	res := map[string]interface{}{
		"items":  data,
		"is_end": len(data) < base.MsgPageRowLimit,
	}
	eT := funcs.GetMillis()
	fmt.Println("getMessageIds request endTime:", eT, ":-diffTime-:", eT-sT, "data len :", len(data))
	response.RespOriginData(ctx, res)
	return
}

func GetUserMaxSequence(ctx *gin.Context) {
	userId, _ := ctx.Get("uid")
	seq := message.GetMaxSequence(userId.(string))
	response.RespData(ctx, map[string]interface{}{
		"sequence": seq,
	})
	return
}
