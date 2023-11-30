package contacts

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/friend"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
)

type ApplyRequest = friend.ApplyRequest
type ReadApplyFriendsRequest = friend.ReadApplyFriendsRequest

func Apply(ctx *gin.Context) {
	var params friend.ApplyRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	uid, _ := ctx.Get("uid")
	_, err = friend.AddApply(ctx, uid.(string), params)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	// errcode
	response.RespErr(ctx, err)
	return
}

func ReadApplyFriends(ctx *gin.Context) {
	//userId, _ := ctx.Get("uid")
	var params ReadApplyFriendsRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err := friend.ReadApplyFriends(ctx, params)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}

func GetList(ctx *gin.Context) {
	data, _ := friend.GetApplyLists(ctx)
	response.RespListData(ctx, data)
	return
}

func Agree(ctx *gin.Context) {
	var params friend.AgreeRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err = friend.AgreeApply(ctx, params); err == nil {
		response.RespSuc(ctx)
	} else {
		response.RespErr(ctx, err)
	}
	return
}

func Refuse(ctx *gin.Context) {
	var params friend.RefuseRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	params.UID = ctx.Value("uid").(string)
	if err = friend.RefuseApply(ctx, params); err == nil {
		response.RespSuc(ctx)
	} else {
		response.RespErr(ctx, err)
	}
	return
}
