package unreadstock

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/redis"
)

// client report
type UnreadStock struct {
	ID        string `bson:"_id" json:"id"` // uid
	AKey      string `bson:"akey" json:"akey,omitempty"`
	Sequence  int64  `bson:"sequence" json:"sequence"` // user's latest sequence
	Num       uint16 `bson:"num" json:"num"`           // the number of user unread
	CreatedAt int64  `bson:"create_time" json:"create_time"`
	UpdatedAt int64  `bson:"update_time" json:"update_time"`
}

func New() *UnreadStock {
	return new(UnreadStock)
}

func (u UnreadStock) TableName() string {
	return "unread_stock"
}

func (u *UnreadStock) WithAKey(aKey string) *UnreadStock {
	u.AKey = aKey
	return u
}

func (u UnreadStock) collection() *mongo.CollectionInfo {
	return mongo.Database().SetTable(u.TableName())
}

func (u UnreadStock) Add(uid string, sequence int64, num uint16) bool {
	data := UnreadStock{
		ID:       uid,
		Sequence: sequence,
		Num:      num,
	}
	_, err := u.collection().InsertOne(&data)
	if err != nil {
		// 根据错误判断，是否是已经存在的数据导致错误的，如果是则返回成功，不是则响应失败
		return true
	}
	redis.Client.HSet(base.UnreadInfoHash, base.ImUnreadKey+"_"+uid, num)
	return true
}

func (u UnreadStock) Update(uid string, sequence int64, num uint16) (int64, error) {
	uData := map[string]interface{}{
		"sequence": sequence,
		"num":      num,
	}
	res, err := u.collection().Where(bson.M{"_id": uid}).UpdateOne(uData)
	if res.MatchedCount == 0 {
		addRes := u.Add(uid, sequence, num)
		if !addRes {
			return 0, nil
		}
	}
	redis.Client.HSet(base.UnreadInfoHash, base.ImUnreadKey+"_"+uid, num)
	return res.ModifiedCount, err
}

func (u UnreadStock) GetDetails(uIds []string) []UnreadStock {
	data := make([]UnreadStock, 0)
	where := bson.M{"_id": bson.M{"$in": uIds}}
	u.collection().Where(where).FindMany(&data)
	return data
}

func (u UnreadStock) GetNewestStockInfo(uid string) UnreadStock {
	where := bson.M{"_id": uid}
	sort := bson.M{"sequence": -1}
	var data UnreadStock
	u.collection().Where(where).Sort(sort).FindOne(&data)
	return data
}

func (u UnreadStock) GetUsersStockInfo(ids []string) map[string]map[string]interface{} {
	var data []UnreadStock
	where := bson.M{"_id": bson.M{"$in": ids}}
	u.collection().Where(where).FindMany(&data)
	var res = make(map[string]map[string]interface{}, 0)
	if len(data) == 0 {
		return res
	}
	for _, v := range data {
		res[v.ID] = map[string]interface{}{
			"num":      v.Num,
			"sequence": v.Sequence,
		}
	}
	return res
}
