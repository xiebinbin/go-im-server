package group

import (
	"context"
	"fmt"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao/chat/changelogs"
	user2 "imsdk/internal/common/dao/user"
	"imsdk/internal/common/dao/user/infov2/info"
	"imsdk/internal/common/model/chat"
	"imsdk/internal/common/model/user"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"sort"
	"strings"

	"imsdk/internal/common/dao/group/detail"
	"imsdk/internal/common/dao/group/members"

	"imsdk/pkg/errno"
)

type CreateRequest struct {
	Id      string   `json:"id" binding:"required"`
	Avatar  string   `json:"avatar"`
	PubKey  string   `json:"pub_key"`
	Name    string   `json:"name" binding:"required,gt=3,lte=64"`
	Members []string `json:"members" binding:"required"`
}

type GetGroupMemberInfoByUIdsRequest struct {
	GroupMember []GetGroupMemberInfoByUIdsItem `json:"group_member" binding:"required"`
}

type GetGroupMemberInfoByUIdsItem struct {
	GroupID string   `json:"gid" binding:"required"`
	UIDs    []string `json:"uids"`
}

type BaseDetail struct {
	ID     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	IsIn   int    `json:"is_in"`
}

const (
	MemExceedMax = 100100
)

func Create(ctx context.Context, uid string, params CreateRequest) (err error) {
	targetMembers := params.Members
	t := funcs.GetMillis()
	membersDao := members.New()
	//memList := make([]members.Members, 0)
	allMemIds := make([]string, len(params.Members))
	copy(allMemIds, params.Members)
	allMemIds = append(allMemIds, uid)
	allMemIds = funcs.RemoveDuplicatesAndEmpty(allMemIds)

	// save group detail data
	groupDetail := detail.Detail{
		ID:         params.Id,
		CreatorUid: uid,
		OwnerUid:   uid,
		Name:       params.Name,
		Avatar:     params.Avatar,
		Total:      len(allMemIds),
		CreatedAt:  t,
		UpdatedAt:  t,
	}
	if len(allMemIds) < 3 {
		return errno.Add("Group of at least 3 user", errno.MissingParams)
	}
	if err = detail.New().Add(groupDetail); err != nil && !mongoDriver.IsDuplicateKeyError(err) {
		//app.Logger().WithField("action", "create group").WithField("desc", "save group detail unsuccessfully").Error(err)
		return errno.Add("failed to save group", errno.Exception)
	}
	// save members of group data
	var uIds []string
	var delIds []string
	for _, u := range allMemIds {
		if err = user.GetUserErr(u); err != nil {
			delIds = append(delIds, u)
			continue
		}
		item := members.Members{
			ID:        membersDao.GetId(u, params.Id),
			GroupID:   params.Id,
			UID:       u,
			Role:      members.RoleCommonMember,
			JoinType:  members.JoinTypeInvite,
			InviteUID: uid,
			Status:    members.StatusYes,
			CreatedAt: t,
			UpdatedAt: t,
		}
		if u == uid {
			item.Role = members.RoleOwner
		}
		er := membersDao.UpsertOne(item)
		if er != nil {
			fmt.Println("item and err:", item, err)
			//app.Logger().WithField("action", "create group").WithField("desc", "save members unsuccessfully").Error(err)
			return errno.Add("failed to save group", errno.Exception)
		}
		uIds = append(uIds, u)
		//memList = append(memList, item)
	}
	if len(uIds) == 1 {
		return errno.Add("user-status-delete", errno.UserDelete)
	}

	if err != nil {
		return err
	}
	_, er1 := chat.Create(ctx, chat.CreateParams{
		Id:       params.Id,
		Creator:  uid,
		OwnerUid: uid,
		UIds:     allMemIds,
		Avatar:   params.Avatar,
		Name:     params.Name,
	})
	if er1 != nil {
		return err
	}
	//_, err = group.SetDefaultInfoByUIds(uIds, params.Id)
	//if err != nil {
	//	return err
	//}
	// send system message by hollow man （模板消息：invite-join-group）
	targetMembers = funcs.SliceMinus(targetMembers, delIds)
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   targetMembers,
	//	"number":   len(targetMembers),
	//	"temId":    "invite-join-group",
	//}
	// todo should use queue and must be send successfully
	//actionData := map[string]interface{}{
	//	"name": "update_group_member",
	//}
	//err = message.HollowManSendGroupSystemMsg(params.Id, sysMsgData, actionData)
	//if err != nil {
	//	return err
	//}
	return nil
}

func QuitAllGroup(ctx context.Context, uid string, request QuitAllRequest) error {
	data, err := members.New().GetMyGroupIdList(uid)
	if err != nil {
		return errno.Add("sys-err", errno.SysErr)
	}
	if len(data) > 0 {
		for _, datum := range data {
			err = QuitGroup(ctx, uid, QuitRequest{
				GroupID:   datum.ID,
				IsDelChat: request.IsDelChat,
			})
			if err != nil {
				return errno.Add("sys-err", errno.SysErr)
			}
		}
	}
	return nil
}

