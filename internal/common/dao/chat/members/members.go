package members

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao"
	"imsdk/internal/common/dao/chat"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
	"imsdk/pkg/redis"
	"imsdk/pkg/sdk"
	"time"
)

type JoinType = sdk.JoinChatType
type RoleType = uint8
type Members struct {
	ID         string   `bson:"_id" json:"id"`
	ChatId     string   `bson:"chat_id" json:"chat_id,omitempty"`
	UID        string   `bson:"uid" json:"uid,omitempty"`
	InviteUID  string   `bson:"invite_uid" json:"invite_uid,omitempty"`
	Role       RoleType `bson:"role" json:"role,omitempty"`
	JoinType   JoinType `bson:"join_type" json:"join_type,omitempty"`
	MyAlias    string   `bson:"my_alias" json:"my_alias,omitempty"`
	AdminAt    int64    `bson:"admin_time" json:"admin_time,omitempty"`
	OnlineTime int64    `bson:"online_time" json:"online_time,omitempty"`
	MaxReadSeq int64    `bson:"max_read_seq" json:"max_read_seq,omitempty"`
	Status     int8     `bson:"status" json:"status"` //1-normal 2-forbidden send msg
	CreatedAt  int64    `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt  int64    `bson:"update_time" json:"update_time,omitempty"`
}

const (
	StatusNormal = 1
	StatusBlock  = 2
	StatusDelete = 3
)

type MyChatRes struct {
	ID     string `json:"id"`
	UID    string `json:"uid"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}

type ChatIDsRes struct {
	ID string `bson:"chat_id" json:"id"`
}

type ChatUIDsRes struct {
	UID string `bson:"uid" json:"id"`
}

type ChatMembersInfoRes struct {
	ID       string   `bson:"_id" json:"id"`
	UID      string   `bson:"uid" json:"uid"`
	ChatId   string   `bson:"chat_id" json:"chat_id"`
	Role     RoleType `bson:"role" json:"role"`
	Avatar   string   `bson:"my_alias" json:"my_alias"`
	AdminAt  int64    `bson:"admin_time" json:"admin_time"`
	CreateAt int64    `bson:"create_time" json:"create_time"`
}

type BaseInfoResponse struct {
	UID      string   `bson:"uid" json:"uid"`
	Role     RoleType `bson:"role" json:"role"`
	CreateAt int64    `bson:"create_time" json:"create_time"`
}

const (
	RoleDef RoleType = iota
	RoleOwner
	RoleAdministrator
	RoleCommonMember
)

func New() *Members {
	return new(Members)
}

func (m *Members) GetId(uid, chatId string) string {
	return funcs.Md516(uid + chatId)
}

func (m *Members) TableName() string {
	return "chat_members"
}

func (m *Members) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(m.TableName(), rpMode)
}

