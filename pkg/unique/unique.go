package unique

import (
	"github.com/globalsign/mgo/bson"
	"github.com/sony/sonyflake"
	"time"
)

var flake = sonyflake.NewSonyflake(sonyflake.Settings{
	StartTime: time.Date(2019, 8, 7, 0, 0, 0, 0, time.Local),
})

// ID 获取一个新的唯一id
func ID() uint64 {
	id, err := flake.NextID()
	if err != nil {
		return 0
	}
	return id
}

func Id12() string {
	return bson.NewObjectId().Hex()
}
