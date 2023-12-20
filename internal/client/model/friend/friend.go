package friend

import (
	"context"
	friend2 "imsdk/internal/common/dao/friendv2/friend"
	"imsdk/internal/common/dao/user"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
)

type Friend struct {
	ID         string   `json:"_id"`
	ObjID      string   `json:"obj_id"`
	IsStar     int8     `json:"is_star"`
	Alias      string   `json:"alias"`
	Phone      []string `json:"phone"`
	RemarkText string   `json:"remark_text"`
	Tag        []string `json:"tag"`
	ViewDetail int64    `json:"view_detail"`
	IsRead     int8     `json:"is_read"`
}

type GetFriendListsResponse struct {
	UId    string `json:"uid"`
	ChatId string `json:"chat_id"`
	Remark string `json:"remark"`
	PubKey string `json:"pub_key"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Sign   string `json:"sign"`
}

func GetFriendsIds(uid string) []string {
	data, _ := friend2.New().GetFriendIds(uid)
	return data
}

func GetFriendLists(ctx context.Context, uid string, request ListFriendsRequest) []GetFriendListsResponse {
	uIds, remarkInfo, _ := friend2.New().GetFriendInfos(uid, request.UIds)
	uInfos, _ := user.New().GetByIDs(uIds)
	res := make([]GetFriendListsResponse, 0)
	for _, datum := range uInfos {
		res = append(res, GetFriendListsResponse{
			UId:    datum.ID,
			Remark: remarkInfo[datum.ID],
			ChatId: funcs.CreateSingleChatId(uid, datum.ID),
			PubKey: datum.PubKey,
			Avatar: datum.Avatar,
			Name:   datum.Name,
			Gender: datum.Gender,
			Sign:   datum.Sign,
		})
	}
	return res
}

func GetMineByFriendIds(uid string, ids []string) []string {
	data := friend2.New().GetMineByFriendIds(uid, ids)
	return data
}

func GetFriendsDetail(uid string, ids []string) []friend2.Friend {
	data, _ := friend2.New().GetFriendsInfo(uid, ids)
	return data
}

func GetRelationInfo(uid string, ids []string) []GetRelationInfoResponse {
	friendIds, otherSideIds := friend2.New().GetRelationInfo(uid, ids)
	result := make([]GetRelationInfoResponse, 0)
	publicRelationUid := friend2.GetPublicRelationUIds()
	for _, v := range ids {
		isFriend := 0
		if funcs.In(v, publicRelationUid) {
			isFriend = 1
		} else if funcs.In(v, friendIds) && funcs.In(v, otherSideIds) {
			isFriend = 1
		} else if funcs.In(v, friendIds) && !funcs.In(v, otherSideIds) {
			isFriend = 2
		} else if !funcs.In(v, friendIds) && funcs.In(v, otherSideIds) {
			isFriend = 3
		}
		item := GetRelationInfoResponse{
			ID:       v,
			IsFriend: int8(isFriend),
		}
		result = append(result, item)
	}
	res := make([]GetRelationInfoResponse, 0)
	if len(result) > 0 {
		dao := user.New()
		for _, v := range result {
			u, _ := dao.GetInfoById(v.ID)
			v.Name = u.Name
			v.Avatar = u.Avatar
			v.Gender = u.Gender
			v.Sign = u.Sign
			res = append(res, v)
		}
		//mineBlockIds, _ := block.New().GetMineBlockedIds(uid, ids)
		//blockMineIds, _ := block.New().GetBlockedMineIds(uid, ids)
		//for _, v := range data {
		//	v["is_block_me"] = 0
		//	v["is_me_block"] = 0
		//	if funcs.In(v["id"].(string), mineBlockIds) {
		//		v["is_me_block"] = 1
		//	}
		//	if funcs.In(v["id"].(string), blockMineIds) {
		//		v["is_block_me"] = 1
		//	}
		//}
	}
	return res
}

func DelFriend(uid, objUid string) error {
	err := friend2.New().DelFriend(uid, objUid)
	return err
}

// DelFriendsUnilateral 单方面
func DelFriendsUnilateral(ctx context.Context, uid string, request DeleteFriendsRequest) error {
	err := friend2.New().DelFriendsUnilateral(uid, request.UIDs)
	return err
}

func DelAllUnilateral(ctx context.Context, uid string) error {
	err := friend2.New().DelAllUnilateral(uid)
	return err
}

// DelFriendsBilateral 双方面
func DelFriendsBilateral(ctx context.Context, uid string, request DeleteFriendsRequest) error {
	err := friend2.New().DelFriendsBilateral(uid, request.UIDs)
	return err
}

func DelAllBilateral(ctx context.Context, uid string) error {
	err := friend2.New().DelAllBilateral(uid)
	return err
}

func UpdateRemark(ctx context.Context, uid string, request UpdateRemarkRequest) error {
	upData := map[string]interface{}{
		"remark": request.Remark,
	}
	_, err := friend2.New().UpdateRemark(uid, request.Remark, upData)
	return err
}

func GetFriendInfoById(senderId, target string) friend2.Friend {
	data, _ := friend2.New().GetByID(friend2.New().GetID(senderId, target))
	return data
}

func GetRelationErrCode(senderId, target string) int {
	relationInfo := GetRelationInfo(senderId, []string{target})
	errCode := errno.OK
	if relationInfo[0].IsFriend == 2 {
		errCode = friend2.OtherSideNotFriend
	} else if relationInfo[0].IsFriend == 3 {
		errCode = friend2.MemberNotFriend
	} else if relationInfo[0].IsFriend == 0 {
		errCode = friend2.EachOtherStrangers
	}
	return errCode
}
