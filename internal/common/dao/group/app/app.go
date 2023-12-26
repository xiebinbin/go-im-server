package app

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
)

type App struct {
	ID        string `bson:"_id" json:"id"`
	GroupID   string `bson:"gid" json:"gid"`
	Name      string `bson:"name" json:"name"`
	Icon      string `bson:"icon" json:"icon"`
	Desc      string `bson:"desc" json:"desc"`
	Url       string `bson:"url" json:"url"`
	Sort      int    `bson:"sort" json:"sort"`
	Status    int    `bson:"status" json:"status"`
	CreatedAt int64  `bson:"created_at" json:"created_at"`
	UpdatedAt int64  `bson:"updated_at" json:"updated_at"`
}

const (
	StatusYes       = 1
	StatusForbidden = 2
)

func (a *App) TableName() string {
	return "group_app"
}

func New() *App {
	return new(App)
}

func (a *App) collection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(a.TableName(), rpMode)
}

func (a *App) Add(data App) error {
	_, err := a.collection().InsertOne(&data)
	return err
}

func (a *App) GetByID(id string, fields string) (App, error) {
	var data App
	err := a.collection().Fields(fields).Where(map[string]interface{}{
		"_id": id,
	}).FindOne(&data)
	return data, err
}

func (a *App) GetByIds(ids []string, status []int) ([]App, error) {
	var data []App
	where := bson.M{"_id": bson.M{"$in": ids}}
	if len(status) > 0 {
		where["status"] = bson.M{"$in": status}
	}
	err := New().collection().Where(where).FindMany(&data)
	return data, err
}

func (a *App) GetByGroupIds(gIds []string, status []int) ([]App, error) {
	var data []App
	where := bson.M{
		"gid": bson.M{"$in": gIds},
	}
	if len(status) > 0 {
		where["status"] = bson.M{"$in": status}
	}
	err := a.collection().Where(where).FindMany(&data)
	return data, err
}

func (a *App) UpdateById(id string, data App) error {
	err := a.collection().Where(bson.M{"_id": id}).Upsert(&data)
	return err.Err()
}

func (a *App) DeleteById(id string) error {
	_, err := a.collection().Where(bson.M{"_id": id}).Delete()
	return err
}

func (a *App) DeleteByIds(ids []string) error {
	_, err := a.collection().Where(bson.M{"_id": bson.M{"$in": ids}}).Delete()
	return err
}

func (a *App) DeleteByGroupIds(gIds []string) error {
	_, err := a.collection().Where(bson.M{"gid": bson.M{"$in": gIds}}).Delete()
	return err
}
