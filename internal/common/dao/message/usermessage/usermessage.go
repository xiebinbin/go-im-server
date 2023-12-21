package usermessage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"imsdk/internal/common/dao"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"strconv"
	"strings"
)

const (
	StatusBurn    = -1
	StatusYes     = 0
	StatusDelete  = 1
	StatusDisable = 2
	IsReadNo      = 0
	IsReadYes     = 1
	IsFinishNo    = 0
	IsFinishYes   = 1
)

type User struct {
	ID        string `bson:"_id" json:"id"`
	MID       string `bson:"mid" json:"mid"` // message id
	UID       string `bson:"uid" json:"uid"`
	ChatId    string `bson:"chat_id" json:"chat_id"`
	Sequence  int64  `bson:"sequence" json:"sequence"`
	Status    int8   `bson:"status" json:"status"`   //
	IsRead    int8   `bson:"is_read" json:"is_read"` //
	CreatedAt int64  `bson:"create_time" json:"create_time"`
	UpdatedAt int64  `bson:"update_time" json:"update_time"`
}

type MsgIds struct {
	ID       string `bson:"_id" json:"id"` // message id
	Sequence int64  `bson:"sequence" json:"sequence"`
	IsRead   int8   `bson:"is_read" json:"is_read"` //
	Status   int8   `bson:"status" json:"status"`   //
}

func New() *User {
	return new(User)
}

func (u *User) TableName(uid string) string {
	t := strconv.Itoa(int(funcs.Str2Modulo(uid)))
	fmt.Println("TableName uid:", uid, ":uid-s:", t)
	return "message_user_" + t
}

func (u *User) getCollection(uid string, mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(u.TableName(uid), rpMode)
}

