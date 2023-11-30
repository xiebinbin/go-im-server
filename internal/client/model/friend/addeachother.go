package friend

import (
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao/friendv2/friend"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
)

func AddEachOther(request AddFriendRequest) error {
	friendDao := friend.New()
	t := funcs.GetMillis()
	addDataBase := map[string]interface{}{
		"status":      friend.StatusYes,
		"agree_time":  t,
		"create_time": t,
		"update_time": t,
	}
	data1 := addDataBase
	data1["uid"] = request.ObjUID
	data1["obj_uid"] = request.UID

	err1 := friendDao.AddFriendByMap(friend.GetId(request.ObjUID, request.UID), data1)
	if err1 != nil && !mongoDriver.IsDuplicateKeyError(err1) {
		return errno.Add("failed to save object friend", errno.DefErr)
	}
	data2 := addDataBase
	data2["uid"] = request.UID
	data2["obj_uid"] = request.ObjUID
	err2 := friendDao.AddFriendByMap(friend.GetId(request.UID, request.ObjUID), data2)

	if err2 != nil && !mongoDriver.IsDuplicateKeyError(err2) {
		return errno.Add("failed to save friend", errno.DefErr)
	}
	return nil
}