func QuitGroup(ctx context.Context, uid string, request QuitRequest) error {
	groupInfo, err := detail.New().GetByID(request.GroupID, "id,status,name")
	if err != nil || groupInfo.ID == "" || groupInfo.Status != detail.StatusYes {
		return errno.Add("group-not-exist", errno.DataNotExist)
	}
	memCount, _ := members.New().GetGroupMemberIds(request.GroupID)
	if len(memCount) == 1 {
		err = ChangeGroupStatus(uid, request.GroupID, detail.StatusDel)
		if members.New().Delete(uid, request.GroupID) != nil {
			return errno.Add("fail", errno.DefErr)
		}
		return err
	}
	memInfo, _ := members.New().GetByUidAndGid(uid, request.GroupID, "role")
	if memInfo.Role == members.RoleOwner {
		// todo transferOwner
		newOwner := GetRandomGroupOwner(request.GroupID)
		memDao := members.New()
		if err = detail.New().OwnerQuit(request.GroupID, newOwner); err != nil {
			return errno.Add("failed to update group info", errno.DefErr)
		}

		// new owner
		err = memDao.UpdateMemberRole(newOwner, request.GroupID, members.RoleOwner)
		err = memDao.UpdateMemberRole(uid, request.GroupID, members.RoleCommonMember)
		if err != nil {
			return errno.Add("failed to update new owner", errno.DefErr)
		}
		er1 := changelogs.New().UpdateMemberInfo(request.GroupID, newOwner)
		if er1 != nil {
			return er1
		}
		//target := []string{newOwner}
		//sysMsgData := map[string]interface{}{
		//	"operator": uid,
		//	"target":   target,
		//	"temId":    "transfer-owner",
		//}
		redis.Client.Del(funcs.Md516(request.GroupID))
		// todo should use queue and must be send successfully
		//err = message.HollowSendOnlyGroupOwnerMsg(gid, sysMsgData, "")
		//if err != nil {
		//	return err
		//}
		//message.HollowSendSingleMsgInGroup(gid, newOwner, sysMsgData)
	}
	if members.New().Delete(uid, request.GroupID) != nil {
		return errno.Add("fail", errno.DefErr)
	}
	//conf, _ := config.GetIMSdkKey()
	//er1 := imsdkChat.NewClient(&imsdkChat.Options{
	//	Credentials: imsdk.NewStaticCredentials(conf.AK, conf.SK),
	//}).Remove(gid, []string{uid})
	//if er1 != nil {
	//	return er1
	//}
	err = chat.RemoveMember(context.Background(), chat.RemoveMembersParams{
		ChatId: request.GroupID,
		UIds:   []string{uid},
	})
	fmt.Println("err----", err)
	if err != nil {
		return err
	}

	detail.New().UpTotal(request.GroupID, -1)
	groupInfo.Total = groupInfo.Total - 1
	if request.IsDelChat == 1 {
		//_ = generateSnapshot(uid, groupInfo)
		//todo : send delete chat msg
		_ = QuitGroupDelChat(context.Background(), uid, request.GroupID)
	}
	er2 := changelogs.New().DelMemberInfo(request.GroupID, uid)
	if er2 != nil {
		return er2
	}

	// send system message by hollow man （format：quit-group）
	//var target []string
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   target,
	//	"temId":    "quit-group",
	//}
	// todo should use queue and must be send successfully
	//err = message.HollowSendOnlyGroupOwnerMsg(gid, sysMsgData, "")
	//if err != nil {
	//	return err
	//}
	return nil
}

func AgreeJoin(ctx context.Context, request AgreeJoinRequest) error {
	_, err := verifyGroupAndIsOwner(request.UID, request.GroupID, true)
	if err != nil {
		return err
	}
	dao := members.New()
	if len(request.ObjUid) > 0 {
		var ids []string
		for _, s := range request.ObjUid {
			ids = append(ids, dao.GetId(s, request.GroupID))
		}
		err = dao.UpByIDs(ids, members.Members{Status: members.StatusYes})
		if err != nil {
			return errno.Add("sys-err", errno.SysErr)
		}
	}
	return nil
}

func InviteJoin(ctx context.Context, uid string, request InviteJoinRequest) error {
	groupInfo, err := detail.New().GetByID(request.GroupID, "id,status,total")
	if err != nil || groupInfo.ID == "" || groupInfo.Status != detail.StatusYes {
		return errno.Add("group-not-exist", errno.DataNotExist)
	}
	if !members.New().IsExist(uid, request.GroupID) {
		return errno.Add("user-not-exist", errno.UserNotExist)
	}
	groupMemberInfo, _ := members.New().GetGroupMember(request.GroupID)
	ids, err := user2.New().GetInfoByStatus(request.ObjUid, []int{info.StatusDelete}, "_id,status")
	if err != nil {
		return err
	}
	var deleteIds []string
	if len(ids) > 0 {
		for _, v := range ids {
			deleteIds = append(deleteIds, v.ID)
		}
	}
	objUid := funcs.SliceMinus(request.ObjUid, deleteIds)
	if len(objUid) == 0 {
		return errno.Add("user-status-delete", errno.UserDelete)
	}
	if (len(groupMemberInfo) + len(objUid)) > detail.TotalMax {
		return errno.Add("group-exceed-limit", MemExceedMax)
	}
	membersDao := members.New()
	memList := make([]members.Members, 0)
	t := funcs.GetMillis()
	status := members.StatusIng
	if _, err = verifyGroupAndIsOwner(uid, request.GroupID, true); err == nil {
		status = members.StatusYes
	}
	for _, v := range objUid {
		item := members.Members{
			ID:        membersDao.GetId(v, request.GroupID),
			GroupID:   request.GroupID,
			UID:       v,
			Role:      members.RoleCommonMember,
			JoinType:  members.JoinTypeInvite,
			InviteUID: uid,
			Status:    int8(status),
			CreatedAt: t,
			UpdatedAt: t,
		}
		if v == uid {
			item.Role = members.RoleOwner
		}
		err = membersDao.Save(item)
		if err != nil {
			//app.Logger().WithField("action", "invite join group").WithField("desc", "failed to save members").Error(err)
			continue
		}
		memList = append(memList, item)
	}
	// join chat
	err = chat.JoinChat(ctx, chat.JoinChatParams{
		ChatId:    request.GroupID,
		InviteUID: uid,
		UIds:      objUid,
	})
	if err != nil {
		return err
	}

	// send system message by hollow man （模板消息：invite-join-group）
	detail.New().UpTotal(request.GroupID, len(objUid))
	defer func() {
		if err1 := recover(); err1 != nil {
			//app.Logger().WithField("action", "invite join group").WithField("desc", "send msg unsuccessfully").Error(err)
		}
	}()
	redis.Client.Del(funcs.Md516(request.GroupID))
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   objUid,
	//	"temId":    "invite-join-group",
	//}
	//actionData := map[string]interface{}{
	//	"name": "update_group_member",
	//}
	//// todo should use queue and must be send successfully
	//err = message.HollowManSendGroupSystemMsg(gid, sysMsgData, actionData)
	//if err != nil {
	//	return err
	//}
	changelogs.New().UpdateManyMemberInfo(request.GroupID, objUid)
	return nil
}

