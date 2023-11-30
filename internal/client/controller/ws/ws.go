package ws

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"imsdk/internal/common/model/socket"
	"imsdk/internal/common/model/user/device"
	"imsdk/internal/common/pkg/req/request"
	"imsdk/pkg/errno"
	"imsdk/pkg/funcs"
	"imsdk/pkg/log"
	"imsdk/pkg/response"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	addr string
	//wg   sync.WaitGroup
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	heartbeatInterval = 30

	// Time allowed to read the next pong message from the peer.
	pongWait = heartbeatInterval * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	closeBufSize = 512

	closeCodeTokenErr           = 401
	closeTextTokenErr           = "token-err"
	closeCodeAnotherDeviceLogin = 1101
	closeTextAnotherDeviceLogin = "another-device-login"

	tokenErrCode = 401
	tokenErrText = "token error"
	sysErrCode   = 1102
	sysErrText   = "fail"
)

var (
	userClient = make(map[string]*Client) //
	wsUpGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type PushMessageToClientParams struct {
	Address string      `json:"address"`
	Data    interface{} `json:"data"`
	OriData string      `json:"ori_data"`
	Devices []string    `json:"devices"`
}

func Connect(ctx *gin.Context) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "wsConnect"})
	address := ctx.Query("address")
	deviceId := ctx.Query("device_id")
	if address == "" || deviceId == "" {
		return
	}
	if reqId := ctx.Query("req-id"); reqId == "" {
		reqId = funcs.UniqueId16()
	}
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			//log.Logger().Error(logCtx, "panic error : ", string(buf))
		}
	}()
	conn, err := wsUpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Logger().Info(logCtx, " wsUpGrader err : ", address)
		return
	}

	clientInfo := ClientInfo{
		UID:      address,
		Address:  address,
		DeviceId: deviceId,
	}
	socketHost := os.Getenv("SOCKET_HOST") // docker run -e
	// Determine whether the terminal has established a long connection on other devices （login on the other device）
	// if login , disconnect
	host, err := socket.GetSocketConnectionHost(clientInfo.Address, clientInfo.DeviceId)
	fmt.Println("clientInfo--> ", clientInfo, socketHost)
	if err != nil {
		log.Logger().Error(logCtx, "failed to get old socket host , err: ", err)
		ctx.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		closeByConn(ctx, conn, websocket.CloseInternalServerErr, "unknown error ")
		return
	}

	if host != "" {
		if host == socketHost { // login on this server
			if cli, ok := userClient[clientInfo.Address+clientInfo.DeviceId]; ok {
				closeSocketConnection(logCtx, cli, closeCodeAnotherDeviceLogin, closeTextAnotherDeviceLogin, "another device login")
				releaseResource(logCtx, cli, "another device login")
			}
		} else { // login on other server
			requestCloseConnection(logCtx, host, clientInfo.Address, clientInfo.DeviceId)
		}
	}
	client := &Client{
		conn:       conn,
		send:       make(chan []byte, 256),
		clientInfo: clientInfo,
		aliveTime:  time.Now().Unix(),
	}
	log.Logger().Info(logCtx, " client info : ", client)
	if _, err = socket.SaveSocketConnection(clientInfo.Address, clientInfo.DeviceId, socketHost, clientInfo.ReqId); err != nil {
		log.Logger().Error(logCtx, " failed to save socket connect ,err : ", err)
		ctx.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		closeByConn(ctx, conn, websocket.CloseInternalServerErr, "unknown error ")
		return
	}
	userClient[clientInfo.Address+clientInfo.DeviceId] = client
	err = device.Add(ctx, device.AddRequest{
		UID:      clientInfo.Address,
		DeviceId: clientInfo.DeviceId,
		OS:       "",
	})
	if err != nil {
		return
	}
	log.Logger().Info(logCtx, "client connect successfully")
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writeMsg()
	go client.readMsg()
}

// PushMessageToClient
// api : request from other server by http
func PushMessageToClient(ctx *gin.Context) {
	var params PushMessageToClientParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}
	_, failedDeviceIds := PushMsgToClient(ctx, params)
	//response.RespData(ctx, failedDeviceIds)
	fmt.Println("failedDeviceIds:", failedDeviceIds)
	return
}