// Init : init user message table (create index: sequence)
func (u *User) Init(uid string) error {
	isUnique := true
	opts := options.IndexOptions{Unique: &isUnique}
	indexModel := mongoDriver.IndexModel{
		Keys: bson.D{
			{"mid", -1},
			{"sequence", -1},
			{"chat_id", -1},
			{"uid", -1},
		},
		Options: &opts,
	}

	ctx := context.Background()
	_, err := u.getCollection(uid).Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func GetId(mid, uid string) string {
	return funcs.Md5Str(mid + uid)
}

func (u *User) UpdateMapById(uid, id string, upData map[string]interface{}) error {
	_, err := u.getCollection(uid).Where(bson.M{"_id": id}).UpdateOne(upData)
	return err
}

func GetSequenceNew(uid string) (int64, error) {
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

func (u *User) GetMsgStatus(uid string, ids []string) ([]MsgIds, error) {
	data := make([]MsgIds, 0)
	where := bson.M{"_id": bson.M{"$in": ids}}
	//count := u.getCollection(uid, mongo.SecondaryPreferredMode).Where(where).Count()
	err := u.getCollection(uid, mongo.SecondaryPreferredMode).Where(where).Fields(dao.GetMongoFieldsBsonByString("_id,sequence,is_read,status")).FindMany(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (u *User) GetUnreadIds(uid string, ids []string) ([]string, error) {
	data := make([]MsgIds, 0)
	where := bson.M{"_id": bson.M{"$in": ids}, "is_read": IsReadNo, "uid": uid}
	//count := u.getCollection(uid, mongo.SecondaryPreferredMode).Where(where).Count()
	err := u.getCollection(uid, mongo.SecondaryPreferredMode).Where(where).Fields(dao.GetMongoFieldsBsonByString("_id,sequence,is_read,status")).FindMany(&data)
	var res []string
	if err != nil || len(data) == 0 {
		return res, nil
	}
	for _, datum := range data {
		res = append(res, datum.ID)
	}
	return res, nil
}

func (u *User) GetById(uid, id string) MsgIds {
	var data MsgIds
	err := u.getCollection(uid, mongo.SecondaryPreferredMode).Fields(dao.GetMongoFieldsBsonByString("_id,sequence,is_read,status")).FindByID(id, &data)
	if err != nil {
		return MsgIds{}
	}
	return data
}

func (u *User) GetMsgIds(ctx context.Context, uid string, sequence int64) []MsgIds {
	data := make([]MsgIds, 0)
	where := bson.M{"sequence": bson.M{"$gt": sequence}, "uid": uid}
	logCtx := log.WithFields(ctx, map[string]string{"action": "GetMsgIds"})
	log.Logger().Info(logCtx, "where: ", where, ", uid: ", uid, " sequence: ", sequence)
	err := u.getCollection(uid, mongo.SecondaryPreferredMode).Where(where).Fields(dao.GetMongoFieldsBsonByString("_id,sequence,mid,is_read,status")).FindMany(&data)
	if err != nil {
		log.Logger().Info(logCtx, "failed to query ,err: ", err)
		return nil
	}
	log.Logger().Info(logCtx, "where: ", where, ", uid: ", uid, " sequence: ", sequence, " resp message id amount :", len(data))
	return data
}

func (u *User) GetMsgIdsV2(ctx context.Context, uid string, sequence int64) []MsgIds {
	data := make([]MsgIds, 0)
	where := bson.M{"sequence": bson.M{"$gt": sequence}, "uid": uid}
	logCtx := log.WithFields(ctx, map[string]string{"action": "GetMsgIds"})
	//log.Logger().Info(logCtx, "where: ", where, ", uid: ", uid, " sequence: ", sequence)
	field := dao.GetMongoFieldsBsonByString("_id,mid,sequence,is_read,status")
	err := u.getCollection(uid, mongo.SecondaryPreferredMode).Where(where).Fields(field).Sort(bson.M{"sequence": 1}).Limit(base.MsgPageRowLimit).FindMany(&data)
	if err != nil {
		log.Logger().Error(logCtx, "failed to query ,err: ", err)
		return nil
	}
	//log.Logger().Info(logCtx, "where: ", where, ", uid: ", uid, " sequence: ", sequence, " resp message id amount :", len(data))
	return data
}

func (u *User) Delete(uid string, ids []string) (int64, error) {
	uData := map[string]interface{}{
		"status": StatusDelete,
	}
	where := bson.M{"uid": uid, "mid": bson.M{"$in": ids}}
	res, err := u.getCollection(uid).Where(where).UpdateMany(&uData)
	return res.ModifiedCount, err
}

func (u *User) DeleteSelfAll(uid string) (int64, error) {
	uData := map[string]interface{}{
		"status": StatusDelete,
	}
	res, err := u.getCollection(uid).Where(bson.M{"uid": uid}).UpdateMany(uData)
	return res.ModifiedCount, err
}

func (u *User) DeleteSelfMsgInChat(uid, chatId string) (int64, error) {
	uData := map[string]interface{}{
		"status":      StatusDelete,
		"update_time": funcs.GetMillis(),
	}
	res, err := u.getCollection(uid).Where(bson.M{"chat_id": chatId, "uid": uid}).UpdateMany(uData)
	return res.ModifiedCount, err
}

func (u *User) DeleteSelfMsgInChatIds(uid string, chatIds, exceptIds []string) (int64, error) {
	uData := map[string]interface{}{
		"status":      StatusDelete,
		"update_time": funcs.GetMillis(),
	}
	where := bson.M{"chat_id": bson.M{"$in": chatIds}, "uid": uid}
	if len(exceptIds) > 0 {
		where["_id"] = bson.M{"$nin": exceptIds}
	}
	res, err := u.getCollection(uid).Where(where).UpdateMany(uData)
	return res.ModifiedCount, err
}

func (u *User) SetMsgDisable(uid, msgId string, status int) (int64, error) {
	uData := map[string]interface{}{
		"status": status,
	}
	res, err := u.getCollection(uid).Where(bson.M{"mid": msgId}).UpdateOne(uData)
	return res.ModifiedCount, err
}

func (u *User) UpdateIsReadAll(uid string, isRead int8) (int64, error) {
	uData := map[string]interface{}{
		"is_read": isRead,
	}
	res, err := u.getCollection(uid).Where(bson.M{"uid": uid}).UpdateMany(&uData)
	return res.ModifiedCount, err
}

func (u *User) UpdateIsRead(uid string, ids []string, isRead int8) (int64, error) {
	uData := map[string]interface{}{
		"is_read": isRead,
	}
	where := bson.M{"mid": bson.M{"$in": ids}}
	res, err := u.getCollection(uid).Where(where).UpdateMany(&uData)
	return res.ModifiedCount, err
}

func (u *User) GetMaxSeq(uid string) (int64, error) {
	var data User
	err := u.getCollection(uid).Where(bson.M{"uid": uid}).Sort(bson.M{"sequence": -1}).FindOne(&data)
	return data.Sequence, err
}

func (u *User) UpdateSeq(uid string, seq int64) {
	data := bson.M{"$set": bson.M{"_id": uid, "seq": seq}}
	mongo.Database().SetTable("counters").Where(bson.M{"_id": uid}).UpsertByBson(data)
}

func (u *User) Count(uid string) int64 {
	return u.getCollection(uid).Where(bson.M{"uid": uid}).Count()
}

func (u *User) DropTable(uid string) error {
	return mongo.Database().Where(bson.M{"uid": uid}).Database.Collection(u.TableName(uid)).Drop(context.Background())
}

func (u *User) Add(uid, mid, chatId string, isRead int8) (int64, error) {
	t := funcs.GetMillis()
	sequence, err := GetSequenceNew(uid)
	fmt.Println("usermessage-seq:", sequence, "     ", uid, "   ", mid, "    ", isRead)
	if err != nil {
		return 0, err
	}
	id := GetId(mid, uid)
	data := User{
		ID:        id,
		MID:       mid,
		UID:       uid,
		Sequence:  sequence,
		Status:    StatusYes,
		ChatId:    chatId,
		IsRead:    isRead,
		CreatedAt: t,
		UpdatedAt: t,
	}
	_, er := u.getCollection(uid).InsertOne(&data)
	return sequence, er
}

func (u *User) Save(uid, mid, chatId string, isRead, status int8) (int64, error) {
	t := funcs.GetMillis()
	sequence, err := GetSequenceNew(uid)
	fmt.Println("AddPrepareMsgByStatus save user message-seq:", sequence, "     ", uid, "   ", mid, "    ", isRead, status)
	if err != nil {
		return 0, err
	}
	id := GetId(mid, uid)
	data := User{
		ID:        id,
		MID:       mid,
		UID:       uid,
		Sequence:  sequence,
		Status:    status,
		ChatId:    chatId,
		IsRead:    isRead,
		CreatedAt: t,
		UpdatedAt: t,
	}
	_, er := u.getCollection(uid).Where(bson.M{"_id": id}).InsertOne(&data)
	if er != nil && mongoDriver.IsDuplicateKeyError(er) {
		res := u.GetById(uid, mid)
		return res.Sequence, nil
	}
	return sequence, er
}

func (u *User) SaveNew(uid, messageId, chatId string, isRead int8, isRetry bool) (int64, error) {
	t := funcs.GetMillis()
	sequence, err := GetSequenceNew(uid)
	fmt.Println("save user message-seq:", sequence, "     ", uid, "   ", messageId, "    ", isRead)
	if err != nil {
		return 0, err
	}
	id := GetId(messageId, uid)
	data := User{
		ID:        id,
		MID:       messageId,
		Sequence:  sequence,
		Status:    StatusYes,
		ChatId:    chatId,
		IsRead:    isRead,
		CreatedAt: t,
		UpdatedAt: t,
	}
	_, er := u.getCollection(uid).Where(bson.M{"_id": id}).InsertOne(&data)
	if er != nil && mongoDriver.IsDuplicateKeyError(er) {
		if strings.Contains(err.Error(), "index: sequence_") {
			if isRetry {
				return sequence, er
			}
			return u.SaveNew(uid, messageId, chatId, isRead, true)
		} else {
			res := u.GetById(uid, messageId)
			return res.Sequence, nil
		}
	}
	return sequence, er
}
