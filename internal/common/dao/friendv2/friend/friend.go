package friend

import (
	"context"
	"imsdk/internal/common/dao"
	"imsdk/pkg/app"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type Friend struct {
	ID          string `bson:"_id" json:"id"`
	UId         string `bson:"uid" json:"uid,omitempty"`
	ObjUId      string `bson:"obj_uid" json:"obj_uid,omitempty"`
	RemarkIndex string `bson:"remark_index" json:"remark_index,omitempty"`
	Remark      string `bson:"remark" json:"remark,omitempty"`
	Status      int8   `bson:"status" json:"status,omitempty"`
	IsStar      int8   `bson:"is_star" json:"is_star,omitempty"`
	Rank        int8   `bson:"rank" json:"rank,omitempty"`
	ApplyUID    string `bson:"apply_uid" json:"apply_uid,omitempty"`
	AgreeAt     int64  `bson:"agree_time" json:"agree_time,omitempty"`
	CreatedAt   int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt   int64  `bson:"update_time" json:"update_time,omitempty"`
}

const (
	StatusYes                                 = 1
	StatusDeleted                             = -1
	StatusLocked                              = 0
	STATUS_DELETED, STATUS_LOCKED, STATUS_YES = -1, 0, 1

	OtherSideNotFriend = 100000
	MemberNotFriend    = 100001
	EachOtherStrangers = 100002
)

func New() *Friend {
	return new(Friend)
}

func (f *Friend) TableName() string {
	return "user_friend"
}
func (f *Friend) GetID(uid, objId string) string {
	return GetId(uid, objId)
}
func GetId(uid, objId string) string {
	return funcs.Md516(uid + objId)
}

func (f *Friend) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(f.TableName(), rpMode)
}

