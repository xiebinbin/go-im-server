package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"imsdk/internal/common/dao"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
)

type User struct {
	ID        string `bson:"_id" json:"address"`
	Avatar    string `bson:"avatar" json:"avatar"`
	Name      string `bson:"name" json:"name"`
	Gender    string `bson:"gender" json:"gender"` // 0-未知 1-男 2-女
	Sign      string `bson:"sign" json:"sign"`
	PubKey    string `bson:"pub_key" json:"pub_key"`
	Status    int8   `bson:"status" json:"status"`
	DBIdx     int64  `bson:"db_idx" json:"db_idx"`
	CreatedAt int64  `bson:"create_time" json:"create_time,omitempty"`
}

type ListResponse struct {
	ID     string `bson:"_id" json:"id"`
	Avatar string `bson:"avatar" json:"avatar"`
	Name   string `bson:"name" json:"name"`
	PubKey string `bson:"pub_key" json:"pub_key"`
	Gender string `bson:"gender" json:"gender"` // 0-未知 1-男 2-女
	Sign   string `bson:"sign" json:"sign"`
}

const (
	StatusNormal = 0
	StatusDelete = -1
	StatusForbid = -2
)

func New() *User {
	return new(User)
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(u.TableName(), rpMode)
}

// Init : init user table (create index: auid)
func (u *User) Init() error {
	isUnique := true
	opts := options.IndexOptions{Unique: &isUnique}
	indexModel := mongoDriver.IndexModel{
		Keys: bson.D{
			{"uid", -1},
		},
		Options: &opts,
	}

	ctx := context.Background()
	_, err := u.collection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (u *User) GetNameInfo(uIds []string) ([]User, error) {
	var data []User
	err := u.collection().Where(bson.M{"_id": bson.M{"$in": uIds}}).Fields(dao.GetMongoFieldsBsonByString("_id,name")).FindMany(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *User) GetById(id string) (string, error) {
	var data User
	err := u.collection().Where(bson.M{"_id": id}).FindOne(&data)
	if err != nil {
		return "", err
	}
	return data.ID, nil
}

func (u *User) GetByIds(Ids []string) (map[string]string, error) {
	var data []User
	result := make(map[string]string)
	err := u.collection().Where(bson.M{"_id": bson.M{"$in": Ids}}).FindMany(&data)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		result[v.ID] = v.ID
	}
	return result, nil
}

func (u *User) GetInfoByIds(Ids []string) (map[string]User, error) {
	var data []User
	result := make(map[string]User)
	err := u.collection().Where(bson.M{"_id": bson.M{"$in": Ids}}).FindMany(&data)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		result[v.ID] = v
	}
	return result, nil
}

func (u *User) GetSliceAUIds(uIds []string) ([]string, error) {
	var data []User
	result := make([]string, 0)
	err := u.collection(mongo.SecondaryPreferredMode).Where(bson.M{"_id": bson.M{"$in": uIds}}).FindMany(&data)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		result = append(result, v.ID)
	}
	return result, nil
}

func (u *User) Add(data User) error {
	data.DBIdx = int64(funcs.Str2Modulo(data.ID))
	_, err := u.collection().InsertOne(&data)
	if err != nil {
		if dao.DataIsSaveSuccessfully(err) {
			return nil
		}
		return err
	}
	return nil
}

func (u *User) Update(uid string, data User) error {
	_, err := u.collection().Where(bson.M{"_id": uid}).UpdateOne(&data)
	if err != nil {
		if dao.DataIsSaveSuccessfully(err) {
			return nil
		}
		return err
	}
	return nil
}

func (u *User) GetInfoById(addr string) (User, error) {
	var data User
	//cacheTag := "im:user:info:" + addr
	//res := redis.Client.Get(cacheTag).Val()
	//if res == "" {
	err := u.collection().Where(bson.M{"_id": addr}).FindOne(&data)
	if err != nil {
		return data, err
	}
	//	re, _ := json.Marshal(data)
	//	redis.Client.Set(cacheTag, string(re), time.Second*1800)
	//	return data, nil
	//}
	//err := json.Unmarshal([]byte(res), &data)
	//if err != nil {
	//	return data, err
	//}
	return data, nil
}

func (u *User) GetAddressByUID(uid string) (User, error) {
	var data User
	err := u.collection().Where(bson.M{"uid": uid}).FindOne(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (u *User) GetByIDs(uIds []string) ([]ListResponse, error) {
	var data []ListResponse
	err := u.collection().Where(bson.M{"_id": bson.M{"$in": uIds}}).FindMany(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (u *User) GetByID(id string) (User, error) {
	var data User
	err := u.collection().Where(bson.M{"_id": id}).FindOne(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (u *User) GetAll() ([]User, error) {
	var data []User
	err := u.collection().Where(bson.M{"status": StatusNormal}).FindMany(&data)
	return data, err
}

func (u *User) VerifyUserStatus(uid string) error {
	info, _ := u.GetInfoById(uid)
	if info.Status == StatusDelete {
		return errno.Add("user-status-delete", errno.UserDelete)
	}
	if info.Status == StatusForbid {
		return errno.Add("user-status-forbid", errno.UserUnavailable)
	}
	return nil
}

func (u *User) GetInfoByStatus(ids []string, status []int, fields string) ([]User, error) {
	var data []User
	where := bson.M{"_id": bson.M{"$in": ids}, "status": bson.M{"$in": status}}
	err := u.collection().Where(where).Fields(dao.GetMongoFieldsBsonByString(fields)).FindMany(&data)
	return data, err
}
