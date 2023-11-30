package detail

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"imsdk/internal/common/dao"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"time"
)

type Detail struct {
	ID         string          `bson:"_id" json:"id"`
	ChatId     string          `bson:"chat_id" json:"chat_id,omitempty"`
	FromUID    string          `bson:"from_uid" json:"from_uid"`
	Content    string          `bson:"content" json:"content"`
	Extra      base.JsonString `bson:"extra" json:"extra"`
	Action     interface{}     `bson:"action" json:"action,omitempty"`
	Status     int8            `bson:"status" json:"status"`
	Type       int8            `bson:"type" json:"type"`
	Sequence   int64           `bson:"sequence" json:"sequence,omitempty"`
	CreatedAt  int64           `bson:"create_time" json:"create_time"`
	ReceiveIds []string        `bson:"receive_ids" json:"receive_ids,omitempty"`
	UpdatedAt  int64           `bson:"update_time" json:"update_time"`
}

type GetMessageListResponse struct {
	Id          string `bson:"_id" json:"id"`
	Sequence    int64  `bson:"sequence" json:"sequence"`
	Content     string `bson:"content" json:"content"`
	FromAddress string `bson:"from_uid" json:"from_uid"`
	CreateTime  int64  `bson:"create_time" json:"create_time"`
}

const (
	StatusYes       = 0
	StatusDelete    = 1
	TypeSystem      = 1
	TypeApplication = 2
)

func (d *Detail) TableName() string {
	return "message_detail"
}

func New() *Detail {
	return new(Detail)
}

func (d *Detail) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(d.TableName(), rpMode)
}

