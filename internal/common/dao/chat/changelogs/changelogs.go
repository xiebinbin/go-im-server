package changelogs

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"strconv"
)

type ChangeLogs struct {
	ID        string `bson:"_id" json:"id"`
	ChatID    string `bson:"chat_id" json:"chat_id,omitempty"`
	Uid       string `bson:"uid" json:"uid,omitempty"`
	Status    int8   `bson:"status" json:"status,omitempty"`     // value: 1: normal , 2: del
	Sequence  int64  `bson:"sequence" json:"sequence,omitempty"`
	UniSeqIdx string `bson:"uni_seq_idx" json:"uni_seq_idx,omitempty"` // unique index
	CreatedAt int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt int64  `bson:"update_time" json:"update_time,omitempty"`
}

const (
	UidInfo      = "chat_info"
	StatusNormal = 1
	StatusDel    = 2
)

// chat_id  ,max_sequence

func New() *ChangeLogs {
	return new(ChangeLogs)
}

func (c *ChangeLogs) GetId(uid, gid string) string {
	return funcs.Md516(uid + gid)
}

func (c *ChangeLogs) TableName() string {
	return "chat_change_logs"
}

func (c *ChangeLogs) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(c.TableName(), rpMode)
}

func (c *ChangeLogs) UpdateMemberInfo(chatId, uid string) error {
	data := ChangeLogs{
		ID:     c.GetId(uid, chatId),
		ChatID: chatId,
		Uid:    uid,
		Status: StatusNormal,
	}
	return c.AddLogs(data)
}

func (c *ChangeLogs) UpdateManyMemberInfo(chatId string, uids []string) {
	for _, uid := range uids {
		data := ChangeLogs{
			ID:     c.GetId(uid, chatId),
			ChatID: chatId,
			Uid:    uid,
			Status: StatusNormal,
		}
		c.AddLogs(data)
	}
}

func (c *ChangeLogs) DelMemberInfo(chatId, uid string) error {
	data := ChangeLogs{
		ID:     c.GetId(uid, chatId),
		ChatID: chatId,
		Uid:    uid,
		Status: StatusDel,
	}
	return c.AddLogs(data)
}

func (c *ChangeLogs) UpdateChatInfo(chatId string) error {
	uid := UidInfo
	data := ChangeLogs{
		ID:     c.GetId(uid, chatId),
		ChatID: chatId,
		Uid:    uid,
		Status: StatusNormal,
	}

	return c.AddLogs(data)
}

func (c *ChangeLogs) AddLogs(data ChangeLogs) error {
	return c.Save(data, 0)
}

func (c *ChangeLogs) Save(data ChangeLogs, retryTimes int8) error {
	seq, err := c.GetSeq(data.ChatID)
	if err != nil {
		return err
	}

	data.Sequence = seq
	data.UniSeqIdx = data.ChatID + strconv.FormatInt(seq, 10)
	if err = c.collection().Where(bson.M{"_id": data.ID}).Upsert(&data).Err(); err == nil {
		return nil
	}

	if mongo2.IsDuplicateKeyError(err) {
		err = c.Save(data, retryTimes)
		retryTimes = retryTimes + 1
		if retryTimes > 3 || err == nil {
			return err
		}
	}
	return err
}

func (c *ChangeLogs) GetSeq(chatId string) (int64, error) {
	maxSeq, err := c.GetMaxSequence(chatId)
	if err != nil {
		return 0, err
	}
	for i := 0; i < 3; i++ {
		seq, er := GetSequence(chatId)
		if er != nil {
			return 0, er
		} else if seq > maxSeq {
			return seq, nil
		}
	}
	return 0, errors.New("fail")
}

func (c *ChangeLogs) GetMaxSequence(chatId string) (int64, error) {
	var data ChangeLogs
	where := bson.M{"chat_id": chatId}
	if err := c.collection().Where(where).Sort(bson.M{"sequence": -1}).Limit(1).FindOne(&data); err != nil {
		if !dao.IsNoDocumentErr(err) {
			return 0, err
		}
	}
	return data.Sequence, nil
}

func GetSequence(chatId string) (int64, error) {
	ctx := context.Background()
	upsert := true
	var isReturn options.ReturnDocument = 1
	var res struct {
		ID  string `bson:"_id" json:"id"`
		Seq int64  `bson:"seq" json:"seq"`
	}
	opt := options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &isReturn}
	err := mongo.Database().SetTable("chat_change_logs_sequence").Collection.FindOneAndUpdate(ctx, bson.M{"_id": chatId}, bson.M{"$inc": bson.M{"seq": 1}}, &opt).Decode(&res)
	//fmt.Println("get seq ", res, "  ", err)
	if err != nil {
		return 0, err
	}
	return res.Seq, nil
}

func (c *ChangeLogs) GetLogAmount(gid string, sequence int64) int64 {
	where := bson.M{"chat_id": gid, "sequence": bson.M{"$gt": sequence}}
	return c.collection().Where(where).Count()
}

func (c *ChangeLogs) List(gid string, sequence int64) ([]ChangeLogs, error) {
	var data []ChangeLogs
	where := bson.M{"chat_id": gid, "sequence": bson.M{"$gt": sequence}}
	fields := dao.GetMongoFieldsBsonByString("uid,status,sequence")
	err := c.collection(mongo.SecondaryPreferredMode).Where(where).Sort(bson.M{"sequence": -1}).Fields(fields).FindMany(&data)
	return data, err
}