func verifyGroupAndIsAdministrator(uid, gid string, isVerifyAdministrator bool) (detail.Detail, error) {
	groupInfo, err := detail.New().GetByID(gid, "id,status,name")
	if err != nil || groupInfo.ID == "" || groupInfo.Status != detail.StatusYes {
		return groupInfo, errno.Add("group-not-exist", errno.DataNotExist)
	}
	memInfo, err := members.New().GetByUidAndGid(uid, gid, "id,role")
	if memInfo == (members.Members{}) || err != nil || (isVerifyAdministrator && memInfo.Role == members.RoleCommonMember) {
		return groupInfo, errno.Add("forbidden", errno.FORBIDDEN)
	}
	return groupInfo, nil
}

func verifyGroupAndIsOwner(uid, gid string, isVerifyOwner bool) (detail.Detail, error) {
	groupInfo, err := detail.New().GetByID(gid, "id,status")
	if err != nil || groupInfo.ID == "" || groupInfo.Status != detail.StatusYes {
		return groupInfo, errno.Add("group-not-exist", errno.DataNotExist)
	}

	memInfo, err := members.New().GetByUidAndGid(uid, gid, "id,role")

	if memInfo == (members.Members{}) || err != nil || (isVerifyOwner && memInfo.Role != members.RoleOwner) {
		return groupInfo, errno.Add("forbidden", errno.FORBIDDEN)
	}
	return groupInfo, nil
}

func UpdateName(ctx context.Context, uid string, request UpdateNameRequest) error {
	_, err := verifyGroupAndIsAdministrator(uid, request.GroupID, true)
	if err != nil {
		return err
	}
	//if !members.Save().IsExist(uid, gid) {
	//	return errno.Save("forbidden", errno.FORBIDDEN, "*")
	//}
	uData := detail.Detail{
		Name:      request.Name,
		UpdatedAt: funcs.GetMillis(),
	}
	if err = detail.New().UpByID(request.GroupID, uData); err != nil {
		return errno.Add("fail", errno.DefErr)
	}

	//err = changelogs.New().UpdateGroupInfo(request.GID)
	//if err != nil {
	//	return err
	//}
	//conf, _ := config.GetIMSdkKey()
	//imsdkChat.NewClient(&imsdkChat.Options{
	//	Credentials: imsdk.NewStaticCredentials(conf.AK, conf.SK),
	//}).UpdateName(gid, name)

	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   []string{},
	//	"oldName":  groupInfo.Name,
	//	"newName":  name,
	//	"temId":    "mod-group-name",
	//}
	//actionData := map[string]interface{}{}
	//// todo should use queue and must be send successfully
	//er1 := message.HollowManSendGroupSystemMsg(gid, sysMsgData, actionData)
	//if er1 != nil {
	//	return er1
	//}
	return nil

}

func UpdateAvatar(ctx context.Context, uid string, request UpdateAvatarRequest) error {
	_, err := verifyGroupAndIsAdministrator(uid, request.GroupID, true)
	if err != nil {
		return err
	}
	//var avatarMap map[string]interface{}
	//err = json.Unmarshal([]byte(avatar), &avatar)
	detailDao := detail.New()
	data, _ := detailDao.GetByID(request.GroupID, "_id,avatar")
	if data.Avatar == "" {
		request.IsNotice = 0
	}
	uData := detail.Detail{
		Avatar:    request.Avatar,
		UpdatedAt: funcs.GetMillis(),
	}

	if err = detailDao.UpByID(request.GroupID, uData); err != nil {
		return err
	}
	//err = changelogs.New().UpdateGroupInfo(request.GroupID)
	//if err != nil {
	//	// record error log
	//}
	//conf, _ := config.GetIMSdkKey()
	//imsdkChat.NewClient(&imsdkChat.Options{
	//	Credentials: imsdk.NewStaticCredentials(conf.AK, conf.SK),
	//}).UpdateAvatar(request.GroupID, request.Avatar)

	if request.IsNotice == 0 {
		return nil
	}
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   []string{},
	//	"temId":    "mod-group-avatar",
	//}
	//actionData := map[string]interface{}{}
	//// todo should use queue and must be send successfully
	//
	//if er := message.HollowManSendGroupSystemMsg(request.GroupID, sysMsgData, actionData); er != nil {
	//	return er
	//}
	return nil
}

