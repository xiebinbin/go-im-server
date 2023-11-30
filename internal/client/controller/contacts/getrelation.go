package contacts

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imsdk/internal/client/model/friend"
	"imsdk/pkg/response"
)

type RelationRequest struct {
	UIDs []string `json:"uids" binding:"required"`
}

//type ApplyRequest = contacts.ApplyRequest

func GetRelationList(ctx *gin.Context) {
	userId, _ := ctx.Get("uid")
	uid := userId.(string)
	var params RelationRequest
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	res := friend.GetRelationInfo(uid, params.UIDs)
	response.RespListData(ctx, res)
	return
}
