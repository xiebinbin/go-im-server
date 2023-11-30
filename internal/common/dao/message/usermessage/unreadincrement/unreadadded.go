package unreadincrement

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/internal/common/dao/message/usermessage/unreadstock"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"strconv"
)

type UnreadIncrement struct {
	ID int64 `bson:"_id" json:"id"` // sequence
}

func New() *UnreadIncrement {
	return new(UnreadIncrement)
}

func (u UnreadIncrement) TableName(uid string) string {
	t := strconv.Itoa(int(funcs.Str2Modulo(uid)))
	return "unread_added_" + t
}

func (u UnreadIncrement) getCollection(uid string, mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(u.TableName(uid), rpMode)
}

func (u UnreadIncrement) Add(uid string, sequence int64) error {
	data := UnreadIncrement{
		ID: sequence,
	}
	_, err := u.getCollection(uid).InsertOne(&data)
	badge := u.GetPushBadge(uid)
	redis.Client.HSet(base.UnreadInfoHash, base.ImUnreadKey+"_"+uid, badge)
	return err
}

func (u UnreadIncrement) Update(uid string, sequence int64) (int64, error) {
	uData := UnreadIncrement{
		ID: sequence,
	}
	res, err := u.getCollection(uid).Where(bson.M{"_id": uid}).UpdateOne(uData)
	badge := u.GetPushBadge(uid)
	redis.Client.HSet(base.UnreadInfoHash, base.ImUnreadKey+"_"+uid, badge)
	return res.MatchedCount, err
}

func (u UnreadIncrement) GetAddedCount(uid string, sequence int) int {
	where := bson.M{"_id": bson.M{"$gt": sequence}}
	count := u.getCollection(uid).Where(where).Count()
	return int(count)
}

func (u UnreadIncrement) GetPushBadge(uid string) int64 {
	badgeInfo := unreadstock.New().GetUsersStockInfo([]string{uid})
	sequence := 0
	num := 0
	if uInfo, ok := badgeInfo[uid]; ok {
		sequence = int(uInfo["sequence"].(int64))
		num = int(uInfo["num"].(uint16))
	}
	addedNum := u.GetAddedCount(uid, sequence)
	badge := num + addedNum
	fmt.Println(uid, "im unreadInfo =====", addedNum, badgeInfo, badge)
	return int64(badge)
}
