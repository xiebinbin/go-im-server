package terminalinfo

import (
	"go.mongodb.org/mongo-driver/bson"
	"imsdk/pkg/database/mongo"
)

type TerminalInfo struct {
	ID        string `bson:"_id" json:"id"` // userId+source+date
	UID       string `bson:"uid" json:"uid"`
	Version   string `bson:"version" json:"version"`
	Os        string `bson:"os" json:"os"`
	OsVer     string `bson:"os_ver" json:"os_ver"`
	Ip        string `bson:"ip" json:"ip"`
	Lang      string `bson:"lang" json:"lang"`
	Country   string `bson:"country" json:"country"`
	Region    string `bson:"region" json:"region"`
	City      string `bson:"city" json:"city"`
	BandInfo  string `bson:"band_info" json:"band_info"`
	Band      string `bson:"band" json:"band"`
	BandModel string `bson:"band_model" json:"band_model"`
	Source    int8   `bson:"source" json:"source"`
	DateIdx   string `bson:"date_idx" json:"date_idx"`
	Date      string `bson:"date" json:"date"`
	CreateAt  int64  `bson:"create_time" json:"create_time"`
	UpdatedAt int64  `bson:"update_time" json:"update_time"`
}

type AggregateResp struct {
	ID    string `bson:"_id" json:"id"`
	Count uint64 `bson:"count" json:"count"`
	Type  []uint `bson:"type" json:"type"`
}
type SourceType = int8

const (
	SourceReg    SourceType = 1
	SourceLogin  SourceType = 2
	SourceActive SourceType = 3
)

func New() *TerminalInfo {
	return new(TerminalInfo)
}

func (t TerminalInfo) TableName() string {
	return "user_terminal_info"
}

func (t TerminalInfo) collection() *mongo.CollectionInfo {
	return mongo.Database().SetTable(t.TableName())
}

func (t TerminalInfo) Add(addData TerminalInfo) (string, error) {
	res, err := t.collection().InsertOne(addData)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(string), nil
}

func (t TerminalInfo) GetRegDate(uid string) (TerminalInfo, error) {
	var data TerminalInfo
	//, "source": SourceReg
	err := t.collection().Where(bson.M{"uid": uid}).FindOne(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (t TerminalInfo) InsertOrUpdate(data TerminalInfo) error {
	return t.collection().Where(bson.M{"_id": data.ID}).Upsert(&data).Err()
}

func (t TerminalInfo) GetCountByCond(where bson.M) int64 {
	return t.collection().Where(where).Count()
}

func (t TerminalInfo) GetCountByGroupField(where bson.M, field string) []AggregateResp {
	pipeline := []bson.M{
		bson.M{
			"$match": where,
		},
		bson.M{
			"$group": bson.M{
				"_id":   "$" + field,
				"count": bson.M{"$sum": 1},
			},
		},
	}
	var data []AggregateResp
	t.collection().Aggregate(pipeline, &data)
	return data
}
