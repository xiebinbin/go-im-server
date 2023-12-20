package contacts

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/friend"
	"imsdk/internal/client/model/user/contact"
	friend2 "imsdk/internal/common/dao/friendv2/friend"
	contact2 "imsdk/internal/common/dao/user/contact"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/response"
	"strings"
)

type InfoRequest struct {
	Ids []string `json:"ids" binding:"required"`
}

type InfoResponse struct {
	ID         string                   `json:"id"`
	ObjUID     string                   `json:"obj_uid"`
	IsStar     int8                     `json:"is_star"`
	Alias      string                   `json:"alias"`
	Phone      []string                 `json:"phone"`
	RemarkText string                   `json:"remark_text"`
	RemarkImg  string                   `json:"remark_img"`
	Tag        []map[string]interface{} `json:"tag"`
	Phones     []map[string]interface{} `json:"phones,omitempty"`
	Relations  []map[string]interface{} `json:"relations,omitempty"`
	Emails     []map[string]interface{} `json:"emails,omitempty"`
	Dates      []map[string]interface{} `json:"dates,omitempty"`
	Companies  []map[string]interface{} `json:"companies,omitempty"`
	Schools    []map[string]interface{} `json:"schools,omitempty"`
	Address    []map[string]interface{} `json:"addrs,omitempty"`
	ViewDetail int64                    `json:"view_detail"`
	IsRead     int8                     `json:"is_read"`
	FromWay    int8                     `json:"from_way"`
}

func GetContactsIds(ctx *gin.Context) {
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	data := friend.GetFriendsIds(uid)
	response.RespData(ctx, data)
	return
}

func GetFriendList(ctx *gin.Context) {
	userId, _ := ctx.Value(base.HeaderFieldUID).(string)
	var params friend.ListFriendsRequest
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	res := friend.GetFriendLists(ctx, userId, params)
	response.RespListData(ctx, res)
	return
}

func DelFriendsUnilateral(ctx *gin.Context) {
	var params friend.DeleteFriendsRequest
	userId, _ := ctx.Get(base.HeaderFieldUID)
	uid := userId.(string)
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	err := friend.DelFriendsUnilateral(ctx, uid, params)
	if len(params.UIDs) == 0 {
		err = friend.DelAllUnilateral(ctx, uid)
	}
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func DelFriendsBilateral(ctx *gin.Context) {
	var params friend.DeleteFriendsRequest
	userId, _ := ctx.Get(base.HeaderFieldUID)
	uid := userId.(string)
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	err := friend.DelFriendsBilateral(ctx, uid, params)
	if len(params.UIDs) == 0 {
		err = friend.DelAllBilateral(ctx, uid)
	}
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func UpdateRemark(ctx *gin.Context) {
	var params friend.UpdateRemarkRequest
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	uid, _ := ctx.Get(base.HeaderFieldUID)
	err := friend.UpdateRemark(ctx, uid.(string), params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func GetContactsListInfo(ctx *gin.Context) {
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	var params InfoRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	resData := make([]InfoResponse, 0)
	if len(params.Ids) == 0 {
		response.RespData(ctx, resData)
		return
	}
	contactsData := contact.GetContactsDetail(uid, params.Ids)
	friendsData := friend.GetFriendsDetail(uid, params.Ids)
	tmp, tmpContact := make(map[string]interface{}, 0), make(map[string]interface{}, 0)
	if friendsData != nil {
		for _, v := range friendsData {
			tmp[v.ObjUId] = v
		}
	}
	if contactsData != nil {
		for _, v := range contactsData {
			tmpContact[v.ObjUId] = v
		}
	}
	for _, v := range params.Ids {
		var c contact2.Contact
		var info friend2.Friend
		tmpResp := InfoResponse{
			ID:     uid,
			ObjUID: v,
			Tag:    []map[string]interface{}{},
		}
		if tmpContact[v] != nil {
			c = tmpContact[v].(contact2.Contact)
			tmpResp.Alias = c.Alias
			tmpResp.Phone = strings.Split(c.Phone, ",")
			tmpResp.RemarkText = c.RemarkText
			tmpResp.Phones = c.Phones
			tmpResp.Relations = c.Relations
			tmpResp.Emails = c.Emails
			tmpResp.Dates = c.Dates
			tmpResp.Schools = c.Schools
			tmpResp.Companies = c.Companies
			tmpResp.Address = c.Address
		}
		if tmp[v] != nil {
			info = tmp[v].(friend2.Friend)
			tmpResp.IsStar = info.IsStar
		}
		//if c.Tag != "" {
		//	tagTmp, _ := tag.New().GetTagInfo(strings.Split(c.Tag, ","))
		//}
		resData = append(resData, tmpResp)
	}
	response.ResData(ctx, resData)
	return
}
