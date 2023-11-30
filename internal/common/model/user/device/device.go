package device

import (
	"context"
	"imsdk/internal/common/dao/user/device"
	"imsdk/pkg/funcs"
)

type AddRequest struct {
	UID      string `json:"uid"`
	DeviceId string `json:"device_id"`
	OS       string `json:"os"`
}

type GetUserDeviceRequest struct {
	UID string `json:"uid"`
}

func Add(ctx context.Context, request AddRequest) error {
	t := funcs.GetMillis()
	err := device.New().UpsertOne(device.Device{
		ID:        device.GetId(request.UID, request.DeviceId, request.OS),
		DeviceId:  request.DeviceId,
		UID:       request.UID,
		CreateAt:  t,
		UpdatedAt: t,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetUserDevice(ctx context.Context, request GetUserDeviceRequest) []string {
	var res []string
	devices, err := device.New().GetDeviceIdsByUID(request.UID)
	if err == nil {
		for _, d := range devices {
			res = append(res, d.DeviceId)
		}
	}
	return res
}
