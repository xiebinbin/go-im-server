package pkg

import (
	"context"
	"imsdk/internal/common/pkg/base"
)

func GetDeviceId(ctx context.Context) string {
	if deviceId := ctx.Value(base.HeaderFieldDeviceId); deviceId != nil {
		return deviceId.(string)
	}
	return ""
}
