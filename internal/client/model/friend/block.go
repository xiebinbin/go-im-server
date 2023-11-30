package friend

import (
	"imsdk/internal/common/dao/block"
)

type BlockRequest struct {
	ObjUID string `json:"obj_uid" binding:"required"`
}

func GetMineBlockList(uid string) ([]string, error) {
	data, err := block.New().GetMineBlockIds(uid)
	return data, err
}

func GetBlockMineList(uid string) ([]string, error) {
	data, err := block.New().GetBlockMineIds(uid)
	return data, err
}

func IsBlockedMine(uid, objUId string) (bool, error) {
	return block.New().IsBlockedMine(uid, objUId)
}

func IsMineBlocked(uid, objUId string) (bool, error) {
	return block.New().IsMineBlocked(uid, objUId)
}
