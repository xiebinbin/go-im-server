package forward

import (
	"context"
	"errors"
	json "github.com/json-iterator/go"
	"imsdk/internal/client/controller/ws"
	"imsdk/internal/client/model/user/device"
	"imsdk/internal/common/model/socket"
	"imsdk/internal/common/pkg/base"
	"imsdk/internal/common/pkg/req/request"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/redis"
	"os"
)

type ReceiveUser struct {
	Uid       string   `json:"uid"`
	DeviceIds []string `json:"device_ids"`
}

type SocketData struct {
	Cmd  string      `json:"cmd,omitempty"`
	Data interface{} `json:"items,omitempty"`
}

type WebsocketParams struct {
	SocketData  SocketData    `json:"socket_data"`
	ReceiveUser []ReceiveUser `json:"receive_user,omitempty"`
}

type SenderInfo struct {
	SenderUid      string `json:"sender_uid"`
	SenderDeviceId string `json:"sender_device_id"`
}

// PushMessageToUserSocket
// The first attempt to use the socket to forward messages to the user, if unsuccessful, join the offline queue
func PushMessageToUserSocket(ctx context.Context, pushParams PushMessageParams, isMq bool) (offlineDeviceIds []string, err error) {
	senderInfo := pushParams.SenderInfo
	receiveUid := pushParams.Uid
	cmd := pushParams.Cmd
	data := pushParams.Data
	logCtx := log.WithFields(ctx, map[string]string{
		"cmd":              cmd,
		"action":           "socketPush",
		"sender":           senderInfo.SenderUid,
		"sender_device_id": senderInfo.SenderDeviceId,
		"receiveUid":       receiveUid,
	})
	// get receiver online device socket connect
	isSender := senderInfo.SenderUid == receiveUid
	connections := socket.GetUserConnections(receiveUid)
	log.Logger().Info(logCtx, "PushMessageToUserSocket: rece: ", receiveUid, ", conn: ", connections)
	userDevices := make(map[string]string)
	if len(connections) <= 0 {
		// No need to forward offline messages to the sender
		if isSender && data.Type != base.MsgTypeMeeting {
			return
		}
		return
	}
	//log.Logger().Info(logCtx, "PushMessageToUserSocket: rece: ", receiveUid, ", user dev: ", userDevices)
	//if len(userDevices) == 0 {
	//	log.Logger().Error(logCtx, "receiver's device not exit")
	//	return nil, errors.New("device not exit")
	//}

	thisSocketHost := os.Getenv("SOCKET_HOST") // docker run -e

	//log.Logger().Info(logCtx, "PushMessageToUserSocket: rece: ", receiveUid, ", user dev: ", userDevices, " , host: ", thisSocketHost)
	connMap := make(map[string][]string) // map[socketHost][]{"deviceId"}
	for _, conn := range connections {
		if senderInfo.SenderUid == receiveUid && conn.DeviceId == senderInfo.SenderDeviceId { // do not forward message to sender device again
			continue
		}
		isFront := device.GetDeviceStatus(receiveUid, conn.DeviceId) == 1
		log.Logger().Info(logCtx, "PushMessageToUserSocket isFront:", isFront, connMap, " receiveUid : ", receiveUid)
		if conn.DeviceId == base.AppDeviceId && !isFront {
			continue
		}
		if _, ok := connMap[conn.Host]; ok {
			connMap[conn.Host] = append(connMap[conn.Host], conn.DeviceId)
		} else {
			connMap[conn.Host] = []string{conn.DeviceId}
		}
	}
	socketData := SocketData{
		Cmd:  cmd,
		Data: data,
	}
	// todo remember to delete this after two version
	oriSocketData := ""
	if cmd == CmdApplicationCmd {
		socketData.Data = data.Content
		oriSocketData = data.Content
	}
	log.Logger().Info(logCtx, "PushMessageToUserSocket:-2 : ", socketData, connMap, isMq)
	for host, deviceIds := range connMap {
		if isMq {
			// todo must update
			failedDeviceIds, er := PushMessageToUserSocketWithApi(ctx, host, receiveUid, deviceIds, socketData, oriSocketData)
			if er != nil { // all failed
				continue
			}
			sucDeviceIds := funcs.DifferenceSetString(deviceIds, failedDeviceIds)
			if len(sucDeviceIds) != 0 {
				for _, id := range sucDeviceIds {
					delete(userDevices, id)
				}
			}

		} else if host == thisSocketHost {
			log.Logger().Info(logCtx, "PushMessageToUserSocket:-3 equal: ", socketData, receiveUid, deviceIds, host, thisSocketHost)
			sucDeviceIds, _ := ws.PushMsgToClient(ctx, ws.PushMessageToClientParams{
				Address: receiveUid,
				Data:    socketData,
				OriData: oriSocketData,
				Devices: deviceIds,
			})
			if len(sucDeviceIds) > 0 && !isSender {
				for _, id := range sucDeviceIds {
					delete(userDevices, id)
				}
			}
		} else {
			pushParams.Host = host
			for _, id := range deviceIds {
				isFront := device.GetDeviceStatus(receiveUid, id) == 1
				// todo singapore
				version := redis.Client.Get("user:version:" + receiveUid).Val()
				if funcs.CompareVersion(version, "1.1.0") == -1 {
					isFront = true
				}
				if id == base.AppDeviceId && !isFront {
					continue
				}
				pushParams.Devices = []string{id}
				if err = ProduceSocketMessage(ctx, pushParams); err != nil {
					//connect to other server , produce data to mq
					log.Logger().Error(logCtx, "failed to produce message,  PushMessageToUserSocket-4, err: ", err)
					return nil, errors.New("failed to send message")
				}
				delete(userDevices, id)
			}
		}
	}
	if isSender {
		return nil, nil
	}
	log.Logger().Info(logCtx, "PushMessageToUserSocket:-5 userDevices:", userDevices)
	if len(userDevices) == 0 {
		return
	}
	for id, _ := range userDevices {
		offlineDeviceIds = append(offlineDeviceIds, id)
	}
	return offlineDeviceIds, nil
}

