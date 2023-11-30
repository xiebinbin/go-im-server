package button

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"imsdk/pkg/database/mongo"
	"imsdk/pkg/funcs"
)

type Button struct {
	ID        string         `bson:"_id" json:"id"` // message id
	AK        string         `bson:"ak" json:"ak,omitempty"`
	Button    map[string]int `bson:"button" json:"button,omitempty"`
	CreatedAt int64          `bson:"create_time" json:"create_time"`
}

func New() *Button {
	return new(Button)
}

func (b *Button) TableName() string {
	return "message_button"
}

func (b *Button) WithAKey(ak string) *Button {
	b.AK = ak
	return b
}

func (b *Button) getCollection(mode ...mongo.Mode) *mongo.CollectionInfo {
	rpMode := mongo.PrimaryMode
	if len(mode) > 0 && mode[0] > 0 {
		rpMode = mode[0]
	}
	return mongo.Database().SetTable(b.TableName(), rpMode)
}

// Init : init user message table (create index: sequence)
func (b *Button) Init(uid string) error {
	ctx := context.Background()
	indexModel := mongoDriver.IndexModel{
		Keys: bson.D{
			{"button", -1},
		},
	}
	_, err := b.getCollection().Collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (b *Button) Add(button Button) (interface{}, error) {
	t := funcs.GetMillis()
	button.CreatedAt = t
	data, er := b.getCollection().InsertOne(&button)
	if er != nil {
		return "", er
	}
	return data.InsertedID, er
}

func (b *Button) InsertOne(button Button) error {
	return b.getCollection().Where(bson.M{"_id": button.ID}).Upsert(&button).Err()
}

func (b *Button) GetById(id string) (Button, error) {
	var data Button
	err := b.getCollection(mongo.SecondaryPreferredMode).FindByID(id, &data)
	return data, err
}

func (b *Button) Count() int64 {
	return b.getCollection().Count()
}

func (b *Button) DropTable() error {
	return mongo.Database().Database.Collection(b.TableName()).Drop(context.Background())
}
