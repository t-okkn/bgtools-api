package ws

import (
	"fmt"
	"net"
	"net/http"

	"bgtools-api/db"
	"bgtools-api/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsConnection struct {
	*websocket.Conn
}

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool {
			return true
		},
	}

	chWsReq = make(chan models.WsRequest)

	// <summary>: 接続情報のプール
	WsConnPool = map[string]*WsConnection{}

	// <summary>: 部屋情報のプール
	RoomPool = map[string]models.RoomInfoSet{}

	// <summary>: DB接続のコネクションプール
	BgRepo *db.BgRepository
)

// <summary>: Websocket接続時に行われる動作
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	h, _, _ := net.SplitHostPort(r.RemoteAddr)
	logp := logParams{
		ClientIP:    h,
		ConnId:      "",
		Method:      models.NONE,
		IsProcError: false,
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("WebSocketのUpgradeに失敗しました：%s", err))

		return
	}

	wsconn := &WsConnection{conn}
	obj, _ := uuid.NewRandom()
	connid := fmt.Sprintf("ws.%s", obj.String())

	WsConnPool[connid] = wsconn
	logp.ConnId = connid
	logp.Method = models.CONNCTED

	res := models.WsResponse{
		Method: models.CONNCTED.String(),
		Params: models.ConnectedResponse{
			ConnId: connid,
		},
	}

	if err := conn.WriteJSON(res); err == nil {
		logp.log(fmt.Sprintf("<CONNECTED> 送信完了：%+v", res))

	} else {
		logp.IsProcError = true
		logp.log(fmt.Sprintf(
			"<CONNECTED> メッセージの送信に失敗しました：%s", err))
	}

	go readRequests(connid, wsconn)
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
			elogp := newLogParams(id, conn.RemoteAddr())
			elogp.IsProcError = true

			deleteConnection(id)
			elogp.log(fmt.Sprintf("予期せぬエラーが発生しました：%s", r))
		}
	}()

	var req models.WsRequest

	for {
		logp := newLogParams(id, conn.RemoteAddr())

		if err := conn.ReadJSON(&req); err == nil {
			req.ClientIP = conn.RemoteAddr()
			logp.Method = models.ParseMethod(req.Method)
			logp.log(fmt.Sprintf("メッセージ受信：%+v", req))

			chWsReq <- req

		} else {
			logp.IsProcError = true
			logp.log(fmt.Sprintf("メッセージの受信に失敗しました：%s", err))

			// TODO: 他のCloseCodeのときはどうするか検討
			// そもそもどういう状況でどんなCloseCodeになるか要調査
			if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				deleteConnection(id)

				logp.IsProcError = false
				logp.log("接続が切断されました")

				break
			}
		}
	}
}

// <summary>: 接続情報を削除します
func deleteConnection(id string) {
	check := ""

	// TODO: RoomPoolのロック制御しないと、非同期で読み書きし放題は・・・
	for roomid, room := range RoomPool {
		_, exsit := room.Players[id]

		if exsit {
			delete(room.Players, id)

			if len(room.Players) == 0 {
				check = roomid
			}
			// TODO: 部屋にいる人に通知が必要
			// TODO: また、最小プレー人数を下回ったら待機状態にすべき

			break
		}
	}

	if check != "" {
		delete(RoomPool, check)
	}

	delete(WsConnPool, id)
}
