package ws

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"imsdk/internal/common/model/active"
	"imsdk/internal/common/model/socket"
	"imsdk/pkg/log"
	"net"
	"time"
)

type ClientInfo struct {
	UID      string
	Address  string
	DeviceId string
	ReqId    string
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	clientInfo ClientInfo

	wsType socket.WebsocketType

	aliveTime int64
}

// readMsg pumps messages from the websocket connection to the hub.
//
// The application runs readMsg in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (cli *Client) readMsg() {
	//defer wg.Done()
	logFields := map[string]string{
		"action": "readMsg",
		"uid":    cli.clientInfo.Address,
	}
	ctx := log.WithFields(context.Background(), logFields)
	releaseReason := "read message end"
	defer func() {
		releaseResource(ctx, cli, releaseReason)
	}()
	cli.conn.SetReadLimit(maxMessageSize)
	cli.conn.SetReadDeadline(time.Now().Add(pongWait))

	// client send and server response
	cli.conn.SetPingHandler(func(appData string) error {
		log.Logger().Info(ctx, "ping handler ", " alive time: ", cli.aliveTime)
		if _, ok := userClient[cli.clientInfo.Address+cli.clientInfo.DeviceId]; !ok {
			return errors.New("exception connection")
		}
		//cli.aliveTime = time.Now().Unix()

		nowTime := time.Now().Unix()
		userClient[cli.clientInfo.Address+cli.clientInfo.DeviceId].aliveTime = nowTime
		active.SetOnlineByUid(cli.clientInfo.Address, nowTime)

		err := cli.conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(writeWait))
		if errors.Is(err, websocket.ErrCloseSent) {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	// server send and client response
	cli.conn.SetPongHandler(func(string) error {
		log.Logger().Info(ctx, "pong handler, ", " alive time: ", cli.aliveTime)
		cli.conn.SetReadDeadline(time.Now().Add(pongWait))
		//cli.aliveTime = time.Now().Unix()

		nowTime := time.Now().Unix()
		userClient[cli.clientInfo.Address+cli.clientInfo.DeviceId].aliveTime = nowTime
		active.SetOnlineByUid(cli.clientInfo.Address, nowTime)
		return nil
	})

	cli.conn.SetCloseHandler(func(code int, text string) error {
		log.Logger().Info(ctx, "close handler,", " alive time: ", cli.aliveTime)
		releaseReason = "close handle"
		return nil
	})

	for {
		_, message, err := cli.conn.ReadMessage()
		log.Logger().Info(ctx, "data is : ", string(message))
		if err != nil {
			releaseReason = "read message error"
			log.Logger().Info(ctx, "reader error: ", err, " alive time: ", cli.aliveTime)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				releaseReason = "read message abnormal error"
			}
			return
		}
		//cli.aliveTime = time.Now().Unix()
		userClient[cli.clientInfo.Address+cli.clientInfo.DeviceId].aliveTime = time.Now().Unix()
		// handle message
		//if message != nil {
		//}
	}
	return
}

// writeMsg pumps messages from the hub to the websocket connection.
//
// A goroutine running writeMsg is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (cli *Client) writeMsg() {
	//defer wg.Done()
	logFields := map[string]string{
		"action":    "writeMsg",
		"device-id": cli.clientInfo.DeviceId,
	}
	ctx := log.WithFields(context.Background(), logFields)
	ticker := time.NewTicker(pingPeriod)
	reason := "write message end"
	defer func() {
		ticker.Stop()
		closeSocketConnection(ctx, cli, websocket.CloseNormalClosure, "", reason)
		releaseResource(ctx, cli, reason)
	}()
	log.Logger().Info(ctx, "start, ", " alive time: ", cli.aliveTime)
	t := time.Now().Unix()
	for {
		select {
		case message, ok := <-cli.send:
			if !ok { // closed the channel.
				return
			}
			log.Logger().Info(ctx, "start write, ", " message: ", string(message), " alive time: ", cli.aliveTime)

			cli.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := cli.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Logger().Error(ctx, "write error: ", err, " alive time: ", cli.aliveTime)
				reason = "write message error"
				return
			}
			seq, er := w.Write(message)
			log.Logger().Info(ctx, "socket write message result: seq is: ", seq, " , error is:", er, " alive time: ", cli.aliveTime)

			if err = w.Close(); err != nil {
				log.Logger().Error(ctx, "write pump close error: ", err, " alive time: ", cli.aliveTime)
				reason = "failed to close write pump"
				return
			}
			//cli.aliveTime = time.Now().Unix()
			userClient[cli.clientInfo.Address+cli.clientInfo.DeviceId].aliveTime = time.Now().Unix()
		case <-ticker.C:
			cli.conn.SetWriteDeadline(time.Now().Add(writeWait))
			t = time.Now().Unix()
			aliveT := userClient[cli.clientInfo.Address+cli.clientInfo.DeviceId].aliveTime
			if t-aliveT > 2*heartbeatInterval { // didn't received message in twice pongwait, and close connect
				log.Logger().Info(ctx, "expired , now: ", t, " alive time: ", cli.aliveTime, " did not receive message in twice pong wait")
				reason = "did not receive message in twice pong wait"
				return
			}
			if err := cli.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Logger().Error(ctx, "heart beat error: ", err, " alive time: ", cli.aliveTime)
			}
		}
	}
	return
}