func UpdateNotice(ctx context.Context, uid string, request UpdateNoticeRequest) error {
	groupInfo, err := verifyGroupAndIsAdministrator(uid, request.GroupID, true)
	if err != nil {
		return err
	}
	noticeId := funcs.SHA1(request.Notice)
	if strings.Trim(request.Notice, "") == "" {
		request.Notice = ""
	}
	uData := map[string]interface{}{
		"notice":      request.Notice,
		"notice_id":   noticeId,
		"update_time": funcs.GetMillis(),
	}
	if err = detail.New().UpMapByID(groupInfo.ID, uData); err != nil {
		return err
	}

	//changelogs.New().UpdateGroupInfo(request.GID)

	if request.Notice == "" {
		return nil
	}
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   []string{},
	//	"temId":    "mod-group-notice",
	//}
	//actionData := map[string]interface{}{}
	//// todo should use queue and must be send successfully
	//err = message.HollowManSendGroupSystemMsg(gid, sysMsgData, actionData)
	//if err != nil {
	//	return err
	//}
	return nil
}

func GetNotice(uid string, request GetNoticeRequest) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	if !members.New().IsExist(uid, request.GroupID) {
		return res, errno.Add("forbidden", errno.FORBIDDEN)
	}
	data, err := detail.New().GetByID(request.GroupID, "_id,notice_id,notice")
	if err != nil {
		return res, err
	}
	res = map[string]interface{}{
		"id":        request.GroupID,
		"notice_id": data.NoticeId,
		"notice":    "",
		"is_update": 0,
	}
	if data.NoticeId != request.NoticeID {
		res["notice"] = data.Notice
		res["notice_id"] = data.NoticeId
		res["is_update"] = 1
	}
	return res, err
}

func UpdateAlias(ctx context.Context, uid string, request UpdateAliasRequest) error {
	groupInfo, err := detail.New().GetByID(request.GroupID, "id,status")
	if err != nil || groupInfo.ID == "" || groupInfo.Status != detail.StatusYes {
		return errno.Add("group-not-exist", errno.DataNotExist)
	}

	memDao := members.New()
	memInfo, err := memDao.GetByUidAndGid(uid, request.GroupID, "id,role")
	if err != nil {
		return errno.Add("forbidden", errno.FORBIDDEN)
	}
	uData := map[string]interface{}{
		"my_alias":    request.Alias,
		"update_time": funcs.GetMillis(),
	}
	if memDao.UpAliasByID(memInfo.ID, uData) {
		changelogs.New().UpdateMemberInfo(request.GroupID, uid)
		return nil
	}

	return errno.Add("fail", errno.DefErr)
}

func GetListByUid(ctx context.Context, uid string, request ListRequest) ([]members.MyGroupRes, error) {
	return members.New().GetMyGroupList(uid, request.GroupIDs)
}

func GetIdListByUid(uid string) ([]string, error) {
	res := make([]string, 0)
	data, err := members.New().GetMyGroupIdList(uid)
	if data == nil {
		return res, err
	}
	for _, v := range data {
		res = append(res, v.ID)
	}
	return res, err
}

func GetMembersInfo(ctx context.Context, uid string, request MembersRequest) ([]members.GroupMembersInfoRes, error) {
	return members.New().GetMembersInfo(request.GroupID, request.ObjUid)
}

func KickOutGroup(ctx context.Context, uid string, request KickOutRequest) error {
	groupInfo, err := detail.New().GetByID(request.GroupID, "id,status,name")
	if err != nil || groupInfo.ID == "" || groupInfo.Status != detail.StatusYes {
		return errno.Add("group-not-exist", errno.DataNotExist)
	}
	memInfo, _ := members.New().GetByUidAndGid(uid, request.GroupID, "id,role")
	if memInfo.Role != members.RoleOwner && memInfo.Role != members.RoleAdministrator {
		return errno.Add("forbid", errno.FORBIDDEN)
	}
	// can not kick out administrator
	allMembers := append(request.ObjUid, uid)
	memList, err := members.New().GetByGidAndUids(request.GroupID, allMembers, "id,role,uid,admin_time")
	if err != nil {
		return errno.Add("wrong-opt", errno.FORBIDDEN)
	}
	//if len(memList) != len(allMembers) {
	//	return errno.Add("wrong-opt", errno.FORBIDDEN)
	//}
	idArr := make([]string, 0)
	memArr := make([]string, 0)
	for _, v := range memList {
		if v.UID == uid {
			continue
		}
		if memInfo.Role == members.RoleAdministrator {
			// Verify whether the requester(self) is an administrator or group owner
			if v.Role == members.RoleOwner || v.Role == members.RoleAdministrator {
				continue
			}
		}
		idArr = append(idArr, v.ID)
		memArr = append(memArr, v.UID)
	}
	if len(memArr) == 0 {
		return nil
	}
	// remove members of group
	// todo should use queue and must be send successfully
	if err = members.New().RemoveMembers(idArr); err != nil {
		fmt.Println("KickOutGroup err:", err)
		return errno.Add("fail", errno.DefErr)
	}
	for _, id := range memArr {
		er1 := changelogs.New().DelMemberInfo(request.GroupID, id)
		if er1 != nil {
			return er1
		}
	}

	groupInfo.Total = groupInfo.Total - len(idArr)
	//_ = generateSnapshot(uid, groupInfo)
	// send system message by hollow man （tmpMsg：kick-out-group-member）
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   memArr,
	//	"temId":    "kick-out-group-member",
	//}
	//actionData := map[string]interface{}{
	//	"name": "update_group_member",
	//}

	detail.New().UpTotal(request.GroupID, -len(idArr))
	defer func() {
		if er := recover(); er != nil {
			//app.Logger().WithField("action", "fail to send msg when inviting users to join group").Error(err)
		}
	}()
	cacheTag := funcs.Md516(request.GroupID)
	redis.Client.SRem(cacheTag, idArr)
	//err = message.HollowManSendGroupSystemMsg(gid, sysMsgData, actionData)
	//if err != nil {
	//	return err
	//}
	//// remove member
	//conf, _ := config.GetIMSdkKey()
	//err = imsdkChat.NewClient(&imsdkChat.Options{
	//	Credentials: imsdk.NewStaticCredentials(conf.AK, conf.SK),
	//}).Remove(gid, memArr)
	chat.RemoveMember(ctx, chat.RemoveMembersParams{
		ChatId: request.GroupID,
		UIds:   request.ObjUid,
	})
	return nil
}