func (m *Members) Init() error {
	ctx := context.Background()
	indexModel := mongoDriver.IndexModel{
		Keys: []bson.M{{"uid": 1}, {"gid": 1}},
	}
	_, err := m.collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (m *Members) GetFriendIds(uid string) ([]string, error) {
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

func (m *Members) AddMany(data []Members) error {
	_, err := m.collection().InsertMany(&data)
	return err
}

func (m *Members) AddOne(data Members) error {
	_, err := m.collection().InsertOne(&data)
	return err
}

func (m *Members) UpsertOne(data Members) error {
	err := m.collection().Where(bson.M{"_id": data.ID}).Upsert(&data).Err()
	return err
}

func (m *Members) UpByID(id string, uData Members) bool {
	_, err := m.collection().UpByID(id, uData)
	return err == nil
}

func (m *Members) UpAliasByID(id string, uData map[string]interface{}) bool {
	_, err := m.collection().UpByID(id, uData)
	return err == nil
}

func (m *Members) GetByUidAndGid(uid, gid, fields string) (Members, error) {
	id := m.GetId(uid, gid)
	return m.GetByID(id, fields)
}

func (m *Members) Delete(uid, chatId string) error {
	id := m.GetId(uid, chatId)
	_, err := m.collection().Where(bson.M{"_id": id}).Delete()
	return err
}

func (m *Members) DeleteMany(ids []string) (int64, error) {
	count, err := m.collection().Where(bson.M{"_id": bson.M{"$in": ids}}).Delete()
	return count, err
}

func (m *Members) DeleteByChatId(chatId string) error {
	_, err := m.collection().Where(bson.M{"chat_id": chatId}).Delete()
	return err
}

func (m *Members) GetByID(id string, fields string) (Members, error) {
	var data Members
	err := m.collection(mongo.SecondaryPreferredMode).Fields(dao.GetMongoFieldsBsonByString(fields)).FindByID(id, &data)
	return data, err
}

func (m *Members) GetByGidAndUids(gid string, uids []string, fields string) ([]Members, error) {
	ids := make([]string, 0)
	for _, uid := range uids {
		id := m.GetId(uid, gid)
		ids = append(ids, id)
	}

	var data []Members
	m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": ids}}).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, nil
}

func (m *Members) GetMyChatList(uid string, gids []string) ([]MyChatRes, error) {
	ids := make([]string, 0)
	for _, id := range gids {
		id = m.GetId(uid, id)
		ids = append(ids, id)
	}
	var items []Members
	res := make([]MyChatRes, 0)
	m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": ids}}).Fields(bson.M{"gid": 1}).FindMany(&items)
	if len(items) == 0 {
		return res, nil
	}
	gidSlice := make([]string, 0)
	for _, item := range items {
		gidSlice = append(gidSlice, item.ChatId)
	}
	chats, err := chat.New().GetInfoByIds(gidSlice, "_id,avatar,name")
	if chats == nil {
		return res, err
	}
	for _, v := range chats {
		item := MyChatRes{
			ID:     v.ID,
			UID:    uid,
			Avatar: v.Avatar,
			Name:   v.Name,
		}
		res = append(res, item)
	}
	return res, err
}

func (m *Members) GetMyChatIdList(uid string) ([]Members, error) {
	res := make([]Members, 0)
	err := m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"uid": uid}).Fields(bson.M{"chat_id": 1}).FindMany(&res)
	return res, err
}

func (m *Members) GetChatMemberIds(gid string) ([]ChatUIDsRes, error) {
	res := make([]ChatUIDsRes, 0)
	//where := bson.M{"chat_id": gid, "status": StatusNormal}
	where := bson.M{"chat_id": gid}
	m.collection(mongo.SecondaryPreferredMode).Where(where).Fields(bson.M{"uid": 1}).FindMany(&res)
	return res, nil
}

func (m *Members) GetChatMemberIdsSortByOnline(gid string) ([]ChatUIDsRes, error) {
	res := make([]ChatUIDsRes, 0)
	m.collection(mongo.SecondaryPreferredMode).Sort(bson.M{"online_time": -1}).Where(bson.M{"chat_id": gid}).Fields(bson.M{"uid": 1}).FindMany(&res)
	return res, nil
}

// GetChatOtherMemberId GetChatMemberIds --temporary use
func (m *Members) GetChatOtherMemberId(gid, senderId string) (ChatUIDsRes, error) {
	var res ChatUIDsRes
	m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"chat_id": gid, "uid": bson.M{"$ne": senderId}}).Fields(bson.M{"uid": 1}).FindOne(&res)
	return res, nil
}

func (m *Members) GetChatsMemberInfo(chatId []string, fields string) ([]Members, error) {
	var data []Members
	m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"chat_id": bson.M{"$in": chatId}}).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, nil
}

func (m *Members) GetMembersInfo(gid string, uids []string) ([]ChatMembersInfoRes, error) {
	data := make([]ChatMembersInfoRes, 0)
	fields := dao.GetMongoFieldsBsonByString("_id,uid,gid,role,my_alias,admin_time,create_time")
	where := bson.M{"gid": gid, "uid": bson.M{"$in": uids}}
	m.collection(mongo.SecondaryPreferredMode).Where(where).Fields(fields).FindMany(&data)
	return data, nil
}

func (m *Members) IsExist(uid, chatId string) bool {
	where := bson.M{"_id": m.GetId(uid, chatId)}
	count := m.collection(mongo.SecondaryPreferredMode).Where(where).Count()
	return count > 0
}

func (m *Members) UpdateRoleByUIds(uIds []string, chatId string, role uint8) error {
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
	where := bson.M{"chat_id": chatId, "uid": bson.M{"$in": uIds}}
	_, err := m.collection().Where(where).UpdateMany(uData)
	return err
}

