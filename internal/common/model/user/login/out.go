package login

import (
	"context"

	"imsdk/internal/common/model/user/token"
)

type OutRequest struct {
	UId        string `json:"uid"`
	Os         string `json:"os"`
	DeviceId   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}

func Out(ctx context.Context, request OutRequest) bool {
	token.RemoveToken(ctx, request.UId, request.DeviceId)
	return true
}
