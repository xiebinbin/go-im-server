package message

import (
	"context"
	"fmt"
)

func (c *Client) SetMessageDisable(ctx context.Context, aMId string, buttonIds []string, aUIds ...string) (uint64, error) {
	disableData := DisableParams{
		Amid:      aMId,
		AUIds:     aUIds,
		ButtonIds: buttonIds,
	}
	fmt.Println("disableData:", disableData)
	_, err := c.SetDisable(ctx, disableData)
	if err != nil {
		return 0, err
	}
	return 0, nil
}
