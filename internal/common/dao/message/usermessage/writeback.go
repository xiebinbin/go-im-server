package usermessage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const SequenceHoldPrefix = "sequence-hold-"

func (u *User) fetchBreakMsgIds(uid string, sequences []int64) []MsgIds {
	data := make([]MsgIds, 0)
	where := bson.M{"sequence": bson.M{"$in": sequences}}

	fmt.Println("[fetchBreakMsgIds] sequences: ", sequences)
	err := u.getCollection(uid, mongo.PrimaryMode).Where(where).Fields(dao.GetMongoFieldsBsonByString("_id,sequence,is_read,status")).FindMany(&data)
	if err != nil {
		return nil
	}
	return data
}

func (u *User) saveFakeSequence(uid string, sequence int64) (MsgIds, error) {
	t := funcs.GetMillis()
	data := User{
		ID:        SequenceHoldPrefix + strconv.Itoa(int(sequence)),
		Sequence:  sequence,
		Status:    StatusDelete,
		ChatId:    "",
		IsRead:    IsReadYes,
		CreatedAt: t,
		UpdatedAt: t,
	}
	_, err := u.getCollection(uid).InsertOne(&data)
	fmt.Println("[saveFakeSequence] data: ", data)
	if err == nil {
		return MsgIds{ID: data.ID, Sequence: data.Sequence, IsRead: data.IsRead, Status: data.Status}, nil
	}
	fmt.Println("[saveFakeSequence] error: ", err.Error())

	// unique index error
	var ret MsgIds
	if mongo.IsFiledDuplicateKeyError(err, "sequence") {
		err = u.getCollection(uid, mongo.PrimaryMode).Where(bson.M{"sequence": sequence}).
			Fields(dao.GetMongoFieldsBsonByString("_id,sequence,is_read,status")).FindOne(&ret)
		return ret, err
	}
	return ret, err
}

func (u *User) breakpointFilling(uid string, items []MsgIds, seqMin int64) []MsgIds {
	var seqMax int64
	existSeq := make(map[int64]struct{})
	for i := 0; i < len(items); i++ {
		if items[i].Sequence > seqMax {
			seqMax = items[i].Sequence
		}
		existSeq[items[i].Sequence] = struct{}{}
	}
	fmt.Println("[breakpointFilling] existSeq: ", existSeq)

	seqMin++
	lackSeq := make([]int64, 0)
	for ; seqMin < seqMax; seqMin++ {
		if _, ok := existSeq[seqMin]; !ok {
			lackSeq = append(lackSeq, seqMin)
		}
	}
	fmt.Println("[breakpointFilling] lackSeq: ", lackSeq)

	errorSeq := make([]int64, 0)
	if len(lackSeq) > 0 {
		var err error
		for i := 0; i < len(lackSeq); i++ {
			_, err = u.saveFakeSequence(uid, lackSeq[i])
			if err != nil {
				errorSeq = append(errorSeq, lackSeq[i])
			}
		}
	}
	fmt.Println("[breakpointFilling] errorSeq: ", errorSeq)

	if len(errorSeq) > 0 {
		fills := u.fetchBreakMsgIds(uid, errorSeq)
		if len(fills) > 0 {
			items = append(items, fills...)
		}
	}

	eff := 0
	for i := 0; i < len(items); i++ {
		if !strings.Contains(items[i].ID, SequenceHoldPrefix) {
			items[eff] = items[i]
			eff++
		}
	}
	return items[:eff]
}

/////////////// test part ///////////

func (u *User) TestMsgWriteBack(ctx context.Context) []MsgIds {
	uid := "test"
	u.InitTest(uid)
	return []MsgIds{}

	for i := 1; i <= 20; i++ {
		rand.Seed(int64(i) + time.Now().UnixNano())
		if rand.Intn(5) != 0 {
			u.SaveTest(int64(i), uid, "test-msg-"+strconv.Itoa(i), "", IsReadYes, false)
		}
	}

	var seqStart int64 = 0
	ret := u.GetMsgIds(ctx, uid, seqStart)
	return u.breakpointFilling(uid, ret, seqStart)
}

func (u *User) SaveTest(sequence int64, uid, messageId, chatId string, isRead int8, isRetry bool) (int64, error) {
	t := funcs.GetMillis()
	fmt.Println("save user message-seq:", sequence, "     uid:", uid, "   message-id:", messageId, "    is-read:", isRead)

	data := User{
		ID:        messageId,
		Sequence:  sequence,
		Status:    StatusYes,
		ChatId:    chatId,
		IsRead:    isRead,
		CreatedAt: t,
		UpdatedAt: t,
	}
	_, err := u.getCollection(uid).Where(bson.M{"_id": messageId}).InsertOne(&data)
	if err != nil && mongoDriver.IsDuplicateKeyError(err) {
		fmt.Println("[SaveTest] error: ", err.Error())
		if !strings.Contains(err.Error(), "index: sequence_") {
			res := u.GetById(uid, messageId)
			return res.Sequence, nil
		}
	}
	return sequence, err
}

func (u *User) InitTest(uid string) error {
	var indexes []bson.M
	ctx := context.Background()
	if cursor, err := u.getCollection(uid).Collection.Indexes().List(ctx); err == nil {
		cursor.All(ctx, &indexes)
	} else {
		fmt.Println("[InitTest] indexes error", err)
		return err
	}

	const (
		IndexNotExist = iota + 1
		IndexNotUnique
		IndexIsUnique
	)
	indexStatus := IndexNotExist
	for _, index := range indexes {
		if strings.Contains(index["name"].(string), "sequence_") {
			indexStatus = IndexNotUnique
			if unique, ok := index["unique"]; ok && unique.(bool) {
				indexStatus = IndexIsUnique
			}
		}
	}
	fmt.Println(indexStatus, indexes)

	var err error
	if indexStatus != IndexIsUnique {
		var (
			indexUnique = true
			indexName   = "sequence_-1"
		)
		if indexStatus == IndexNotUnique {
			u.getCollection(uid).Collection.Indexes().DropOne(ctx, indexName)
		}

		opts := options.IndexOptions{Unique: &indexUnique, Name: &indexName}
		indexModel := mongoDriver.IndexModel{
			Keys: bson.D{
				{"sequence", -1},
			},
			Options: &opts,
		}
		_, err = u.getCollection(uid).Collection.Indexes().CreateOne(ctx, indexModel)
		fmt.Println(err)
	}
	return err
}
