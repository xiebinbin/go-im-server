package contacts

import (
	"encoding/json"
	"imsdk/internal/client/model/friend"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/response"

	"github.com/gin-gonic/gin"
)

type RelationRequest struct {
	UIDs []string `json:"uids" binding:"required"`
}

//type ApplyRequest = contacts.ApplyRequest

func GetRelationList(ctx *gin.Context) {
	userId, _ := ctx.Get(base.HeaderFieldUID)
	uid := userId.(string)
	var params RelationRequest
	data, _ := ctx.Get("data")
	json.Unmarshal([]byte(data.(string)), &params)
	res := friend.GetRelationInfo(uid, params.UIDs)
	response.RespListData(ctx, res)
	return
}
