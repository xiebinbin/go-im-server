package block

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
)

type Block struct {
	ID        string `bson:"_id" json:"id"`
	UId       string `bson:"uid" json:"uid,omitempty"`
	ObjUId    string `bson:"obj_uid" json:"obj_uid,omitempty"`
	Status    int8   `bson:"status" json:"status,omitempty"`
	CreatedAt int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt int64  `bson:"update_time" json:"update_time,omitempty"`
}

const (
	StatusBlock             = 1
	StatusNormal            = 2
	ErrCodeMineBlocked      = 100021
	ErrCodeBlockedMine      = 100022
	ErrCodeEachOtherBlocked = 100023
)

func New() *Block {
	return new(Block)
}

func (b Block) TableName() string {
	return "user_block"
}

func (b *Block) GetID(uid, objId string) string {
	return GetId(uid, objId)
}

func GetId(uid, objId string) string {
	return funcs.Md516(uid + objId)
}

func (b *Block) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(b.TableName(), rpMode)
}

func (b *Block) Init() error {
	ctx := context.Background()
	indexModel := mongoDriver.IndexModel{
		Keys: bson.M{"uid": 1},
	}
	_, err := b.Collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (b *Block) GetMineBlockIds(uid string) ([]string, error) {
	var data []Block
	where := bson.M{"uid": uid}
	b.Collection(mongo.SecondaryPreferredMode).Where(where).Fields(bson.M{"obj_uid": 1}).FindMany(&data)
	res := make([]string, 0)
	if data == nil {
		return res, nil
	}
	for _, v := range data {
		res = append(res, v.ObjUId)
	}
	return res, nil
}

func (b *Block) GetBlockMineIds(uid string) ([]string, error) {
	var data []Block
	where := bson.M{"obj_uid": uid}
	b.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	res := make([]string, 0)
	if len(data) > 0 {
		for _, v := range data {
			res = append(res, v.UId)
		}
	}
	return res, nil
}

func (b *Block) SaveBlock(data Block) error {
	return b.Collection().Where(bson.M{"_id": data.ID}).Upsert(&data).Err()
}

func (b *Block) CancelBlock(data Block) error {
	_, err := b.Collection().Where(bson.M{"_id": data.ID}).Delete()
	return err
}

func (b *Block) SaveBlocks(data []Block) (interface{}, error) {
	res, err := b.Collection().UpdateOrInsert(data)
	return res.UpsertedCount, err
}

func (b *Block) GetSuccessIDs(uid string, objUIds []string) (interface{}, error) {
	var res []Block
	where := bson.M{"uid": uid, "obj_uid": bson.M{"$in": objUIds}}
	err := b.Collection().Where(where).FindMany(&res)
	return res, err
}

func (b *Block) DelBlocks(uid string, objUIds []string) (interface{}, error) {
	where := bson.M{"uid": uid, "obj_uid": bson.M{"$in": objUIds}}
	res, err := b.Collection().Where(where).Delete()
	return res, err
}

func (b *Block) IsBlockedMine(uid, objUId string) (bool, error) {
	var res Block
	id := GetId(objUId, uid)
	where := bson.M{"_id": id}
	err := b.Collection().Where(where).FindOne(&res)
	if mongoDriver.IsNetworkError(err) {
		return false, errno.Add("net-err", errno.SysErr)
	}
	if err == mongoDriver.ErrNoDocuments {
		return false, nil
	}
	return true, err
}

func (b *Block) IsMineBlocked(uid, objUId string) (bool, error) {
	var res Block
	id := GetId(uid, objUId)
	where := bson.M{"_id": id}
	err := b.Collection().Where(where).FindOne(&res)
	if mongoDriver.IsNetworkError(err) {
		return false, errno.Add("net-err", errno.SysErr)
	}
	if err == mongoDriver.ErrNoDocuments {
		return false, nil
	}
	return true, err
}

func (b *Block) GetMineBlockedIds(uid string, ObjUIds []string) ([]string, error) {
	ids := make([]string, 0)
	for _, v := range ObjUIds {
		ids = append(ids, GetId(uid, v))
	}
	var data []Block
	where := bson.M{"_id": bson.M{"$in": ids}}
	err := b.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	var res []string
	if len(data) > 0 {
		for _, v := range data {
			res = append(res, v.ObjUId)
		}
	}
	return res, err
}

func (b *Block) GetBlockedMineIds(uid string, ObjUIds []string) ([]string, error) {
	ids := make([]string, 0)
	for _, v := range ObjUIds {
		ids = append(ids, GetId(v, uid))
	}
	var data []Block
	where := bson.M{"_id": bson.M{"$in": ids}}
	err := b.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	var res []string
	if len(data) > 0 {
		for _, v := range data {
			res = append(res, v.UId)
		}
	}
	return res, err
}

func (b *Block) GetByIds(Ids []string) ([]Block, error) {
	var res []Block
	where := bson.M{"_id": bson.M{"$in": Ids}}
	err := b.Collection().Where(where).FindMany(&res)
	return res, err
}
