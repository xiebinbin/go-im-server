package message

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"imsdk/internal/common/dao/message/button"
	"imsdk/internal/common/dao/message/detail"
	"imsdk/internal/common/dao/message/usermessage"
	"imsdk/internal/common/model/errors"
	"imsdk/internal/common/model/forward"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"math"
)

type SetDisableParams struct {
	MId       string   `json:"mid" binding:"required"`
	UId       string   `json:"uid"`
	ButtonIds []string `json:"button_ids"`
}

type SetMessageDisableMultiUIdParams struct {
	MId       string   `json:"mid" binding:"required"`
	UIds      []string `json:"uids"`
	ButtonIds []string `json:"button_ids"`
}

func DealMessageButton(params SendMessageParams) error {
	if params.Type != base.MsgTypeCard && params.Type != base.MsgTypeRedEnvelope {
		return nil
	}
	buttonRelation := make(map[string]int, 0)
	if params.Type == base.MsgTypeCard {
		parseData := detail.ParseContentType19(params.Content)
		if len(parseData.Buttons) > 0 {
			for k, buttons := range parseData.Buttons {
				buttonRelation[buttons.ButtonId] = int(math.Pow(2, float64(k+1)))
			}
		}
	}

	if params.Type == base.MsgTypeRedEnvelope {
		buttonRelation = map[string]int{
			"RedEnvelope": int(math.Pow(2, 1)),
		}
	}

	t := funcs.GetMillis()
	err := button.New().InsertOne(button.Button{
		ID:        params.Mid,
		Button:    buttonRelation,
		CreatedAt: t,
	})
	if err != nil && mongo.IsNetworkError(err) {
		return errors.ErrSdkDefErr
	}
	return nil
}

func DealMessageButtonTemp(params SendMessageParamsTemp) error {
	if params.Type != base.MsgTypeCard && params.Type != base.MsgTypeRedEnvelope {
		return nil
	}
	buttonRelation := make(map[string]int, 0)
	if params.Type == base.MsgTypeCard {
		parseData := detail.ParseContentType19(params.Content)
		if len(parseData.Buttons) > 0 {
			for k, buttons := range parseData.Buttons {
				buttonRelation[buttons.ButtonId] = int(math.Pow(2, float64(k+1)))
			}
		}
	}

	if params.Type == base.MsgTypeRedEnvelope {
		buttonRelation = map[string]int{
			"RedEnvelope": int(math.Pow(2, 1)),
		}
	}

	err := button.New().InsertOne(button.Button{
		ID:        params.Mid,
		Button:    buttonRelation,
		CreatedAt: params.CreateTime,
	})
	if err != nil && mongo.IsNetworkError(err) {
		return errors.ErrSdkDefErr
	}
	return nil
}

func GetButtonInfo(mid string) map[string]int {
	data, err := button.New().GetById(mid)
	if err != nil && err == mongo.ErrNoDocuments {
		return map[string]int{}
	}
	return data.Button
}

func SetDisable(ctx *gin.Context, params SetDisableParams) error {
	msgDetail, _ := detail.New().GetDetailById(params.MId, "type")
	fmt.Println("msgDetail:", msgDetail)
	btNums := make([]int, 0)
	btInfos := make(map[string]int)
	btInfos = GetButtonInfo(params.MId)
	for _, id := range params.ButtonIds {
		if v, ok := btInfos[id]; ok {
			btNums = append(btNums, v)
		}
	}
	if len(btNums) == 0 {
		return nil
	}
	status := btNums[0]
	for _, num := range btNums {
		status = status | num
	}
	dao := usermessage.New()
	msgInfo := dao.GetById(params.UId, params.MId)
	_, err := dao.SetMsgDisable(params.UId, params.MId, status|int(msgInfo.Status))
	if err != nil {
		return err
	}
	pushParams := forward.PushMessageParams{
		Cmd: forward.CmdCardMessageStatus,
		Uid: params.UId,
		Data: MsgItem{
			ID: params.MId,
		},
		NoPushOffline: true,
	}
	err = forward.PushMessageToUserSocketDirectly(ctx, pushParams, false)
	if err != nil {
		return err
	}
	return nil
}

func SetMessageDisableMultiUId(ctx context.Context, params SetMessageDisableMultiUIdParams) error {
	//params.UId = "66594b182bdb8230fcbd6da5850a8487"
	msgDetail, _ := detail.New().GetDetailById(params.MId, "chat_id,type")
	btNums := make([]int, 0)
	btInfos := make(map[string]int, 0)
	logCtx := log.WithFields(ctx, map[string]string{"action": "SetMessageDisableMultiUId"})
	log.Logger().Info(logCtx, "msgDetail:", msgDetail, params.UIds)
	btInfos = GetButtonInfo(params.MId)
	for _, id := range params.ButtonIds {
		if v, ok := btInfos[id]; ok {
			btNums = append(btNums, v)
		}
	}
	if len(btNums) == 0 {
		return nil
	}
	status := btNums[0]
	for _, num := range btNums {
		status = status | num
	}
	dao := usermessage.New()
	for _, id := range params.UIds {
		msgInfo := dao.GetById(id, params.MId)
		_, err := dao.SetMsgDisable(id, params.MId, status|int(msgInfo.Status))
		if err != nil {
			log.Logger().Error(logCtx, "failed to set msg disable, err : ", err)
			//return err
		}
		pushParams := forward.PushMessageParams{
			Cmd: forward.CmdCardMessageStatus,
			Uid: id,
			Data: MsgItem{
				ID: params.MId,
			},
			NoPushOffline: true,
		}

		if err = forward.PushMessageToUserSocketDirectly(ctx, pushParams, false); err != nil {
			log.Logger().Error(logCtx, "failed to forward msg, err : ", err)
		}
	}
	return nil
}
