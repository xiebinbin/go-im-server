package apply

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"strconv"
)

type Apply struct {
	ID        string `bson:"_id" json:"id,omitempty"`
	UId       string `bson:"uid" json:"uid,omitempty"`
	ObjUId    string `bson:"obj_uid" json:"obj_uid,omitempty"`
	Remark    string `bson:"remark" json:"remark,omitempty"`
	Status    int8   `bson:"status" json:"status,omitempty"`
	IsRead    int8   `bson:"is_read" json:"is_read,omitempty"`
	CreatedAt int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt int64  `bson:"update_time" json:"update_time,omitempty"`
}

const (
	StatusIng               = 0
	StatusPass              = 1
	IsReadNo                = 0
	IsReadYes               = 1
	FromWayUnknown     int8 = 0
	FromWaySearchTmmId int8 = 1
	FromWaySearchPhone int8 = 2
	FromWayGroup       int8 = 3
	FromWayScanQrCode  int8 = 4
	FromWayCard        int8 = 5
	FromWayContact     int8 = 6
	FromWayNearByUser  int8 = 7
	FromWayShake       int8 = 8
	FromWayNearby      int8 = 9
	FromWayMoments     int8 = 10
	FromWayOfficial    int8 = 11
)

func New() *Apply {
	return new(Apply)
}

func (a Apply) TableName() string {
	return "user_friend_apply"
}

func GetId(uid, objId string) string {
	return funcs.Md516(uid + objId + strconv.FormatInt(funcs.GetTimeSecs(), 10))
}

func (a Apply) Collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(a.TableName(), rpMode)
}

func (a Apply) Init() error {
	ctx := context.Background()
	indexModel := mongoDriver.IndexModel{
		Keys: bson.M{"obj_uid": 1, "uid": 1},
	}
	_, err := a.Collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (a Apply) AddApply(data Apply) error {
	t := funcs.GetMillis()
	addData := map[string]interface{}{
		"_id":         GetId(data.UId, data.ObjUId),
		"uid":         data.UId,
		"obj_uid":     data.ObjUId,
		"remark":      data.Remark,
		"status":      StatusIng,
		"is_read":     IsReadNo,
		"create_time": t,
		"update_time": t,
	}
	_, err := a.Collection().InsertOne(&addData)
	return err
}

func (a Apply) Add(data Apply) error {
	_, err := a.Collection().InsertOne(&data)
	return err
}

func (a Apply) Delete(id string) error {
	_, err := a.Collection().Where(bson.M{"_id": id}).Delete()
	return err
}

func (a Apply) GetApplyLists(uid string) ([]Apply, error) {
	var data []Apply
	fields := "_id,uid,remark,status"
	where := bson.M{"obj_uid": uid}
	err := a.Collection().Fields(dao.GetMongoFieldsBsonByString(fields)).Where(where).FindMany(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (a Apply) GetApplyInfo(address, objAddress string) (Apply, error) {
	where := bson.M{"obj_uid": objAddress, "uid": address}
	var data Apply
	err := a.Collection().Where(where).Sort(bson.M{"create_time": -1}).FindOne(&data)
	return data, err
}

func (a Apply) GetNotAgreeInfo(uid, objUID string) (Apply, error) {
	where := bson.M{"obj_uid": objUID, "uid": uid, "status": StatusIng}
	var data Apply
	err := a.Collection(mongo.SecondaryPreferredMode).Where(where).Sort(bson.M{"create_time": -1}).FindOne(&data)
	return data, err
}

func (a Apply) GetInfoByID(id string) (data Apply, err error) {
	err = a.Collection().FindByID(id, &data)
	return
}

func (a Apply) UpdateInfoById(id string, uData map[string]interface{}) error {
	_, err := a.Collection().UpByID(id, uData)
	return err
}

func (a Apply) UpdateInfoByIds(ids []string, uData map[string]interface{}) error {
	where := bson.M{"_id": bson.M{"$in": ids}}
	_, err := a.Collection().Where(where).UpdateMany(uData)
	return err
}
