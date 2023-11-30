package message

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/model/message"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/response"
)

type UserMessageResponse struct {
	ID        string `bson:"_id" json:"id"`
	Type      uint16 `bson:"type" json:"type"`
	ChatId    string `bson:"chat_id" json:"chat_id"`
	AMID      string `bson:"amid" json:"amid"`
	SenderId  string `bson:"sender_id" json:"sender_id"`
	Content   string `bson:"content" json:"content"`
	SendTime  int64  `bson:"send_time" json:"send_time"`
	CreatedAt int64  `bson:"create_time" json:"create_time"`
}

func GetMessageByChatId(ctx *gin.Context) {
	var params message.GetMessageListRequest
	data, _ := ctx.Get("data")
	err := json.Unmarshal([]byte(data.(string)), &params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	data = message.GetMessageList(ctx, params)
	response.RespListData(ctx, data)
	return
}

func DeleteAll(ctx *gin.Context) {
	uid := ctx.Value(base.HeaderFieldUID)
	data := message.DeleteByUID(ctx, uid.(string))
	response.RespListData(ctx, data)
	return
}

func BatchDelete(ctx *gin.Context) {
	var params message.BatchDeleteRequest
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	err := message.BatchDelete(ctx, params.Ids)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}