func (m *Members) UpdateMemberRole(uid, gid string, role uint8) error {
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

func (m *Members) GetAdministrator(gid string) []Members {
	data := make([]Members, 0)
	where := bson.M{"gid": gid, "role": RoleAdministrator}
	m.collection(mongo.SecondaryPreferredMode).Where(where).Fields(bson.M{"uid": 1}).Sort(bson.M{"create_time": 1}).FindMany(&data)
	return data
}

func (m *Members) GetOwnerInfo(gid string) (Members, error) {
	var data Members
	where := bson.M{"gid": gid, "role": RoleOwner}
	err := m.collection().Where(where).Fields(bson.M{"uid": 1}).Sort(bson.M{"create_time": 1}).FindOne(&data)
	return data, err
}

func (m *Members) GetRandomUidNotOwner(gid string) Members {
	var data Members
	where := bson.M{"gid": gid, "role": bson.M{"$ne": RoleOwner}}
	m.collection(mongo.SecondaryPreferredMode).Where(where).Sort(bson.M{"create_time": 1}).FindOne(&data)
	return data
}

func (m *Members) RemoveMembers(ids []string) error {
	where := bson.M{"_id": bson.M{"$in": ids}}
	_, err := m.collection().Where(where).Delete()
	return err
}

func (m *Members) GetChatMemberUIds(gid string) ([]string, error) {
	cacheTag := funcs.Md516(gid)
	res := redis.Client.SMembers(cacheTag).Val()
	if len(res) == 0 {
		//data, err := m.GetChatMemberIds(gid)
		data, err := m.GetChatMemberIdsSortByOnline(gid)
		if data == nil {
			return []string{}, err
		}

		for i := 0; i < len(data); i++ {
			res = append(res, data[i].UID)
		}
		redis.Client.SAdd(cacheTag, res)
		redis.Client.Expire(cacheTag, time.Second*2)
	}
	return res, nil
}

func (m *Members) GetChatMember(gid, field string) ([]Members, error) {
	data := make([]Members, 0)
	fields := dao.GetMongoFieldsBsonByString(field)
	m.collection().Where(bson.M{"gid": gid}).Fields(fields).FindMany(&data)
	return data, nil
}

func (m *Members) Save(data Members) error {
	_, err := m.collection().Where(bson.M{"_id": data.ID}).InsertOne(&data)
	return err
}

func (m *Members) GetChatMyOwner(uid string) ([]ChatMembersInfoRes, error) {
	data := make([]ChatMembersInfoRes, 0)
	where := bson.M{"uid": uid, "role": RoleOwner}
	err := m.collection().Where(where).FindMany(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (m *Members) GetMyChatByGid(uid string, gids []string, fields string) ([]Members, error) {
	var data []Members
	mIds := make([]string, 0)
	for _, id := range gids {
		mIds = append(mIds, m.GetId(uid, id))
	}
	err := m.collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": mIds}}).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, err
}

func (m *Members) GetRoleCount(gid string, roleType RoleType) int64 {
	where := bson.M{"chat_id": gid, "role": roleType}
	return m.collection().Where(where).Fields(bson.M{"uid": 1}).Sort(bson.M{"create_time": 1}).Count()
}

func (m *Members) GetChatMembersCount(gid string) int64 {
	where := bson.M{"chat_id": gid}
	return m.collection().Where(where).Fields(bson.M{"uid": 1}).Count()
}

func (m *Members) GetChatMembers(gid string) ([]string, error) {
	res := make([]string, 0)
	data, err := m.GetChatMemberIds(gid)
	if data == nil {
		return []string{}, err
	}
	for _, v := range data {
		res = append(res, v.UID)
	}
	return res, err
}

func (m *Members) GetChatMemberUIdsWithoutOther(gid string, uIds []string) ([]string, error) {
	cacheTag := funcs.Md516(gid)
	res := redis.Client.SMembers(cacheTag).Val()
	if len(res) == 0 {
		data, err := m.GetChatMemberIds(gid)
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

func (m *Members) UpdateUserOnlineTime(uid string, time int64) error {
	_, err := m.collection().Where(bson.M{"uid": uid}).UpdateMany(map[string]interface{}{
		"online_time": time,
	})
	return err
}