func (f *Friend) Init() error {
	ctx := context.Background()
	indexModel := mongoDriver.IndexModel{
		Keys: bson.M{"uid": 1},
	}
	_, err := f.Collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (f *Friend) GetFriendIds(uid string) ([]string, error) {
	var data []Friend
	where := bson.M{"uid": uid, "status": StatusYes}
	f.Collection().Where(where).Fields(bson.M{"obj_uid": 1}).FindMany(&data)
	res := make([]string, 0)
	if data == nil {
		return res, nil
	}
	for _, v := range data {
		res = append(res, v.ObjUId)
	}
	return res, nil
}

func (f *Friend) GetFriendInfos(uid string, objUIds []string) ([]string, map[string]string, map[string]string, error) {
	var data []Friend
	where := bson.M{"uid": uid, "status": StatusYes}
	if len(objUIds) != 0 {
		where["obj_uid"] = bson.M{"$in": objUIds}
	}
	err := f.Collection().Where(where).Fields(bson.M{"obj_uid": 1}).Sort(bson.M{"update_time": -1}).FindMany(&data)
	if err != nil {
		return nil, nil, nil, err
	}
	res := make([]string, 0)
	remarkInfo := make(map[string]string)
	remarkIndexInfo := make(map[string]string)
	if data == nil {
		return res, remarkInfo, remarkIndexInfo, nil
	}
	for _, v := range data {
		res = append(res, v.ObjUId)
		remarkInfo[v.ObjUId] = v.Remark
		remarkIndexInfo[v.ObjUId] = v.RemarkIndex
	}
	return res, remarkInfo, remarkIndexInfo, nil
}

func (f *Friend) GetMineByFriendIds(uid string, ObjUIds []string) []string {
	data := make([]Friend, 0)
	ids := make([]string, 0)
	for _, v := range ObjUIds {
		ids = append(ids, GetId(v, uid))
	}
	where := bson.M{"_id": bson.M{"$in": ids}, "status": StatusYes}
	f.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	res := make([]string, 0)
	if len(data) == 0 {
		return res
	}
	for _, v := range data {
		res = append(res, v.UId)
	}
	return res
}

func (f *Friend) GetFriendsInfo(uid string, ObjUIds []string) (data []Friend, err error) {
	ids := make([]string, 0)
	for _, v := range ObjUIds {
		ids = append(ids, GetId(uid, v))
	}
	where := bson.M{"_id": bson.M{"$in": ids}, "status": StatusYes}
	f.Collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	return
}

func (f *Friend) GetInfoByID(id string) (data Friend, err error) {
	err = f.Collection(mongo.SecondaryPreferredMode).Where(bson.M{"uid": id}).FindOne(&data)
	return
}

func (f *Friend) GetByID(id string) (data Friend, err error) {
	err = f.Collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": id}).FindOne(&data)
	return
}

func (f *Friend) GetByIds(Ids []string) ([]Friend, error) {
	var res []Friend
	where := bson.M{"_id": bson.M{"$in": Ids}}
	err := f.Collection().Where(where).FindMany(&res)
	return res, err
}

func (f *Friend) GetNormalData(Ids []string) ([]Friend, error) {
	var res []Friend
	where := bson.M{"_id": bson.M{"$in": Ids}, "status": StatusYes}
	err := f.Collection().Where(where).FindMany(&res)
	return res, err
}

func (f *Friend) AddFriendByMap(id string, data map[string]interface{}) error {
	return f.Collection().Where(bson.M{"_id": id}).Upsert(&data).Err()
}

func (f *Friend) AddFriend(data Friend) error {
	return f.Collection().Where(bson.M{"_id": data.ID}).Upsert(&data).Err()
}

func (f *Friend) AddMany(data []Friend) ([]interface{}, error) {
	res, err := f.Collection().InsertMany(data)
	if err != nil {
		return []interface{}{}, err
	}
	return res.InsertedIDs, err
}

func (f *Friend) DelFriend(uid, objUId string) error {
	uData := bson.M{"status": StatusDeleted}
	where := bson.M{"_id": GetId(uid, objUId)}
	_, err := f.Collection().Where(where).UpdateOne(&uData)
	return err
}

// DelFriendsUnilateral 单方面
func (f *Friend) DelFriendsUnilateral(uid string, objUIds []string) error {
	uData := bson.M{"status": StatusDeleted}
	if len(objUIds) == 0 {
		return nil
	}
	var ids []string
	for _, i2 := range objUIds {
		ids = append(ids, GetId(uid, i2))
	}
	where := bson.M{"_id": bson.M{"$in": ids}}
	_, err := f.Collection().Where(where).UpdateOne(&uData)
	return err
}

func (f *Friend) DelAllUnilateral(uid string) error {
	uData := bson.M{"status": StatusDeleted}
	where := bson.M{"uid": uid}
	_, err := f.Collection().Where(where).UpdateOne(&uData)
	return err
}

// DelFriendsBilateral 双方面
func (f *Friend) DelFriendsBilateral(uid string, objUIds []string) error {
	uData := bson.M{"status": StatusDeleted}
	if len(objUIds) == 0 {
		return nil
	}
	var ids []string
	for _, i2 := range objUIds {
		ids = append(ids, GetId(uid, i2), GetId(i2, uid))
	}
	where := bson.M{"_id": bson.M{"$in": ids}}
	_, err := f.Collection().Where(where).UpdateOne(&uData)
	return err
}

func (f *Friend) DelAllBilateral(uid string) error {
	uData := bson.M{"status": StatusDeleted}
	where := bson.M{
		"$or": []bson.M{
			{"uid": uid},
			{"ogj_uid": uid},
		},
	}
	_, err := f.Collection().Where(where).UpdateOne(&uData)
	return err
}

func (f *Friend) GetOtherSideRelationInfo(uid string, ObjUIds []string, fields string) (data []Friend, err error) {
	ids := make([]string, 0)
	for _, v := range ObjUIds {
		ids = append(ids, GetId(v, uid))
	}

	where := bson.M{"_id": bson.M{"$in": ids}, "status": StatusYes}
	f.Collection(mongo.SecondaryPreferredMode).Where(where).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return
}

func (f *Friend) GetRelationInfo(uid string, ids []string) ([]string, []string) {
	data, _ := f.GetFriendsInfo(uid, ids)
	otherSideData, _ := f.GetOtherSideRelationInfo(uid, ids, "uid,obj_uid")
	var friendIds []string
	var otherSideIds []string
	tmpInfo := make(map[string]Friend)
	for _, v := range data {
		if v != (Friend{}) {
			friendIds = append(friendIds, v.ObjUId)
		}
		tmpInfo[v.ObjUId] = v
	}
	for _, v := range otherSideData {
		otherSideIds = append(otherSideIds, v.UId)
	}
	return friendIds, otherSideIds
}

func (f *Friend) Add(data Friend) error {
	_, err := f.Collection().InsertOne(&data)
	return err
}

func (f *Friend) UpdateRemark(uid, objUId string, upData map[string]interface{}) (int64, error) {
	id := GetId(uid, objUId)
	where := bson.M{"_id": id}
	res, err := f.Collection().Where(where).UpdateOne(&upData)
	return res.MatchedCount, err
}

func GetPublicRelationUIds() []string {
	var PubUserConf struct {
		PubUid []string `toml:"pub_relation_ids"`
	}
	app.Config().Bind("global", "public_user", &PubUserConf)
	return PubUserConf.PubUid
}
