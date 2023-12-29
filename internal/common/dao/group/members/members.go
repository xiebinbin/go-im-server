package members

import (
	"context"
	"fmt"
	"imsdk/internal/common/dao"
	"imsdk/internal/common/dao/group/detail"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type Members struct {
	ID        string `bson:"_id" json:"id"`
	GroupID   string `bson:"gid" json:"gid,omitempty"`
	UID       string `bson:"uid" json:"uid,omitempty"`
	EncPri    string `bson:"enc_pri" json:"enc_pri,omitempty"`
	EncKey    string `bson:"enc_key" json:"enc_key,omitempty"`
	InviteUID string `bson:"invite_uid" json:"invite_uid,omitempty"`
	Role      uint8  `bson:"role" json:"role,omitempty"`
	JoinType  uint8  `bson:"join_type" json:"join_type,omitempty"`
	MyAlias   string `bson:"my_alias" json:"my_alias,omitempty"`
	AdminAt   int64  `bson:"admin_time" json:"admin_time,omitempty"`
	CreatedAt int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt int64  `bson:"update_time" json:"update_time,omitempty"`
	Status    int8   `bson:"status" json:"status,omitempty"`
}

type MyGroupRes struct {
	ID        string `bson:"_id" json:"id"`
	UID       string `bson:"uid" json:"uid"`
	Avatar    string `bson:"avatar" json:"avatar,omitempty"`
	Name      string `bson:"name" json:"name,omitempty"`
	EncKey    string `bson:"enc_key" json:"enc_key"`
	Status    int8   `bson:"status" json:"status"`
	CreatedAt int64  `bson:"created_at" json:"create_time"`
}

type GroupIDsRes struct {
	ID string `bson:"gid" json:"id"`
}

type GroupUIDsRes struct {
	UID string `bson:"uid" json:"id"`
}

type GroupMembersInfoRes struct {
	ID       string `bson:"_id" json:"id"`
	UID      string `bson:"uid" json:"uid"`
	GID      string `bson:"gid" json:"gid"`
	Role     uint8  `bson:"role" json:"role"`
	Avatar   string `bson:"avatar" json:"avatar,omitempty"`
	MyAlias  string `bson:"my_alias" json:"my_alias,omitempty"`
	AdminAt  int64  `bson:"admin_time" json:"admin_time,omitempty"`
	CreateAt int64  `bson:"create_time" json:"create_time,omitempty"`
}

type BaseInfoResponse struct {
	UID      string `bson:"uid" json:"uid"`
	Role     uint8  `bson:"role" json:"role"`
	CreateAt int64  `bson:"create_time" json:"create_time"`
}

type EncInfoResponse struct {
	GID    string `bson:"gid" json:"gid"`
	EncPri string `bson:"enc_pri" json:"enc_pri,omitempty"`
	EncKey string `bson:"enc_key" json:"enc_key,omitempty"`
}

type ApplyRes struct {
	ID        string `bson:"_id" json:"id"`
	GID       string `bson:"gid" json:"gid"`
	UID       string `bson:"uid" json:"uid"`
	Avatar    string `bson:"avatar" json:"avatar,omitempty"`
	Name      string `bson:"name" json:"name,omitempty"`
	EncKey    string `bson:"enc_key" json:"enc_key"`
	Role      uint8  `bson:"role" json:"role"`
	Status    int8   `bson:"status" json:"status"`
	CreatedAt int64  `bson:"created_at" json:"create_time"`
}

const (
	RoleDef uint8 = iota
	RoleOwner
	RoleAdministrator
	RoleCommonMember
)
const (
	JoinTypeInvite = 2
	JoinTypeSelf   = 1
)
const (
	StatusIng, StatusYes, StatusRefuse = 0, 1, 2
)

func New() *Members {
	return new(Members)
}

func (m Members) GetId(uid, gid string) string {
	return funcs.Md516(uid + gid)
}

func (m Members) TableName() string {
	return "group_members"
}

func (m Members) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(m.TableName(), rpMode)
}

func (m Members) Init() error {
	ctx := context.Background()
	indexModel := mongoDriver.IndexModel{
		Keys: []bson.M{{"uid": 1}, {"gid": 1}},
	}
	_, err := m.collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (m Members) GetFriendIds(uid string) ([]string, error) {
	type Data struct {
		ObjUId string `json:"obj_uid"`
	}
	data := make([]Data, 1000)
	m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"uid": uid}).Fields(bson.M{"obj_uid": 1}).FindMany(&data)
	res := make([]string, 0)
	if len(res) > 0 {
		for _, v := range data {
			res = append(res, v.ObjUId)
		}
	}
	return res, nil
}

func (m Members) AddMany(data []Members) error {
	res, err := m.collection().InsertMany(&data)
	fmt.Println(res, err)
	return err
}

func (m Members) AddOne(data Members) bool {
	_, err := m.collection().InsertOne(&data)
	return err == nil
}

func (m Members) UpByID(id string, uData Members) bool {
	_, err := m.collection().UpByID(id, uData)
	return err == nil
}