func DisbandGroup(ctx context.Context, uid string, request DisbandRequest) error {
	_, err := verifyGroupAndIsAdministrator(uid, request.GroupID, true)
	if err != nil {
		return err
	}
	if ChangeGroupStatus(uid, request.GroupID, detail.StatusDel) == nil {
		//changelogs.New().UpdateGroupInfo(gid)
		return nil
	}
	return errno.Add("fail", errno.DefErr)
}

func GetRandomGroupOwner(gid string) string {
	memDao := members.New()
	memList, _ := memDao.GetGroupMember(gid)
	adminMinTime := funcs.GetMillis()
	memberMinTime := funcs.GetMillis()
	var adminIdArr []string
	var memberIdArr []string
	for _, v := range memList {
		if v.Role == members.RoleAdministrator {
			if adminMinTime == 0 {
				adminMinTime = v.CreateAt
				adminIdArr = append(adminIdArr, v.UID)
			} else if v.CreateAt <= adminMinTime {
				adminIdArr = append(adminIdArr, v.UID)
				adminMinTime = v.CreateAt
			}
		} else if v.Role == members.RoleCommonMember {
			if adminMinTime == 0 {
				memberMinTime = v.CreateAt
			} else if v.CreateAt <= memberMinTime {
				memberIdArr = append(memberIdArr, v.UID)
				memberMinTime = v.CreateAt
			}
		}
	}
	newOwner := ""
	if adminIdArr != nil {
		adminMemInfo, _ := user2.New().GetNameInfo(adminIdArr)
		adminMemInfoTmp := make(map[string]interface{}, 0)
		var adminMemName []string
		for _, v := range adminMemInfo {
			adminMemInfoTmp[v.Name] = v.ID
			adminMemName = append(adminMemName, v.Name)
			sort.Sort(sort.StringSlice(adminMemName))
		}
		if len(adminMemName) > 0 {
			newOwner = adminMemInfoTmp[adminMemName[0]].(string)
		}
	} else {
		memInfo, _ := user2.New().GetNameInfo(memberIdArr)
		memInfoInfoTmp := make(map[string]interface{}, 0)
		var memName []string
		for _, v := range memInfo {
			memInfoInfoTmp[v.Name] = v.ID
			memName = append(memName, v.Name)
			sort.Sort(sort.StringSlice(memName))
		}
		if len(memName) > 0 {
			newOwner = memInfoInfoTmp[memName[0]].(string)
		}
	}
	return newOwner
}

//func RandomTransferOwner(uid, gid string) error {
//	memDao := members.New()
//	memList, _ := memDao.GetGroupMember(gid)
//	adminMinTime := funcs.GetMillis()
//	memberMinTime := funcs.GetMillis()
//	var adminIdArr []string
//	var memberIdArr []string
//	for _, v := range memList {
//		if v.Role == members.RoleAdministrator {
//			if adminMinTime == 0 {
//				adminMinTime = v.CreateAt
//				adminIdArr = append(adminIdArr, v.UID)
//			} else if v.CreateAt <= adminMinTime {
//				adminIdArr = append(adminIdArr, v.UID)
//				adminMinTime = v.CreateAt
//			}
//		} else if v.Role == members.RoleCommonMember {
//			if adminMinTime == 0 {
//				memberMinTime = v.CreateAt
//			} else if v.CreateAt <= memberMinTime {
//				memberIdArr = append(memberIdArr, v.UID)
//				memberMinTime = v.CreateAt
//			}
//		}
//	}
//	newOwner := ""
//	if adminIdArr != nil {
//		adminMemInfo, _ := info.New().GetNameInfo(adminIdArr)
//		adminMemInfoTmp := make(map[string]interface{}, 0)
//		var adminMemName []string
//		for _, v := range adminMemInfo {
//			adminMemInfoTmp[v.Name] = v.ID
//			adminMemName = append(adminMemName, v.Name)
//			sort.Sort(sort.StringSlice(adminMemName))
//		}
//		newOwner = adminMemInfoTmp[adminMemName[0]].(string)
//
//	} else {
//		memInfo, _ := info.New().GetNameInfo(memberIdArr)
//		memInfoInfoTmp := make(map[string]interface{}, 0)
//		var memName []string
//		for _, v := range memInfo {
//			memInfoInfoTmp[v.Name] = v.ID
//			memName = append(memName, v.Name)
//			sort.Sort(sort.StringSlice(memName))
//		}
//		newOwner = memInfoInfoTmp[memName[0]].(string)
//	}
//	err := memDao.UpdateMemberRole(newOwner, gid, members.RoleOwner)
//	if err != nil {
//		return errno.Add("fail", errno.DefErr)
//	}
//	target := []string{newOwner}
//	sysMsgData := map[string]interface{}{
//		"operator": uid,
//		"target":   target,
//		"temId":    "transfer-owner",
//	}
//	// todo should use queue and must be send successfully
//	err = message.HollowSendOnlyGroupOwnerMsg(gid, sysMsgData, "")
//	if err != nil {
//		return err
//	}
//	return nil
//}

