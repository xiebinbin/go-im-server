package message

import (
	"context"
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
	"imsdk/internal/common/dao/message/detail"
	"imsdk/internal/common/dao/message/revokebackup"
	"imsdk/internal/common/dao/message/usermessage"
	"imsdk/internal/common/pkg/base"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
)

type RevokeParams struct {
	UID    string   `json:"uid"`
	ChatId string   `json:"chat_id" binding:"required"`
	MIds   []string `json:"ids"`
	Os     string   `json:"os"`
}

type RevokeByChatIdsParams struct {
	ChatIds    []string `json:"chat_ids" binding:"required"`
	ExceptMIds []string `json:"except_mids"`
	UID        string   `json:"uid"`
	Os         string   `json:"os"`
}

func Revoke(ctx *gin.Context, params RevokeParams) ([]string, bool) {
	msgInfo := GetMsgInfo(ctx, params.MIds)
	logCtx := log.WithFields(ctx, map[string]string{"action": "Revoke"})
	log.Logger().Info(logCtx, "revoke msg ", params)
	undoIds := make([]string, 0)
	if len(msgInfo) > 0 {
		for _, v := range msgInfo {
			undoIds = append(undoIds, funcs.Md516(v.ID+HandleUndo))
		}
	}
	log.Logger().Info(logCtx, "revoke msg-1 ", undoIds)
	err := revokeBackup(ctx, params.UID, []string{}, params.MIds)
	if err != nil {
		return []string{}, false
	}
	_, err = detail.New().Delete(params.UID, params.MIds)
	if err == nil {
		sysMsgData := map[string]interface{}{
			"operator": params.UID,
			"target":   []string{},
			"temId":    "revoke-msg",
			"mids":     params.MIds,
		}
		contentByte, _ := json.Marshal(sysMsgData)
		data := SendMessageParams{
			Mid:     undoIds[0],
			Type:    base.MsgContentTypeHollowDel,
			ChatId:  params.ChatId,
			Content: string(contentByte),
			Extra:   "",
		}

		log.Logger().Info(logCtx, "revoke msg-2 ", data)
		if err = HollowManSendMsg(ctx, data, params.Os); err != nil {
			log.Logger().Error(logCtx, "failed to send ,err :  ", err)
			return undoIds, false
		}
		return undoIds, true
	}
	return undoIds, false
}

func RevokeByChatIds(ctx *gin.Context, params RevokeByChatIdsParams) ([]string, error) {
	msgInfo := detail.New().GetExceptByChatIdAndSenderId(params.UID, params.ChatIds, params.ExceptMIds, "_id")
	undoIds := make([]string, 0)
	if len(msgInfo) == 0 {
		return undoIds, nil
	}
	for _, v := range msgInfo {
		undoIds = append(undoIds, v.ID)
	}
	_, err := detail.New().Delete(params.UID, undoIds)
	logCtx := log.WithFields(ctx, map[string]string{"action": "RevokeByChatIds"})
	log.Logger().Error(logCtx, "RevokeByChatIds err:", err, params)
	if err == nil {
		sysMsgData := map[string]interface{}{
			"operator":    params.UID,
			"target":      []string{},
			"temId":       "revoke-msg-all",
			"chat_id":     params.ChatIds,
			"except_mids": params.ExceptMIds,
		}
		contentByte, _ := json.Marshal(sysMsgData)
		data := SendMessageParams{
			Mid:     funcs.CreateMsgId(params.UID),
			Type:    base.MsgContentTypeHollowDel,
			ChatId:  params.ChatIds[0],
			Content: string(contentByte),
			Extra:   "",
		}
		if err = HollowManSendMsg(ctx, data, params.Os); err != nil {
			return undoIds, err
		}
		return undoIds, nil
	}
	return undoIds, err
}

func ClearByChatIds(ctx *gin.Context, params RevokeByChatIdsParams) ([]string, error) {
	//logCtx := log.WithFields(ctx, map[string]string{"action": "ClearByChatIds"})
	_, err := usermessage.New().DeleteSelfMsgInChatIds(params.UID, params.ChatIds, params.ExceptMIds)
	if err != nil {
		return []string{}, err
	}
	err = revokeBackup(ctx, params.UID, params.ChatIds, []string{})
	if err != nil {
		return []string{}, err
	}
	detailDao := detail.New()
	recentIds, _ := detailDao.GetMyRecentByChatIds(params.UID, params.ChatIds, params.ExceptMIds, 20)
	_, err = detailDao.DeleteSelfByChatIds(params.UID, params.ChatIds, params.ExceptMIds)
	if len(params.ChatIds) > 0 && len(params.ExceptMIds) > 0 {
		chatId := params.ChatIds[0]
		_ = AddPrepareMsgByStatus(context.Background(), AddPrepareMsgByStatusReq{
			ChatId:   chatId,
			SenderId: params.UID,
			MIds:     params.ExceptMIds,
			Tag:      TagRevoke,
		})
	}
	sysMsgData := map[string]interface{}{
		"operator":    params.UID,
		"target":      []string{},
		"temId":       "revoke-msg-all",
		"chat_id":     params.ChatIds,
		"mids":        recentIds,
		"except_mids": []string{},
		//"except_mids": params.ExceptMIds,
	}
	contentByte, _ := json.Marshal(sysMsgData)
	data := SendMessageParams{
		Mid:     funcs.CreateMsgId(params.UID),
		Type:    base.MsgContentTypeHollowDel,
		ChatId:  params.ChatIds[0],
		Content: string(contentByte),
		Extra:   "",
	}
	if err = HollowManSendMsg(ctx, data, params.Os); err != nil {
		return []string{}, err
	}
	return []string{}, nil
}

func revokeBackup(ctx context.Context, uid string, chatIds, mIds []string) error {
	logCtx := log.WithFields(ctx, map[string]string{"action": "ClearByChatIds"})
	detailDao := detail.New()
	res, err := make([]detail.Detail, 0), interface{}(nil)
	t := funcs.GetMillis()
	if len(chatIds) > 0 {
		res, err = detailDao.GetMyMsgByChatIds(uid, chatIds)
		if err != nil {
			log.Logger().Error(logCtx, "RevokeBackup saveMany err:", err)
			return errno.Add("RevokeBackup saveMany err", errno.DefErr)
		}
	}
	if len(mIds) > 0 {
		res = detailDao.GetDetails(mIds)
		if err != nil {
			log.Logger().Error(logCtx, "RevokeBackup saveMany err:", err)
			return errno.Add("RevokeBackup saveMany err", errno.DefErr)
		}
	}
	if len(res) == 0 {
		return nil
	}
	addData := make([]revokebackup.RevokeBackup, 0)
	for _, v := range res {
		addData = append(addData, revokebackup.RevokeBackup{
			ID:       v.ID,
			SenderId: v.FromUID,
			Content:  v.Content,
			CreateAt: t,
		})
	}
	_, err = revokebackup.New().SaveMany(addData)
	if err != nil {
		log.Logger().Error(logCtx, "RevokeBackup saveMany err:", err, addData)
		return errno.Add("RevokeBackup saveMany err", errno.SaveDataFailed)
	}
	return nil
}
