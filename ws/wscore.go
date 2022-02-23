package ws

import (
	"fmt"
	"net/http"

	"bgtools-api/models"

	"github.com/google/uuid"
)

// <summary>: Websocket接続時に行われる動作
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	logp := newLogParams()
	logp.ClientIP = r.RemoteAddr
	logp.Message = "クライアントからの新規接続要求"
	logp.log()

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logp.IsProcError = true
		logp.Message = fmt.Sprintf("WebSocketのUpgradeに失敗しました：%s", err)
		logp.log()

		return
	}

	wsconn := &WsConnection{conn}
	obj, _ := uuid.NewRandom()
	uuid_str := obj.String()

	WsConnPool[uuid_str] = wsconn
	logp.ConnId = uuid_str
	logp.Method = models.CONNCTED

	res := models.WsResponse{
		Method: models.CONNCTED.String(),
		Params: models.ConnectedResponse{
			ConnId: uuid_str,
		},
	}

	if err := conn.WriteJSON(res); err == nil {
		logp.Message = fmt.Sprintf("<CONNECTED> 送信完了：%+v", res)
		logp.log()

	} else {
		logp.IsProcError = true
		logp.Message =
			fmt.Sprintf("<CONNECTED> メッセージの送信に失敗しました：%s", err)
		logp.log()
	}

	go readRequests(uuid_str, wsconn)
}

// <summary>: Websocketでの電文のやり取りを行います
func ListenAndServe() {
	for {
		// メッセージが入るまで、ここでブロック
		e := <-chWsReq
		var action func(models.WsRequest)

		switch models.ParseMethod(e.Method) {
		case models.CREATE_ROOM:
			action = actionCreateRoom

		default:
			action = actionNone
		}

		action(e)
	}
}

// <summary>: 受信した内容を読み取ります
func readRequests(id string, conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			elogp := newLogParams()
			elogp.ClientIP = conn.RemoteAddr().String()
			elogp.ConnId = id

			deleteConnection(id)

			elogp.IsProcError = true
			elogp.Message = fmt.Sprintf("予期せぬエラーが発生しました：%s", r)
			elogp.log()
		}
	}()

	var req models.WsRequest

	for {
		logp := newLogParams()
		logp.ClientIP = conn.RemoteAddr().String()
		logp.ConnId = id

		if err := conn.ReadJSON(&req); err == nil {
			logp.Method = models.ParseMethod(req.Method)
			logp.Message = fmt.Sprintf("メッセージ受信：%+v", req)
			logp.log()

			chWsReq <- req

		} else {
			logp.IsProcError = true
			logp.Message =
				fmt.Sprintf("メッセージの受信に失敗しました：%s", err)
			logp.log()
		}
	}
}

// <summary>: 接続情報を削除します
func deleteConnection(id string) {
	check := ""

	for roomid, room := range RoomPool {
		_, exsit := room.Players[id]

		if exsit {
			delete(room.Players, id)

			if len(room.Players) == 0 {
				check = roomid
			}

			break
		}
	}

	if check != "" {
		delete(RoomPool, check)
	}

	delete(WsConnPool, id)
}