func TransferGroup(ctx context.Context, uid string, request TransferRequest) error {
	_, err := verifyGroupAndIsOwner(uid, request.GroupID, true)
	if err != nil {
		return err
	}
	if uid == request.ObjUId {
		return nil
	}
	memDao := members.New()
	if !memDao.IsExist(request.ObjUId, request.GroupID) {
		return errno.Add("user-not-in-group", errno.DefErr)
	}
	millis := funcs.GetMillis()
	uData := detail.Detail{
		OwnerUid:  request.ObjUId,
		UpdatedAt: millis,
	}
	if err = detail.New().UpByID(request.GroupID, uData); err != nil {
		return errno.Add("failed to update group info", errno.DefErr)
	}

	// new owner
	memUData := members.Members{
		Role:      members.RoleCommonMember,
		UpdatedAt: millis,
	}
	memDao.UpByID(memDao.GetId(request.ObjUId, request.GroupID), memUData)

	// convert old owner to common member
	oldMemUData := members.Members{
		Role:      members.RoleCommonMember,
		UpdatedAt: millis,
	}
	memDao.UpByID(memDao.GetId(uid, request.GroupID), oldMemUData)

	changelogs.New().UpdateManyMemberInfo(request.GroupID, []string{uid, request.ObjUId})
	// send system message by hollow man （tem msgId：transfer-group-owner）
	//target := []string{targetUid}
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   target,
	//	"temId":    "transfer-owner",
	//}
	//// todo should use queue and must be send successfully
	//actionData := map[string]interface{}{}
	//err = message.HollowManSendGroupSystemMsg(gid, sysMsgData, actionData)
	//if err != nil {
	//	return err
	//}
	return nil
}

func ChangeGroupStatus(uid, gid string, status int8) error {
	_, err := verifyGroupAndIsAdministrator(uid, gid, false)
	if err != nil {
		return err
	}
	detailDao := detail.New()
	uData := detail.Detail{
		Status:    status,
		UpdatedAt: funcs.GetMillis(),
	}
	if detailDao.UpByID(gid, uData) == nil {
		return nil
	}
	return errno.Add("fail", errno.DefErr)
}

//func CreateQrCode(uid, gid string) map[string]string {
//	t := funcs.GetMillis()
//	expire, _ := app.Config().GetChildConf("global", "system", "group_qrcode_expire")
//
//	data := qrcode.QrCode{
//		ID:        uuid.New().String(),
//		GroupID:   gid,
//		CreateUID: uid,
//		Expire:    t + int64(expire.(float64))*86400*1000,
//		CreatedAt: t,
//		UpdatedAt: t,
//	}
//	if err := qrcode.New().Add(data); err == nil {
//		qrData := "qr_code=" + data.ID + "&id=" + gid + "&expire=" + strconv.Itoa(int(data.Expire))
//		return map[string]string{"qr": pkg.BuildQrCode(qrData, "group", true)}
//	}
//	return nil
//}

//func JoinV2(uid, gid string, qrCode []byte) error {
//	codeByte, er := qrcode2.VerifyCode(qrCode)
//	if er != nil {
//		return er
//	}
//	return Join(uid, gid, string(codeByte))
//}

func Join(ctx context.Context, uid string, request JoinRequest) error {
	fmt.Println("join group", uid, request.GroupID, request.QrID)
	groupInfo, err := detail.New().GetByID(request.GroupID, "id,status,total")
	if err != nil || groupInfo.ID == "" || groupInfo.Status != detail.StatusYes {
		return errno.Add("group-not-exist", errno.DataNotExist)
	}

	if int(groupInfo.Total) >= detail.TotalMax {
		return errno.Add("group-exceed-limit", MemExceedMax)
	}
	// verify qr code is expire

	memDao := members.New()
	isExist := memDao.IsExist(uid, request.GroupID)
	if isExist {
		return nil
	}
	uInfo, _ := user2.New().GetByID(uid)
	if uInfo.ID == "" {
		return errno.Add("user does not exist", errno.UserNotExist)
	}
	// add member
	myAlias := uInfo.Name
	t := funcs.GetMillis()
	item := members.Members{
		ID:       memDao.GetId(uid, request.GroupID),
		GroupID:  request.GroupID,
		UID:      uid,
		Role:     members.RoleCommonMember,
		JoinType: members.JoinTypeSelf,
		//InviteUID: qrInfo.Uid,
		Status:    members.StatusIng,
		MyAlias:   myAlias,
		CreatedAt: t,
		UpdatedAt: t,
	}
	if ok := memDao.AddOne(item); !ok {
		return errno.Add("fail", errno.DefErr)
	}

	detail.New().UpTotal(request.GroupID, 1)
	// join chat
	err = chat.JoinChat(context.Background(), chat.JoinChatParams{
		ChatId: request.GroupID,
		UIds:   []string{uid},
	})
	if err != nil {
		return err
	}
	// delete snapshot
	//_ = removeSnapshot([]string{uid}, gid)

	defer func() {
		if err := recover(); err != nil {
		}
	}()
	// add change logs
	err = changelogs.New().UpdateMemberInfo(request.GroupID, uid)
	if err != nil {
		return err
	}
	// send system message by hollow man （模板消息：join-group）
	//var target []string
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   target,
	//	"temId":    "join-group",
	//}
	//// todo should use queue and must be send successfully
	//actionData := map[string]interface{}{
	//	"name": "update_group_member",
	//}
	redis.Client.Del(funcs.Md516(request.GroupID))
	//err = message.HollowManSendGroupSystemMsg(gid, sysMsgData, actionData)
	//if err != nil {
	//	return err
	//}
	return nil
}

