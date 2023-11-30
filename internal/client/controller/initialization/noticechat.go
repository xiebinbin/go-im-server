package initialization

import (
	"context"
	"github.com/gin-gonic/gin"
	chat2 "imsdk/internal/common/dao/chat"
	"imsdk/internal/common/model/chat"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/config"
	"imsdk/pkg/app"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/response"
)

func InitNoticeChat(ctx *gin.Context) {
	var params chat.CreateParams
	params.Id = base.OfficialNoticeChatId
	noticeUId, _ := config.GetHollowUId()
	params.Name = "Chat"
	//ak := funcs.GetEnvAk()
	ak, _ := config.GetConfigAk()
	if ak == app.OfficialAK {
		params.Creator = noticeUId
		params.Avatar = ""
	}
	_, err := createNoticeChat(ctx, params)
	if err != nil {
		response.RespErr(ctx, err)
		return
	}
	response.RespSuc(ctx)
	return
}

func createNoticeChat(ctx context.Context, params chat.CreateParams) (chat.CreateResp, error) {
	// todo if group exist
	logCtx := log.WithFields(ctx, map[string]string{"action": "createChat"})
	var res chat.CreateResp
	// save group detail data
	t := funcs.GetMillis()
	chatDetail := chat2.Chat{
		ID:        params.Id,
		Name:      params.Name,
		Total:     0,
		Avatar:    params.Avatar,
		Status:    1,
		CreatedAt: t,
		UpdatedAt: t,
	}
	//log.Logger().Info(logCtx, "chatDetail:", chatDetail)
	if err := chat2.New().Upsert(chatDetail); err != nil {
		log.Logger().Info(logCtx, "save group detail unsuccessfully, err: ", err)
		return res, errno.Add("failed to save chat", errno.Exception)
	}
	return res, nil
}