func (m Members) UpByIDs(ids []string, uData Members) error {
	_, err := m.collection().Where(bson.M{"_id": bson.M{"$in": ids}}).UpdateMany(&uData)
	return err
}

func (m Members) UpAliasByID(id string, uData map[string]interface{}) bool {
	_, err := m.collection().UpByID(id, uData)
	return err == nil
}

func (m Members) GetByUidAndGid(uid, gid, fields string) (Members, error) {
	id := m.GetId(uid, gid)
	return m.GetByID(id, fields)
}

func (m Members) Delete(uid, gid string) error {
	id := m.GetId(uid, gid)
	_, err := m.collection().Where(bson.M{"_id": id}).Delete()
	return err
}

func (m Members) GetByID(id string, fields string) (Members, error) {
	var data Members
	err := m.collection(mongo.SecondaryPreferredMode).Fields(dao.GetMongoFieldsBsonByString(fields)).FindByID(id, &data)
	return data, err
}

func (m Members) UpsertOne(data Members) error {
	err := m.collection().Where(bson.M{"_id": data.ID}).Upsert(data)
	return err.Err()
}

func (m Members) GetByGidAndUids(gid string, uids []string, fields string) ([]Members, error) {
	ids := make([]string, 0)
	for _, uid := range uids {
		id := m.GetId(uid, gid)
		ids = append(ids, id)
	}

	var data []Members
	m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": ids}}).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, nil
}

func (m Members) GetMyGroupList(uid string, gids []string) ([]MyGroupRes, error) {
	ids := make([]string, 0)
	where := bson.M{"uid": uid}
	if len(gids) > 0 {
		for _, id := range gids {
			id1 := m.GetId(uid, id)
			ids = append(ids, id1)
		}
		where = bson.M{"_id": bson.M{"$in": ids}}
	}
	var items []Members
	res := make([]MyGroupRes, 0)
	m.collection(mongo.SecondaryPreferredMode).Where(where).Fields(bson.M{"gid": 1}).FindMany(&items)
	if len(items) == 0 {
		return res, nil
	}
	gidSlice := make([]string, 0)
	for _, item := range items {
		gidSlice = append(gidSlice, item.GroupID)
	}
	groups, err := detail.New().GetInfoByIds(gidSlice)
	if groups == nil {
		return res, err
	}
	for _, v := range groups {
		item := MyGroupRes{
			ID:        v.ID,
			Avatar:    v.Avatar,
			Name:      v.Name,
			CreatedAt: v.CreatedAt,
		}
		res = append(res, item)
	}
	return res, err
}

func (m Members) GetMyGroupIdList(uid string) ([]GroupIDsRes, error) {
	res := make([]GroupIDsRes, 0)
	m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"uid": uid}).Fields(bson.M{"gid": 1}).FindMany(&res)
	return res, nil
}

func (m Members) GetApplyByGIds(gIds []string, status []int8) ([]ApplyRes, error) {
	res := make([]ApplyRes, 0)
	where := bson.M{"gid": bson.M{"$in": gIds}}
	if len(status) > 0 {
		where["status"] = bson.M{"$in": status}
	}
	err := m.collection().Where(where).FindMany(&res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (m Members) GetMyGroupByIds(uid string, ids []string) ([]ApplyRes, error) {
	res := make([]ApplyRes, 0)
	where := bson.M{"uid": uid}
	if len(ids) > 0 {
		where["gid"] = bson.M{"$in": ids}
	}
	m.collection(mongo.SecondaryPreferredMode).Where(where).Fields(bson.M{"gid": 1}).FindMany(&res)
	return res, nil
}

func (m Members) GetListsByStatusAndRole(uid string, ids []string, status []int8, role []int8) ([]ApplyRes, error) {
	res := make([]ApplyRes, 0)
	where := bson.M{"uid": uid}
	if len(ids) > 0 {
		where["gid"] = bson.M{"$in": ids}
	}

	if len(status) > 0 {
		where["status"] = bson.M{"$in": status}
	}
	if len(role) > 0 {
		where["role"] = bson.M{"$in": role}
	}
	m.collection().Where(where).FindMany(&res)
	return res, nil
}

func (m Members) GetGroupMemberIds(gid string) ([]GroupUIDsRes, error) {
	res := make([]GroupUIDsRes, 0)
	m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"gid": gid}).Fields(bson.M{"uid": 1}).FindMany(&res)
	return res, nil
}

func (m Members) GetGroupsMemberInfo(gid []string, fields string) ([]Members, error) {
	var data []Members
	m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"gid": bson.M{"$in": gid}}).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, nil
}

