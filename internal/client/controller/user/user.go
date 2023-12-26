package user

import (
	"encoding/json"
	"fmt"
	"imsdk/internal/client/model/user"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetAuthInfo(ctx *gin.Context) {
	var params user.GetUserInfoRequest
	uid, _ := ctx.Get(base.HeaderFieldUID)
	params.Ids = []string{uid.(string)}
	list, err := user.GetUserInfo(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.ResEnData(ctx, list[0])
	return
}

func GetListInfo(ctx *gin.Context) {
	var params user.GetUserInfoRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	list, err := user.GetUserInfo(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, list)
	return
}

func UpdateName(ctx *gin.Context) {
	var params user.UpdateNameRequest
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	err := user.UpdateName(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func UpdateGender(ctx *gin.Context) {
	var params user.UpdateGenderRequest
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	uid := ctx.Value(base.HeaderFieldUID).(string)
	err := user.UpdateGender(ctx, uid, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func UpdateSign(ctx *gin.Context) {
	var params user.UpdateSignRequest
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	uid := ctx.Value(base.HeaderFieldUID).(string)
	err := user.UpdateSign(ctx, uid, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func UpdateAvatar(ctx *gin.Context) {
	var params user.UpdateAvatarRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = user.UpdateAvatar(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func Unsubscribe(ctx *gin.Context) {
	uid := ctx.Value(base.HeaderFieldUID).(string)
	fmt.Println("Unsubscribe uid:", uid)
	var err error
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
