package group

import (
	"encoding/json"
	"fmt"
	"imsdk/internal/client/model/group"
	"imsdk/internal/common/dao/group/members"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"

	"github.com/gin-gonic/gin"
)

func CreateGroup(ctx *gin.Context) {
	var params group.CreateRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid := ctx.Value(base.HeaderFieldUID).(string)
	if err = group.Create(ctx, uid, params); err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func GetMembersInfo(ctx *gin.Context) {
	var params group.MembersRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	userId, _ := ctx.Get("uid")
	res, er := group.GetMembersInfo(ctx, userId.(string), params)
	if er != nil {
		response.RespErr(ctx, er)
		return
	}
	response.RespListData(ctx, res)
	return
}

func GetMembersByIds(ctx *gin.Context) {
	var params group.GetMembersByIdsRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid := ctx.Value(base.HeaderFieldUID).(string)
	res, er := group.GetMembersByIds(ctx, uid, params)
	if er != nil {
		response.RespErr(ctx, er)
		return
	}
	response.RespListData(ctx, res)
	return
}

func GetEncInfoByIds(ctx *gin.Context) {
	var params group.IdsRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid := ctx.Value(base.HeaderFieldUID).(string)
	res, er := group.GetEncInfoByIds(ctx, uid, params)
	if er != nil {
		response.RespErr(ctx, er)
		return
	}
	response.RespListData(ctx, res)
	return
}

func Join(ctx *gin.Context) {
	var params group.JoinRequest
	uid := ctx.Value(base.HeaderFieldUID).(string)
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	if err = group.Join(ctx, uid, params); err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func AgreeJoin(ctx *gin.Context) {
	var params group.AgreeJoinRequest
	data, _ := ctx.Get("data")
	uid := ctx.Value(base.HeaderFieldUID).(string)
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.AgreeJoin(ctx, uid, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func AdminAgreeJoin(ctx *gin.Context) {
	var params group.AgreeJoinRequest
	data, _ := ctx.Get("data")
	uid := ctx.Value(base.HeaderFieldUID).(string)
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.AdminAgreeJoin(ctx, uid, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func GetList(ctx *gin.Context) {
	var params group.IdsRequest
	uid := ctx.Value(base.HeaderFieldUID).(string)
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	res, _ := group.GetListByUid(ctx, uid, params)
	response.RespListData(ctx, res)
	return
}

func GetMemberIds(ctx *gin.Context) {
	var params group.IdRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, _ := members.New().GetMemberIds(params.Id)
	response.ResData(ctx, data)
	return
}

func GetGroupMemberInfoByIds(ctx *gin.Context) {
	var params group.IdsRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	//data, _ := group.GetGroupsMemberInfoByIds(params.GroupIDs)
	//response.RespData(ctx, data)
	return
}

func GetGroupMemberInfoByUIds(ctx *gin.Context) {
	var params group.GetGroupMemberInfoByUIdsRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	//data, _ := group.GetGroupMemberInfoByUIds(ctx, params)
	//response.RespListData(ctx, data)
	return
}

func GetGroupsMemberIds(ctx *gin.Context) {
	var params group.IdsRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, _ := group.GetGroupsMemberIds(params.Ids)
	response.RespData(ctx, data)
	return
}

func InviteJoin(ctx *gin.Context) {
	var params group.InviteJoinRequest
	uid := ctx.Value(base.HeaderFieldUID).(string)
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.InviteJoin(ctx, uid, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func UpdateName(ctx *gin.Context) {
	var params group.UpdateNameRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.UpdateName(ctx, userId.(string), params)
	if err == nil {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, err)
	return
}

func UpdateAlias(ctx *gin.Context) {
	var params group.UpdateAliasRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.UpdateAlias(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func Transfer(ctx *gin.Context) {
	var params group.TransferRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.TransferGroup(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func Disband(ctx *gin.Context) {
	var params group.IdRequest
	userId, _ := ctx.Get("uid")
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	err = group.DisbandGroup(ctx, userId.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func ClearMessage(ctx *gin.Context) {
	var params group.ClearMessageRequest
	fmt.Println(params)
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

func ApplyList(ctx *gin.Context) {
	var params group.IdsRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid := ctx.Value(base.HeaderFieldUID).(string)
	res, er := group.ApplyList(ctx, uid, params)
	if er != nil {
		response.RespErr(ctx, er)
		return
	}
	response.RespListData(ctx, res)
	return
}

func MyApplyList(ctx *gin.Context) {
	var params group.IdsRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid := ctx.Value(base.HeaderFieldUID).(string)

	res, _ := group.MyApplyList(ctx, uid, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespListData(ctx, res)
	return
}

func DetailByIds(ctx *gin.Context) {
	var params group.IdsRequest
	uid := ctx.Value(base.HeaderFieldUID).(string)
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	res, _ := group.GetDetailByIds(ctx, uid, params)
	if err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	response.RespListData(ctx, res)
	return
}
func GetQrCode(ctx *gin.Context) {
	var params struct {
		GroupID string `json:"id" binding:"required"`
	}
	//uid, _ := ctx.Get("uid")
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	//res := group.CreateQrCode(uid.(string), params.GroupID)
	//if res == nil {
	//	response.RespErr(ctx, errno.Add("fail", errno.DefErr))
	//	return
	//}
	//response.ResData(ctx, res)
	return
}

func GetListNameByIds(ctx *gin.Context) {
	var params struct {
		GIDs []string `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	uid, _ := ctx.Get("uid")
	data, _ := group.GetGroupName(uid.(string), params.GIDs)
	fmt.Println("group controller", data)
	response.RespData(ctx, data)
	return
}

//func GetSnapshotList(ctx *gin.Context) {
//	var params struct {
//		GIds []string `json:"ids" binding:"required"`
//	}
//	if err := ctx.ShouldBindJSON(&params); err != nil {
//		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
//		return
//	}
//	uid, _ := ctx.Get("uid")
//	data := group.GetSnapshotList(uid.(string), params.GIds)
//	response.RespListData(ctx, data)
//	return
//}
