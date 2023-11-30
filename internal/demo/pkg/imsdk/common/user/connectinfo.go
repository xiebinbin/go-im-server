package user

import (
	"context"
	json "github.com/json-iterator/go"
)

type connectInfo struct {
	AUIds []string `json:"auids"`
}

func (c *Client) GetConnectInfo(ctx context.Context, aUIds []string) (interface{}, error) {
	dataByte, _ := json.Marshal(connectInfo{
		AUIds: aUIds,
	})
	data := NewRequest(ActionGetConnectInfo, string(dataByte))
	res, err := RequestIMServer(ctx, c.options, data)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	return res.Data.Items, nil
}
