package forward

import (
	"context"
	"encoding/json"
	"fmt"
	"imsdk/internal/client/controller/ws"
	"imsdk/internal/client/model/user/device"
	"imsdk/internal/common/model/socket"
	"imsdk/internal/common/pkg/base"
)

const (
	CmdApplicationCmd        = "_application_cmd"
	CmdSdkContent            = "_sdk_content"
	CmdNoticeNewMsg          = "CHAT_MSG"
	CmdApplyAddFriend        = "FRIEND_APPLY"
	CmdAgreeAddFriend        = "agree_add_friend"
	CmdNoticeMeetingMember   = "meeting_member"
	CmdNoticeMeetingProfile  = "meeting_profile"
	CmdNoticeMeetingInvite   = "meeting_invite"
	CmdSyncOffline           = "sync_offline"
	CmdSyncMuteSetting       = "sync_setting_mute"
	CmdSyncHideSetting       = "sync_setting_hide"
	CmdSyncSettingIsTop      = "sync_setting_top"
	CmdAddBlockUser          = "add_block_user"
	CmdCancelBlockUser       = "cancel_block_user"
	CmdDeleteUser            = "delete_user"
	CmdCardMessageStatus     = "card_message_status"
	CloseTagAnotherDevice    = 1
	CloseTagAnotherAppDevice = 2
	CloseTagAppLoginOut      = 3
)

type SendCmdParams struct {
	UIds          []string       `json:"uids"`
	Data          string         `json:"data"`
	OfflineData   OfflineRequest `json:"offline"` // {"title":"","body":""}
	NoPushOffline bool           `json:"no_push_offline"`
	AppCtx        string         `json:"app_ctx"`
}

type OldCmdInfo struct {
	Cmd   string                 `json:"cmd"`
	Items map[string]interface{} `json:"items"`
}

func SendCmdMsg(ctx context.Context, params SendCmdParams) error {
	// todo old->new pc cmd start======
	var oldCmdInfo OldCmdInfo
	er := json.Unmarshal([]byte(params.Data), &oldCmdInfo)
	if er != nil {
		fmt.Println("SendCmdMsg er:", er)
	}
	// todo old->new pc cmd end======
	for _, uid := range params.UIds {
		if oldCmdInfo.Cmd == CmdSyncSettingIsTop || oldCmdInfo.Cmd == CmdSyncMuteSetting || oldCmdInfo.Cmd == CmdSyncHideSetting {
			SendSyncChatSettingMsgOld(ctx, oldCmdInfo.Cmd, uid, oldCmdInfo.Items)
			continue
		}
		data := PushMessageParams{
			Cmd: CmdApplicationCmd,
			Uid: uid,
			Data: MsgItem{
				Content: params.Data,
			},
			OfflineData:   params.OfflineData,
			NoPushOffline: params.NoPushOffline,
			AppCtx:        params.AppCtx,
			ReqId:         ctx.Value(base.HeaderFieldReqId).(string),
		}
		err := PushMessageToUserSocketDirectly(ctx, data, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func SendSyncChatSettingMsgOld(ctx context.Context, cmd, uid string, msgData map[string]interface{}) bool {
	connections := socket.GetUserConnections(uid)
	connMap := make(map[string][]string) // map[socketHost][]{"deviceId"}
	for _, conn := range connections {
		isFront := device.GetDeviceStatus(uid, conn.DeviceId) == 1
		if conn.DeviceId == base.AppDeviceId && !isFront {
			continue
		}
		if _, ok := connMap[conn.Host]; ok {
			connMap[conn.Host] = append(connMap[conn.Host], conn.DeviceId)
		} else {
			connMap[conn.Host] = []string{conn.DeviceId}
		}
	}
	contentByte, _ := json.Marshal(msgData)
	socketData := SocketData{
		Cmd:  cmd,
		Data: string(contentByte),
	}
	for _, deviceIds := range connMap {
		sucDeviceIds, _ := ws.PushMsgToClient(ctx, ws.PushMessageToClientParams{
			Address: uid,
			Data:    socketData,
			Devices: deviceIds,
		})
		fmt.Println("sucDeviceIds:", sucDeviceIds)
	}
	return true
}
