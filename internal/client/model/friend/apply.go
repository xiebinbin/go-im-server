package friend

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"imsdk/internal/common/dao/friendv2/apply"
	user2 "imsdk/internal/common/dao/user"
	"imsdk/internal/common/model/chat"
	"imsdk/internal/common/model/user"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
)

type ApplyRequest struct {
	ObjUID string `json:"obj_uid"`
	Remark string `json:"remark"`
}

type AddFriendRequest struct {
	UID    string `json:"uid" binding:"required"`
	ObjUID string `json:"obj_uid" binding:"required"`
}

type ReadApplyFriendsRequest struct {
	Ids []string `json:"ids" binding:"required"`
}

type AddApplyResp struct {
	Result int8 `json:"result"`
}

type AgreeRequest struct {
	ID     string `json:"id"`
	UID    string `json:"uid"`
	ObjUID string `json:"obj_uid"`
}

type RefuseRequest struct {
	ID     string `json:"id"`
	UID    string `json:"uid"`
	ObjUID string `json:"obj_uid"`
}

type ApplyListResponse struct {
	ID        string `bson:"_id" json:"id,omitempty"`
	UId       string `bson:"uid" json:"uid,omitempty"`
	Remark    string `bson:"remark" json:"remark,omitempty"`
	Status    int8   `bson:"status" json:"status,omitempty"`
	IsRead    int8   `bson:"is_read" json:"is_read,omitempty"`
	Avatar    string `bson:"avatar" json:"avatar"`
	Name      string `bson:"name" json:"name"`
	CreatedAt int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt int64  `bson:"update_time" json:"update_time,omitempty"`
}

const (
	StatusUnApply          = errno.DataNotExist
	ErrCodeAddWayNotAllow  = 200000
	DistinctUserInfo       = 1
	DistinctChat           = 2
	ResultApplySuccess     = 1
	ResultAddFriendSuccess = 2
)

func ReadApplyFriends(ctx context.Context, params ReadApplyFriendsRequest) error {
	//uid := ctx.Value("uid")
	uData := map[string]interface{}{
		"is_read": apply.IsReadYes,
	}
	return apply.New().UpdateInfoByIds(params.Ids, uData)
}

func AddApply(ctx context.Context, uid string, params ApplyRequest) (AddApplyResp, error) {
	res := AddApplyResp{
		Result: ResultApplySuccess,
	}
	t := funcs.GetMillis()
	if uid == params.ObjUID {
		return res, errno.Add("Not allowed to add myself as a friend", errno.FORBIDDEN)
	}
	if err := user.GetUserErr(params.ObjUID); err != nil {
		fmt.Println("err:3", err)
		return res, err
	}
	relationInfo := GetRelationInfo(uid, []string{params.ObjUID})
	if relationInfo[0].IsFriend == 1 {
		return res, nil
	}
	applyMeInfo, _ := GetRecentNotAgree(params.ObjUID, uid)
	isApplyMe := (applyMeInfo != apply.Apply{}) && (applyMeInfo.CreatedAt > t-7*24*3600*1000)
	if isApplyMe {
		err := AgreeApply(ctx, AgreeRequest{
			UID:    uid,
			ObjUID: params.ObjUID,
		})
		if err != nil {
			return res, errno.Add("add-fail", errno.SysErr)
		}
		res.Result = ResultAddFriendSuccess
		return res, nil
	}
	if relationInfo[0].IsFriend == 3 {
		_, err1 := chat.Create(ctx, chat.CreateParams{
			Creator: uid,
			UIds:    []string{uid, params.ObjUID},
		})
		if err1 != nil {
			return res, err1
		}
		err := AddEachOther(AddFriendRequest{
			UID:    uid,
			ObjUID: params.ObjUID,
		})
		if err != nil {
			return res, errno.Add("add-fail", errno.SysErr)
		}
		res.Result = ResultAddFriendSuccess
		return res, nil
	}
	applyInfo, _ := GetRecentNotAgree(uid, params.ObjUID)
	isMyApply := (applyInfo != apply.Apply{}) && (applyInfo.CreatedAt > t-7*24*3600*1000)
	lockKey := uid + params.ObjUID
	redis.Lock(lockKey, 1)
	if isMyApply {
		uData := map[string]interface{}{
			"update_time": funcs.GetMillis(),
		}
		return res, apply.New().UpdateInfoById(applyInfo.ID, uData)
	}
	err := apply.New().AddApply(apply.Apply{
		UId:    uid,
		ObjUId: params.ObjUID,
		Remark: params.Remark,
	})
	redis.Unlock(lockKey)
	return res, err
}

func GetApplyLists(ctx context.Context) ([]ApplyListResponse, error) {
	uid := ctx.Value(base.HeaderFieldUID).(string)
	data, err := apply.New().GetApplyLists(uid)
	var uIds []string
	if len(data) > 0 {
		for _, datum := range data {
			uIds = append(uIds, datum.UId)
		}
	}
	res := make([]ApplyListResponse, 0)
	uInfo, _ := user2.New().GetInfoByIds(uIds)
	if err != nil {
		return res, err
	}
	for _, datum := range data {
		res = append(res, ApplyListResponse{
			ID:        datum.ID,
			UId:       datum.UId,
			Remark:    datum.Remark,
			CreatedAt: datum.CreatedAt,
			Status:    datum.Status,
			IsRead:    datum.IsRead,
			Name:      uInfo[datum.UId].Name,
			Avatar:    uInfo[datum.UId].Avatar,
		})
	}
	return res, err
}

func GetApplyInfo(uid, objUId string) (apply.Apply, error) {
	data, err := apply.New().GetApplyInfo(uid, objUId)
	return data, err
}

func GetRecentNotAgree(uid, objUId string) (apply.Apply, error) {
	data, err := apply.New().GetNotAgreeInfo(uid, objUId)
	return data, err
}

func RefuseApply(ctx context.Context, request RefuseRequest) error {
	if err := user.GetUserErr(request.ObjUID); err != nil {
		return err
	}
	err := apply.New().Delete(request.ID)
	return err
}

func AgreeApply(ctx context.Context, request AgreeRequest) error {
	uid := ctx.Value(base.HeaderFieldUID).(string)
	applyDao := apply.New()
	applyInfo, err := applyDao.GetInfoByID(request.ID)
	uData := map[string]interface{}{
		"status":      apply.StatusPass,
		"update_time": funcs.GetMillis(),
	}
	if err != nil {
		// Requirement requires direct return of success
		return nil
		//return errno.Add("unapply", StatusUnApply)
	}

	err = applyDao.UpdateInfoById(applyInfo.ID, uData)
	if err != nil {
		return err
	}
	// Garbage demand
	_, err = chat.Create(ctx, chat.CreateParams{
		Creator: uid,
		UIds:    []string{applyInfo.UId, uid},
	})
	if err != nil {
		return err
	}
	if err = AddEachOther(AddFriendRequest{
		UID:    applyInfo.UId,
		ObjUID: applyInfo.ObjUId,
	}); err != nil {
		return err
	}
	//err = SendSysMsgAfterAgree(request.UID, request.ObjUID)
	//if err != nil {
	//	return err
	//}
	return nil
}

func SendSysMsgAfterAgree(uid, objUId string) error {
	sysMsgData := map[string]interface{}{
		"operator": uid,
		"target":   []string{objUId},
		"temId":    "add-friend",
	}
	contentByte, _ := json.Marshal(sysMsgData)
	fmt.Println("contentByte:", string(contentByte))
	return nil
}