//func AddAdministrator(uid, targetId, gid string) error {
//	_, err := verifyGroupAndIsOwner(uid, gid, true)
//	if err != nil {
//		return err
//	}
//	memDao := members.New()
//	//adminLen := len(memDao.GetAdministrator(gid))
//	//maxAdminNum, _ := app.Config().GetChildConf("global", "group", "max_administrator_num")
//	//if adminLen >= int(maxAdminNum.(float64)) {
//	//	return errno.Add("exceed-maximum-limit", errno.TimesLimited)
//	//}
//	modRes := memDao.UpdateMemberRole(targetId, gid, members.RoleAdministrator)
//	if modRes != nil {
//		return modRes
//	}
//
//	er1 := changelogs.New().UpdateMemberInfo(gid, targetId)
//	if er1 != nil {
//		return er1
//	}
//	// send system message by hollow man （模板消息：add-administrator）
//	var target []string
//	target = append(target, targetId)
//	sysMsgData := map[string]interface{}{
//		"operator": uid,
//		"target":   target,
//		"temId":    "add-administrator",
//	}
//	// todo should use queue and must be send successfully
//	actionData := map[string]interface{}{}
//	err = message.HollowManSendGroupSystemMsg(gid, sysMsgData, actionData)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func AddAdministrators(ctx context.Context, uid string, request AddAdministratorRequest) ([]members.Members, error) {
	_, err := verifyGroupAndIsOwner(uid, request.GroupID, true)
	var res []members.Members
	if err != nil {
		return res, err
	}
	memDao := members.New()
	targetInfo, _ := memDao.GetByGidAndUids(request.GroupID, request.ObjUid, "_id,uid,role")

	if len(targetInfo) == 0 {
		return res, nil
	}
	roleIds := make([]string, 0)
	for _, v := range targetInfo {
		if v.Role == members.RoleCommonMember {
			roleIds = append(roleIds, v.UID)
		}
	}
	modRes := memDao.UpdateRoleByUIds(roleIds, request.GroupID, members.RoleAdministrator)
	if modRes != nil {
		return res, modRes
	}
	changelogs.New().UpdateManyMemberInfo(request.GroupID, request.ObjUid)
	//send system message by hollow man （：add-administrator）
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"temId":    "add-administrator",
	//}
	//for _, v := range roleIds {
	//	sysMsgData["target"] = []string{v}
	//	// todo should use queue and must be send successfully
	//	actionData := map[string]interface{}{}
	//	err = message.HollowSendMsgInGroup(base.MsgTypeCmd, gid, []string{v}, sysMsgData, actionData)
	//	if err != nil {
	//		app.Logger().WithFields(map[string]interface{}{
	//			"action":    "AddAdministrators",
	//			"targetIds": targetIds,
	//			"uid":       uid,
	//			"err":       err,
	//		}).Error("administrator err")
	//	}
	//}
	res, err = memDao.GetByGidAndUids(request.GroupID, request.ObjUid, "_id,uid,admin_time")
	return res, err
}

func RemoveAdministrator(ctx context.Context, uid string, request RemoveAdministratorRequest) error {
	_, err := verifyGroupAndIsOwner(uid, request.GroupID, true)
	if err != nil {
		return err
	}
	memDao := members.New()
	modRes := memDao.UpdateMemberRole(request.ObjUid[0], request.GroupID, members.RoleCommonMember)
	if modRes != nil {
		return modRes
	}
	changelogs.New().UpdateMemberInfo(request.GroupID, request.ObjUid[0])
	// send system message by hollow man （模板消息：remove-administrator）
	//var target []string
	//target = append(target, targetId)
	//sysMsgData := map[string]interface{}{
	//	"operator": uid,
	//	"target":   target,
	//	"temId":    "remove-administrator",
	//}
	//// todo should use queue and must be send successfully
	//actionData := map[string]interface{}{}
	//err = message.HollowManSendGroupSystemMsg(gid, sysMsgData, actionData)
	return nil
}

func GetGroupsMemberIds(gIds []string) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)
	fields := "uid,gid"
	data, err := members.New().GetGroupsMemberInfo(gIds, fields)
	if data == nil {
		return nil, err
	}
	gRes := make(map[string][]string, 0)
	for _, v := range data {
		gRes[v.GroupID] = append(gRes[v.GroupID], v.UID)
	}
	for _, v := range gIds {
		uIds := gRes[v]
		if gRes[v] == nil {
			uIds = []string{}
		}
		tmp := map[string]interface{}{
			"id":   v,
			"uids": uIds,
		}
		res = append(res, tmp)
	}
	return res, nil
}

