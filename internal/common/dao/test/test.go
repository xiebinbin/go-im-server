package test

import (
	"imsdk/pkg/database/mongo"
)

type Test struct {
	ID        string `bson:"_id" json:"id"`
	Seq       int64  `bson:"seq" json:"seq"`
	CreatedAt int64  `bson:"create_time" json:"create_time"`
	UpdatedAt int64  `bson:"update_time" json:"update_time"`
}

func (t Test) TableName() string {
	return "test"
}

func New() *Test {
	return new(Test)
}

func (t Test) GetCollection() *mongo.CollectionInfo {
	return mongo.Database().SetTable(t.TableName())
}

func (t Test) Add() {
	//data := Test{
	//	ID:  strconv.Itoa(int(unique.ID())),
	//	Seq: "getNextSequence(\"test1\")",
	//}
	//collection := t.GetCollection()
	//ctx:=context.Background()
	//collection.Table.Indexes().CreateOne(ctx)
	//_, err := collection.InsertOne(&data)
	//fmt.Println(err)
}

func (t Test) Get() ([]Test, error) {
	var data []Test
	err := t.GetCollection().FindMany(&data)
	return data, err
}

func (t Test) AddMany(addData []Test) ([]interface{}, error) {
	res, err := t.GetCollection().InsertMany(addData)
	if err != nil {
		return []interface{}{}, err
	}
	return res.InsertedIDs, nil
}