// PushMsgToClient execute request
func PushMsgToClient(ctx context.Context, params PushMessageToClientParams) (sucDeviceIds, failedDeviceIds []string) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "PushMsgToClient", "uid": params.Address})
	log.Logger().Info(logCtx, "device : ", params.Devices)
	for _, deviceId := range params.Devices {
		fmt.Println("params.Address+deviceId:", params.Address+deviceId, userClient)
		cli, ok := userClient[params.Address+deviceId]
		log.Logger().Info(logCtx, "single device :", " ,device : ", deviceId, cli)
		if !ok {
			failedDeviceIds = append(failedDeviceIds, deviceId)
			log.Logger().Error(logCtx, " user client info not exist  ", ", device-id: ", deviceId)
			continue
		}
		msg, _ := json.Marshal(params.Data)
		log.Logger().Info(logCtx, " start-forward: ", " deviceId: ", deviceId, " data: ", params.Data, " alive time: ", cli.aliveTime)
		if time.Now().Unix()-cli.aliveTime >= int64(pongWait*2/time.Second) { // didn't received message in twice pongwait, and close connect
			reason := "did not receive message in twice pong wait"
			closeSocketConnection(ctx, cli, websocket.CloseNormalClosure, "", reason)
			releaseResource(logCtx, cli, reason)
			failedDeviceIds = append(failedDeviceIds, deviceId)
			log.Logger().Error(logCtx, " did not receive message in twice pong wait , device-id: ", deviceId)
			continue
		}
		cli.send <- msg
		sucDeviceIds = append(sucDeviceIds, deviceId)
	}
	return
}

func requestCloseConnection(ctx context.Context, host, uid, deviceId string) {
	data := map[string]string{
		"uid":       uid,
		"device_id": deviceId,
	}
	request.PostJson(ctx, "http://"+host+addr+"/ws/closeDeviceConnRequest", data)
}

func closeSocketConnection(ctx context.Context, cli *Client, closeCode int, text, reason string) {
	if text == "" {
		text = reason
	}
	log.Logger().Info(ctx, "reason: ", reason, ", uid: ", cli.clientInfo.Address, " device-id: ", cli.clientInfo.DeviceId)
	closeByConn(ctx, cli.conn, closeCode, text)
	return
}

func closeByConn(ctx context.Context, conn *websocket.Conn, code int, text string) {
	logCtx := log.WithFields(ctx, map[string]string{"action": "closeByConn"})
	defer func() {
		if err := recover(); err != nil {
			log.Logger().Error(logCtx, " panic error ", err)
		}
	}()
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, text))
	conn.Close()
	return
}

func releaseResource(ctx context.Context, cli *Client, reason string) {
	log.WithFields(ctx, map[string]string{"desc": "releaseResource"})
	log.Logger().Info(ctx, "reason: ", reason, ", uid: ", cli.clientInfo.Address, " device-id: ", cli.clientInfo.DeviceId)

	defer func() {
		if err := recover(); err != nil {
			log.Logger().Error(ctx, "panic error: ", err)
		}
	}()
	if _, ok := userClient[cli.clientInfo.Address+cli.clientInfo.DeviceId]; ok {
		delete(userClient, cli.clientInfo.Address+cli.clientInfo.DeviceId)
	}

	socket.DelSocketConnection(ctx, cli.clientInfo.Address, cli.clientInfo.DeviceId, cli.clientInfo.ReqId)
	if _, ok := <-cli.send; ok {
		close(cli.send)
	}

}

func CloseDeviceConnRequest(ctx *gin.Context) {
	var params struct {
		Uid      string `json:"uid"`
		DeviceId string `json:"device_id"`
	}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.RespErr(ctx, errno.Add("params-err", errno.ParamsErr))
		return
	}

	fields := map[string]string{"action": "CloseDeviceConnRequest", "uid": params.Uid, "device": params.DeviceId}
	logCtx := log.WithFields(ctx, fields)
	log.Logger().Info(logCtx, " start")
	if cli, ok := userClient[params.Uid+params.DeviceId]; ok {
		closeSocketConnection(logCtx, cli, closeCodeAnotherDeviceLogin, closeTextAnotherDeviceLogin, "another device login: http req")
		releaseResource(logCtx, cli, "another device login: http req")
	}
}