//func GetGroupsMemberInfoByIds(gIds []string) ([]map[string]interface{}, error) {
//	res := make([]map[string]interface{}, 0)
//	data, err := members.New().GetMemberInfoByGIds(gIds)
//	if data == nil {
//		return nil, err
//	}
//	gRes := make(map[string][]members.GroupMembersInfoRes, 0)
//	for _, v := range data {
//		gRes[v.GID] = append(gRes[v.GID], v)
//	}
//	for _, v := range gIds {
//		uInfo := gRes[v]
//		if gRes[v] == nil {
//			uInfo = []members.GroupMembersInfoRes{}
//		}
//		tmp := map[string]interface{}{
//			"id":      v,
//			"members": uInfo,
//		}
//		res = append(res, tmp)
//	}
//	return res, nil
//}

//func GetGroupMemberInfoByUIds(ctx context.Context, request GetGroupMemberInfoByUIdsRequest) ([]map[string]interface{}, error) {
//	res := make([]map[string]interface{}, 0)
//	ids, gIds := make([]string, 0), make([]string, 0)
//	for _, item := range request.GroupMember {
//		gIds = append(gIds, item.GroupID)
//		for _, uid := range item.UIDs {
//			ids = append(ids, members.New().GetId(uid, item.GroupID))
//		}
//	}
//	data, err := members.New().GetMemberInfoByIds(ids)
//	if data == nil {
//		return nil, err
//	}
//	gRes := make(map[string][]members.GroupMembersInfoRes, 0)
//	for _, v := range data {
//		gRes[v.GID] = append(gRes[v.GID], v)
//	}
//	for _, v := range gIds {
//		uInfo := gRes[v]
//		if gRes[v] == nil {
//			uInfo = []members.GroupMembersInfoRes{}
//		}
//		tmp := map[string]interface{}{
//			"id":      v,
//			"members": uInfo,
//		}
//		res = append(res, tmp)
//	}
//	return res, nil
//}

func GetGroupName(uid string, ids []string) ([]BaseDetail, error) {
	list, err := detail.New().GetInfoById(ids, "_id,name,avatar")
	result := make([]BaseDetail, 0)
	if err != nil {
		return result, err
	}
	myGroupList, _ := members.New().GetMyGroupByGid(uid, ids, "gid")
	myGroups := make(map[string]struct{})
	for _, m := range myGroupList {
		myGroups[m.GroupID] = struct{}{}
	}
	for _, v := range list {
		item := BaseDetail{
			ID:     v.ID,
			Name:   v.Name,
			Avatar: v.Avatar,
			IsIn:   0,
		}
		if _, ok := myGroups[v.ID]; ok {
			item.IsIn = 1
		}
		result = append(result, item)
	}
	return result, nil
}

func GetAvatar(str string) (avatarInfo map[string]interface{}) {
	avatarInfo = map[string]interface{}{
		"height":    240,
		"weight":    240,
		"bucketId":  "2lqxqjg9lnrt1",
		"file_type": "png",
		"text":      "logo",
	}
	return
}

//func generateSnapshot(uid string, groupInfo detail.Detail) error {
//	snapshotDao := snapshot.New()
//	t := funcs.GetMillis()
//	data := snapshot.Snapshot{
//		ID:        snapshotDao.GetId(uid, groupInfo.ID),
//		Gid:       groupInfo.ID,
//		Total:     groupInfo.Total,
//		Name:      groupInfo.Name,
//		Avatar:    groupInfo.Avatar,
//		QuitTime:  t,
//		UpdatedAt: t,
//		CreatedAt: t,
//	}
//	return snapshotDao.Save(data)
//}

func QuitGroupDelChat(ctx context.Context, uid, gid string) error {
	// clear message
	//if _, err := usermessage.New().DeleteSelfMsgInChat(uid, gid); err != nil {
	//	return err
	//}
	err := changelogs.New().UpdateMemberInfo(gid, uid)
	if err != nil {
		return err
	}
	return nil
	// send message of del chat
	//actionData := map[string]interface{}{}
	//sysMsgData := map[string]interface{}{
	//	"achat_id": gid,
	//}
	//er1 := message2.DelUserMessageByAChatIds(ctx, message2.DelUserMessageByAChatIdsParams{})
	//if er1 != nil {
	//	return er1
	//}
	//return message.HollowSendCustomMsgInGroup(base.MsgTypeDelChat, gid, []string{uid}, sysMsgData, actionData)
}

//func removeSnapshot(uids []string, gid string) error {
//	ids := make([]string, 0)
//	snapshotDao := snapshot.New()
//	for _, uid := range uids {
//		ids = append(ids, snapshotDao.GetId(uid, gid))
//	}
//	return snapshot.New().DeleteByIds(ids)
//}

//func GetSnapshotList(uid string, gIds []string) []snapshot.Snapshot {
//	ids := make([]string, 0)
//	snapshotDao := snapshot.New()
//	for _, gid := range gIds {
//		ids = append(ids, snapshotDao.GetId(uid, gid))
//	}
//	list := snapshotDao.GetByIds(ids)
//	tmp := make(map[string]struct{})
//	diffGIds := make([]string, 0)
//	if len(list) == 0 {
//		diffGIds = gIds
//	} else {
//		for _, i := range list {
//			tmp[i.Gid] = struct{}{}
//		}
//		for _, gid := range gIds {
//			if _, ok := tmp[gid]; !ok {
//				diffGIds = append(diffGIds, gid)
//			}
//		}
//	}
//	if len(diffGIds) > 0 {
//		groupList, _ := detail.New().GetInfoById(diffGIds, "_id,total,name,avatar")
//		for _, g := range groupList {
//			list = append(list, snapshot.Snapshot{
//				Gid:      g.ID,
//				Total:    int(g.Total),
//				Name:     g.Name,
//				Avatar:   g.Avatar,
//				QuitTime: 0,
//			})
//		}
//	}
//
//	return list
//}