func (m Members) GetMembersInfo(gid string, uids []string) ([]GroupMembersInfoRes, error) {
	data := make([]GroupMembersInfoRes, 0)
	fields := dao.GetMongoFieldsBsonByString("_id,uid,gid,enc_key,enc_pri,pub,role,my_alias,admin_time,create_time")
	where := bson.M{"gid": gid}
	if len(uids) != 0 {
		where["uid"] = bson.M{"$in": uids}
	}
	err := m.collection(mongo.SecondaryPreferredMode).Where(where).Fields(fields).FindMany(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m Members) GetMembersByIds(gids, uids []string) ([]GroupMembersInfoRes, error) {
	data := make([]GroupMembersInfoRes, 0)
	fields := dao.GetMongoFieldsBsonByString("_id,uid,gid,role,my_alias,admin_time,create_time")
	where := bson.M{"gid": bson.M{"$in": gids}}
	if len(uids) != 0 {
		where["uid"] = bson.M{"$in": uids}
	}
	err := m.collection(mongo.SecondaryPreferredMode).Where(where).Fields(fields).FindMany(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m Members) GetEncInfoByIds(gids, uids []string) ([]EncInfoResponse, error) {
	data := make([]EncInfoResponse, 0)
	where := bson.M{"gid": bson.M{"$in": gids}}
	if len(uids) != 0 {
		where["uid"] = bson.M{"$in": uids}
	}
	err := m.collection(mongo.SecondaryPreferredMode).Where(where).FindMany(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m Members) IsExist(uid, gid string) bool {
	where := bson.M{"_id": m.GetId(uid, gid)}
	count := m.collection(mongo.SecondaryPreferredMode).Where(where).Count()
	return count > 0
}

func (m Members) UpdateRoleByUIds(uIds []string, gid string, role uint8) error {
	t := funcs.GetMillis()
	adminTime := 0
	if role == RoleAdministrator {
		adminTime = int(t)
	}
	uData := map[string]interface{}{
		"role":        role,
		"admin_time":  adminTime,
		"update_time": t,
	}
	where := bson.M{"gid": gid, "uid": bson.M{"$in": uIds}}
	_, err := m.collection().Where(where).UpdateMany(uData)
	return err
}

func (m Members) UpdateMemberRole(uid, gid string, role uint8) error {
	t := funcs.GetMillis()
	adminTime := 0
	if role == RoleAdministrator {
		adminTime = int(t)
	}
	uData := map[string]interface{}{
		"role":        role,
		"admin_time":  adminTime,
		"update_time": t,
	}
	id := m.GetId(uid, gid)
	_, err := m.collection().UpByID(id, uData)
	return err
}

func (m Members) GetAdministrator(gid string) []Members {
	data := make([]Members, 0)
	where := bson.M{"gid": gid, "role": RoleAdministrator}
	m.collection(mongo.SecondaryPreferredMode).Where(where).Fields(bson.M{"uid": 1}).Sort(bson.M{"create_time": 1}).FindMany(&data)
	return data
}

func (m Members) GetOwnerInfo(gid string) (Members, error) {
	var data Members
	where := bson.M{"gid": gid, "role": RoleOwner}
	err := m.collection().Where(where).Fields(bson.M{"uid": 1}).Sort(bson.M{"create_time": 1}).FindOne(&data)
	return data, err
}

func (m Members) GetRandomUidNotOwner(gid string) Members {
	var data Members
	where := bson.M{"gid": gid, "role": bson.M{"$ne": RoleOwner}}
	m.collection(mongo.SecondaryPreferredMode).Where(where).Sort(bson.M{"create_time": 1}).FindOne(&data)
	return data
}

func (m Members) RemoveMembers(ids []string) error {
	where := bson.M{"_id": bson.M{"$in": ids}}
	_, err := m.collection().Where(where).Delete()
	return err
}

func (m Members) GetGroupMembers(gid string) ([]string, error) {
	res := make([]string, 0)
	data, err := m.GetGroupMemberIds(gid)
	if data == nil {
		return []string{}, err
	}
	for _, v := range data {
		res = append(res, v.UID)
	}
	return res, err
}

func (m Members) GetMemberIds(gid string) ([]string, error) {
	cacheTag := funcs.Md516(gid)
	res := redis.Client.SMembers(cacheTag).Val()
	if len(res) == 0 {
		data, err := m.GetGroupMemberIds(gid)
		if data == nil {
			return []string{}, err
		}
		for _, v := range data {
			res = append(res, v.UID)
		}
		redis.Client.SAdd(cacheTag, res)
		redis.Client.Expire(cacheTag, time.Second*2)
	}
	return res, nil
}

func (m Members) GetGroupMember(gid string) ([]BaseInfoResponse, error) {
	data := make([]BaseInfoResponse, 0)
	fields := dao.GetMongoFieldsBsonByString("uid,role,create_time")
	err := m.collection().Where(bson.M{"gid": gid}).Fields(fields).FindMany(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m Members) Save(data Members) error {
	_, err := m.collection().Where(bson.M{"_id": data.ID}).InsertOne(&data)
	return err
}

func (m Members) GetGroupMyOwner(uid string) ([]GroupMembersInfoRes, error) {
	data := make([]GroupMembersInfoRes, 0)
	where := bson.M{"uid": uid, "role": RoleOwner}
	err := m.collection().Where(where).FindMany(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (m Members) GetMyGroupByGid(uid string, gids []string, fields string) ([]Members, error) {
	var data []Members
	mIds := make([]string, 0)
	for _, id := range gids {
		mIds = append(mIds, m.GetId(uid, id))
	}
	err := m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": mIds}}).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, err
}
