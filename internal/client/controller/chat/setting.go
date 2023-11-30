package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imsdk/internal/common/dao/conversation/noticesetting"
	"imsdk/internal/common/model/active"
	"imsdk/internal/common/model/chat"
	"imsdk/pkg/app"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/response"
	"os"
)

type HideRequest = chat.HideRequest
type DeleteChatRequest = chat.DeleteChatRequest
type IsTopRequest = chat.IsTopRequest

func UpdateChatIsTop(ctx *gin.Context) {
	uid, _ := ctx.Get("uid")
	var params IsTopRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	if chat.UpdateIsTop(ctx, uid.(string), params) {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, errno.Add("fail", errno.DefErr))
	return
}

func UpdateChatIsTopV2(ctx *gin.Context) {
	uid, _ := ctx.Get("uid")
	var params IsTopRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	if res, err := chat.UpdateIsTopV2(ctx, uid.(string), params); err == nil {
		response.RespData(ctx, res)
		return
	}
	response.RespErr(ctx, errno.Add("fail", errno.DefErr))
	return
}

func HideChat(ctx *gin.Context) {
	uid, _ := ctx.Get("uid")
	var params HideRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	//os := ctx.Request.Header.Get("DeviceId")
	if chat.HideChat(ctx, uid.(string), params) {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, errno.Add("fail", errno.DefErr))
	return
}

func UpdateChatIsMuteNotify(ctx *gin.Context) {
	uid, _ := ctx.Get("uid")
	var params struct {
		ChatId       string `json:"chat_id" binding:"required"`
		IsMuteNotify uint8  `json:"is_mute_notify" binding:"oneof=0 1"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	if chat.UpdateIsMuteNotify(ctx, uid.(string), params.ChatId, params.IsMuteNotify) {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, errno.Add("fail", errno.DefErr))
	return
}

func UpdateChatBackground(ctx *gin.Context) {
	uid, _ := ctx.Get("uid")
	var params struct {
		ChatId     string                 `json:"chat_id" binding:"required"`
		Background map[string]interface{} `json:"background" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}

	if chat.UpdateBackground(ctx, uid.(string), params.ChatId, params.Background) {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, errno.Add("fail", errno.DefErr))
	return
}

func UpdateChatIsShowName(ctx *gin.Context) {
	uid, _ := ctx.Get("uid")
	var params struct {
		ChatId     string `json:"chat_id" binding:"required"`
		IsShowName uint8  `json:"is_show_name" binding:"oneof=0 1"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}

	if chat.UpdateIsShowName(ctx, uid.(string), params.ChatId, params.IsShowName) {
		response.RespSuc(ctx)
		return
	}
	response.RespErr(ctx, errno.Add("fail", errno.DefErr))
	return
}

func GetChatSetting(ctx *gin.Context) {
	userId, _ := ctx.Get("uid")
	var params struct {
		ChatId string `json:"chat_id" binding:"required,gt=3,lte=64"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	data, err := chat.GetSetting(userId.(string), params.ChatId)
	if err == nil {
		response.RespData(ctx, data)
		return
	}
	response.RespErr(ctx, errno.Add("fail", errno.DefErr))
	return
}

func GetChatSettings(ctx *gin.Context) {
	userId, _ := ctx.Get("uid")
	var params struct {
		ChatIds []string `json:"chat_id"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	if len(params.ChatIds) == 0 {
		response.RespListData(ctx, []string{})
		return
	}
	data, err := chat.GetSettings(userId.(string), params.ChatIds)
	if err == nil {
		response.RespListData(ctx, data)
		return
	}
	response.RespErr(ctx, errno.Add("fail", errno.DefErr))
	return
}

func InitNoticeSetting(ctx *gin.Context) {
	filename := app.Config().GetPublicConfigDir() + "notice_setting.json"
	filePtr, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Open file failed [Err:%s]\n", err.Error())
		return
	}
	defer filePtr.Close()
	var settings []noticesetting.Setting

	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&settings)
	if err != nil {
		fmt.Println("Decoder failed", err.Error())
	} else {
		fmt.Println("Decoder success", len(settings))
		mills := funcs.GetMillis()
		for k, v := range settings {
			settings[k].ID = funcs.Md516(v.Action + v.Lang + v.Role)
			settings[k].CreatedAt = mills
			settings[k].UpdatedAt = mills
		}
		err = noticesetting.New().AddSetting(settings)
		if err != nil {
			response.RespErr(ctx, err)
			return
		}
	}
	response.RespSuc(ctx)
	return
}

func ReportChatActive(ctx *gin.Context) {
	var params active.ChatActiveRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	userId, _ := ctx.Get("uid")
	active.SaveActiveChat(userId.(string), params)
	response.RespSuc(ctx)
	return
}
