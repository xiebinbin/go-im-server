package message

import (
	"context"
	"imsdk/internal/demo/pkg/imsdk/message/msgtype"
)

func (c *Client) SetTemplateContent(ctx context.Context, temId, operator string, targets []string, number, duration int64) *Client {
	//c.Template = _type.Template{
	//	TemId:    temId,
	//	Operator: operator,
	//	Target:   targets,
	//	Duration: duration,
	//	Number:   number,
	//}
	return c
}

func (c *Client) SendTemplateMsg(ctx context.Context, msgType int8, aMId, aChatId string, action map[string]interface{}, receiveUid ...string) (uint64, error) {
	//content := c.GetTemplateCustomizeContent()
	//contentStr, err := json.Marshal(content)
	//if err != nil {
	//	return 0, err
	//}
	if msgType == 0 {
		msgType = msgtype.TypeCmd
	}
	//t := utils.GetMillis()
	msg := SendParams{
		//aMId:       aMId,
		Action: action,
		//AChatId:    aChatId,
		//SendUid:    MiddleMsgSendUID,
		//SendTime:   t,
		//ReceiveIds: receiveUid,
		//Content:    string(contentStr),
		//Type:       msgType,
	}
	_, err := c.sendMessage(ctx, msg)
	if err != nil {
		return 0, err
	}
	return 0, nil
}