// PushMessageToUserSocketWithApi
// return: forward failed device ids, error
func PushMessageToUserSocketWithApi(ctx context.Context, host, receiveUid string, deviceIds []string, socketData SocketData, oriSocketData string) ([]string, error) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "PushMessageToUserSocketWithApi"})
	host = host + "/pushMessageToClient"
	params := ws.PushMessageToClientParams{
		Address: receiveUid,
		Data:    socketData,
		OriData: oriSocketData,
		Devices: deviceIds,
	}
	var resp struct {
		Code  int      `json:"code"`
		Data  []string `json:"data"` // forward failed device ids
		Msg   string   `json:"msg"`
		Field string   `json:"field"`
	}
	response, er := request.InnerReq(ctx, host, params)
	if er != nil {
		return nil, er
	}
	err := json.Unmarshal(response, &resp)
	log.Logger().Info(logCtx, "pushMessageToClient:", receiveUid, response, er)
	if err != nil {
		log.Logger().Error(logCtx, "failed to decode response data , err: ", err)
		return nil, err
	}
	return resp.Data, nil
}

func PushMessageToUserSocketDirectly(ctx context.Context, pushParams PushMessageParams, isMq bool) error {
	if pushParams.ReqId == "" {
		if ReqId := ctx.Value(base.HeaderFieldReqId); ReqId != nil {
			pushParams.ReqId = ReqId.(string)
		}
	}
	logCtx := log.WithFields(ctx, map[string]string{"action": "PushMessageToUserSocketDirectly"})
	offlineDeviceIds, err := PushMessageToUserSocket(ctx, pushParams, isMq)
	log.Logger().Info(logCtx, "uid: ", pushParams.Uid, " after  PushMessageToUserSocket, offline device ids : ", offlineDeviceIds, err, ", no push offline :", pushParams.NoPushOffline)
	if err != nil {
		log.Logger().Error(logCtx, "uid: ", pushParams.Uid, " failed to forward msg to user with socket , err: ", err)
		return errno.Add("fail", errno.DefErr)
	}
	if len(offlineDeviceIds) == 0 || pushParams.NoPushOffline {
		//if pushParams.NoPushOffline {
		return nil
	}
	pushParams.Devices = offlineDeviceIds
	//val, _ := json.Marshal(pushParams)
	//ProcessOfflineMsg(ctx, val)
	//return nil
	//err = ProduceMessageToOfflineMq(ctx, pushParams)
	//if err != nil {
	//	log.Logger().Error(logCtx, "uid: ", pushParams.Uid, " failed to produce msg to offline mq , err: ", err)
	//	return errno.Add("fail", errno.DefErr)
	//}
	return err
}
