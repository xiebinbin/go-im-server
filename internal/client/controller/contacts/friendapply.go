package contacts

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/friend"
	"imsdk/internal/common/dao/friendv2/apply"
	user2 "imsdk/internal/common/dao/user"
	"imsdk/internal/common/pkg/base"
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

func GetApplyLists(ctx context.Context) ([]friend.ApplyListResponse, error) {
	uid := ctx.Value(base.HeaderFieldUID).(string)
	data, err := apply.New().GetApplyLists(uid)
	var uIds []string
	if len(data) > 0 {
		for _, datum := range data {
			uIds = append(uIds, datum.UId)
			uIds = append(uIds, datum.ObjUId)
		}
	}
	res := make([]friend.ApplyListResponse, 0)
	uInfo, _ := user2.New().GetInfoByIds(uIds)
	if err != nil {
		return res, err
	}
	for _, datum := range data {
		var Name string
		var Avatar string
		if uid == datum.UId {
			Name = uInfo[datum.ObjUId].Name
			Avatar = uInfo[datum.ObjUId].Avatar
		} else {
			Name = uInfo[datum.UId].Name
			Avatar = uInfo[datum.UId].Avatar
		}
		res = append(res, friend.ApplyListResponse{
			ID:        datum.ID,
			ObjUId:    datum.ObjUId,
			UId:       datum.UId,
			Remark:    datum.Remark,
			CreatedAt: datum.CreatedAt,
			Status:    datum.Status,
			IsRead:    datum.IsRead,
			Name:      Name,
			Avatar:    Avatar,
		})
	}
	return res, err
}