func (d *Detail) Init() error {
	ctx := context.Background()
	indexModel := mongoDriver.IndexModel{
		Keys: bson.M{"chat_id": 1},
	}
	_, err := d.collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (d *Detail) Save(params Detail) (int8, error) {
	_, err := d.collection().InsertOne(&params)
	//fmt.Println("DeleteChat AddPrepareMsgByStatus--Save--1", err, mongoDriver.IsDuplicateKeyError(err), params)
	if err != nil && mongoDriver.IsDuplicateKeyError(err) {
		res, _ := d.GetDetailById(params.ID, "_id, status")
		params.Status = res.Status
		er := d.collection().Where(bson.M{"_id": params.ID}).Upsert(&params)
		//fmt.Println("DeleteChat AddPrepareMsgByStatus--Save--2", er.Err(), res, params.Status, res.Status)
		return params.Status, er.Err()
	}
	return params.Status, err
}

func (d *Detail) GetDetails(msgIds []string) []Detail {
	data := make([]Detail, 0)
	// https://studygolang.com/articles/11737
	where := bson.M{"_id": bson.M{"$in": msgIds}}
	err := d.collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	fmt.Println("get message detail list: ", len(msgIds), " ", err)
	return data
}

func (d *Detail) GetListByAMIds(amids []string, number int64) []Detail {
	data := make([]Detail, 0)
	// https://studygolang.com/articles/11737
	where := bson.M{"amid": bson.M{"$in": amids}}
	err := d.collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	fmt.Println("get message detail list: ", len(amids), " ", err)
	return data
}

func (d *Detail) GetCount(where bson.M) int64 {
	count := d.collection().Where(where).Count()
	return count
}

// init
func (d *Detail) GetListByLimit(limit int64, where bson.M) []GetMessageListResponse {
	var data []GetMessageListResponse
	d.collection().Where(where).Sort(bson.M{"sequence": -1}).Limit(limit).FindMany(&data)
	return data
}

func (d *Detail) GetListByPage(limit, offset int64, where bson.M) []Detail {
	var data []Detail
	d.collection().Where(where).Sort(bson.M{"sequence": 1}).Skip(offset).Limit(limit).FindMany(&data)
	return data
}

func (d *Detail) GetDetailByChatIds(chatIds []string, fields string) []Detail {
	data := make([]Detail, 0)
	where := bson.M{"chat_id": bson.M{"$in": chatIds}}
	err := d.collection(mongo.SecondaryPreferredMode).Fields(dao.GetMongoFieldsBsonByString(fields)).Where(where).FindMany(&data)
	fmt.Println("get message detail list: ", len(chatIds), " ", err)
	return data
}

func (d *Detail) GetExceptByChatIds(chatIds, exceptIds []string, fields string) []Detail {
	data := make([]Detail, 0)
	if len(chatIds) == 0 {
		return data
	}
	where := bson.M{
		"$or": []bson.M{
			{"status": bson.M{"$exists": false}},
			{"status": bson.M{"$eq": 0}},
		},
		"$and": []bson.M{
			{"chat_id": bson.M{"$in": chatIds}},
		},
	}
	if len(exceptIds) > 0 {
		where["_id"] = bson.M{"$nin": exceptIds}
	}
	err := d.collection(mongo.SecondaryPreferredMode).Fields(dao.GetMongoFieldsBsonByString(fields)).Where(where).FindMany(&data)
	fmt.Println("get message detail list: ", len(chatIds), " ", err, len(data))
	return data
}

func (d *Detail) GetDetailByChatIdAndSenderId(senderId string, chatIds []string, fields string) []Detail {
	data := make([]Detail, 0)
	where := bson.M{"chat_id": bson.M{"$in": chatIds}, "sender_id": senderId, "status": StatusYes}
	err := d.collection(mongo.SecondaryPreferredMode).
		Fields(dao.GetMongoFieldsBsonByString(fields)).Where(where).FindMany(&data)
	fmt.Println("get message detail list: ", len(chatIds), " ", err)
	return data
}

func (d *Detail) GetExceptByChatIdAndSenderId(senderId string, chatIds, exceptIds []string, fields string) []Detail {
	data := make([]Detail, 0)
	where := bson.M{"chat_id": bson.M{"$in": chatIds}, "sender_id": senderId, "status": StatusYes, "_id": bson.M{"$nin": exceptIds}}
	err := d.collection(mongo.SecondaryPreferredMode).
		Fields(dao.GetMongoFieldsBsonByString(fields)).Where(where).FindMany(&data)
	fmt.Println("get message detail list: ", len(chatIds), " ", err)
	return data
}

func (d *Detail) Delete(uid string, ids []string) (int64, error) {
	uData := map[string]interface{}{
		"status":      StatusDelete,
		"content":     "",
		"update_time": funcs.GetMillis(),
	}
	where := bson.M{"_id": bson.M{"$in": ids}, "sender_id": uid}
	res, err := d.collection().Where(where).UpdateMany(uData)
	return res.ModifiedCount, err
}

func (d *Detail) DeleteMyByChatIds(uid string, chatIds []string) (int64, error) {
	where := bson.M{"chat_id": bson.M{"$in": chatIds}, "sender_id": uid}
	res, err := d.collection().Where(where).Delete()
	return res, err
}

func (d *Detail) DeleteByUID(uid string) (int64, error) {
	where := bson.M{"sender_id": uid}
	res, err := d.collection().Where(where).Delete()
	return res, err
}

func (d *Detail) DeleteByIds(ids []string) (int64, error) {
	where := bson.M{"_id": bson.M{"$in": ids}}
	res, err := d.collection().Where(where).Delete()
	return res, err
}

func (d *Detail) GetMyMsgByChatIds(uid string, chatIds []string) ([]Detail, error) {
	data := make([]Detail, 0)
	where := bson.M{
		"$or": []bson.M{
			{"status": bson.M{"$exists": false}},
			{"status": bson.M{"$eq": 0}},
		},
		"$and": []bson.M{
			{"chat_id": bson.M{"$in": chatIds}, "sender_id": uid},
		},
	}
	fields := dao.GetMongoFieldsBsonByString("_id,sender_id,content")
	err := d.collection(mongo.SecondaryPreferredMode).Where(where).Fields(fields).FindMany(&data)
	return data, err
}

func (d *Detail) GetMyRecentByChatIds(uid string, chatIds, exceptIds []string, limit int64) ([]string, error) {
	data := make([]Detail, 0)
	where := bson.M{
		"$or": []bson.M{
			{"status": bson.M{"$exists": false}},
			{"status": bson.M{"$eq": 0}},
		},
		"$and": []bson.M{
			{"chat_id": bson.M{"$in": chatIds}, "sender_id": uid, "type": bson.M{"$nin": base.NotAddUnreadCountType}},
		},
	}
	if len(exceptIds) > 0 {
		where["_id"] = bson.M{"$nin": exceptIds}
	}
	var res []string
	err := d.collection().Where(where).Sort(bson.M{"create_time": -1}).Limit(limit).FindMany(&data)
	for _, datum := range data {
		res = append(res, datum.ID)
	}
	return res, err
}

func (d *Detail) DeleteSelfByChatIds(uid string, chatIds, exceptIds []string) (int64, error) {
	uData := map[string]interface{}{
		"status":      StatusDelete,
		"content":     "",
		"update_time": funcs.GetMillis(),
	}
	where := bson.M{"chat_id": bson.M{"$in": chatIds}, "sender_id": uid}
	if len(exceptIds) > 0 {
		where["_id"] = bson.M{"$nin": exceptIds}
	}
	res, err := d.collection().Where(where).UpdateMany(uData)
	return res.ModifiedCount, err
}

func (d *Detail) UpdateMapById(id string, upData map[string]interface{}) error {
	_, err := d.collection().Where(bson.M{"_id": id}).UpdateOne(upData)
	return err
}

func (d *Detail) GetDetailById(id, fields string) (Detail, error) {
	cacheTag := "tmm:msg:detail:" + funcs.Md5Str(id+fields)
	res, err := redis.Client.Get(cacheTag).Result()
	if err != nil && err != redis.NilErr {
		return Detail{}, errno.Add("sys err", errno.SysErr)
	}
	var data Detail
	if res != "" {
		err = json.Unmarshal([]byte(res), &data)
	} else {
		where := bson.M{"_id": id}
		err = d.collection(mongo.SecondaryPreferredMode).Fields(dao.GetMongoFieldsBsonByString(fields)).Where(where).FindOne(&data)
		//fmt.Println("get message detail list: ", id, " ", err)
		cache, _ := json.Marshal(data)
		redis.Client.Set(cacheTag, string(cache), time.Second*20)
	}
	return data, err
}

func (d *Detail) UpdateIsRead(ids []string, isRead int8) (int64, error) {
	uData := map[string]interface{}{
		"is_read":   isRead,
		"read_time": funcs.GetMillis(),
	}
	where := bson.M{"_id": bson.M{"$in": ids}}
	res, err := d.collection().Where(where).UpdateMany(&uData)
	return res.ModifiedCount, err
}

func GetSequence(uid string) (int64, error) {
	ctx := context.Background()
	upsert := true
	var isReturn options.ReturnDocument = 1
	var res struct {
		ID  string `bson:"_id" json:"id"`
		Seq int64  `bson:"seq" json:"seq"`
	}
	opt := options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &isReturn}
	err := mongo.Database().SetTable("counters").Collection.FindOneAndUpdate(ctx, bson.M{"_id": uid}, bson.M{"$inc": bson.M{"seq": 1}}, &opt).Decode(&res)
	//fmt.Println("get seq ", res, "  ", err)
	if err != nil {
		return 0, err
	}
	return res.Seq, nil
}
