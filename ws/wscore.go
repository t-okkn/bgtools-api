package ws

import (
	"fmt"
	"net"
	"net/http"

	"bgtools-api/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsConnection struct {
	*websocket.Conn
}

var (
	// <summary>: WebSocket開始用パラメータ
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// <summary>: WebSocketのRequest用チャンネル
	chWsReq = make(chan models.WsRequest)

	// <summary>: 接続情報のプール
	WsConnPool = map[string]*WsConnection{}

	// <summary>: 部屋情報のプール
	RoomPool = map[string]models.RoomInfoSet{}
)

// <summary>: WebSocket接続時に行われる動作
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
	logp.Method = models.CONNECT

	res := models.WsResponse{
		Method: models.CONNECT.String(),
		Params: models.ConnectedResponse{
			ConnId: connid,
		},
	}

	wsconn.sendJson("CONN", res, logp)

	go readRequests(connid, wsconn)
}

// <summary>: WebSocketでのリクエストを待ち受けます
func ServeRequest() {
	for {
		// メッセージが入るまで、ここでブロック
		e := <-chWsReq
		var action func(models.WsRequest)

		switch models.ParseMethod(e.Method) {
		case models.CREATE:
			action = actionCreate

		case models.JOIN:
			action = actionJoin

		case models.LEAVE:
			action = actionLeave

		case models.BROADCAST:
			action = actionBroadcast

		default:
			action = actionNone
		}

		if action != nil {
			action(e)
		}
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
			// TODO: 他のCloseCodeのときはどうするか検討
			// そもそもどういう状況でどんなCloseCodeになるか要調査
			if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				logp.log("[close-1005] NoStatusReceived: 接続が切断されました")

				if n := deleteConnection(id); n != "" {
					notifyOtherPlayers(n)
				}

				break

			} else {
				logp.IsProcError = true
				logp.log(fmt.Sprintf("メッセージの受信に失敗しました：%s", err))
			}
		}
	}
}

// <summary>: 接続情報を削除します
func deleteConnection(id string) (notify string) {
	notify = deletePlayerInfo(id)
	delete(WsConnPool, id)
	return
}

// <summary>: プレイヤー情報を部屋情報リストから削除します
func deletePlayerInfo(id string) (notify string) {
	notify = ""
	check := ""
	pos := 0

	// TODO: RoomPoolのロック制御しないと、非同期で読み書きし放題は・・・
	for roomid, room := range RoomPool {
		for i, player := range room.Players {
			if player.ConnId == id {
				check = roomid
				pos = i
				break
			}
		}

		if check != "" {
			break
		}
	}

	if check == "" {
		return
	}

	room := RoomPool[check]

	if len(room.Players) <= 1 {
		delete(RoomPool, check)

	} else {
		room.Players[pos] = room.Players[len(room.Players)-1]
		room.Players = room.Players[:len(room.Players)-1]

		RoomPool[check] = room
		notify = check
	}

	return
}

// <summary>: 部屋にいる他のプレイヤーに通知します
func notifyOtherPlayers(roomid string) {
	if roomid == "" {
		return
	}

	room, exist := RoomPool[roomid]
	if !exist {
		return
	}

	min := models.BgScore[room.GameId].MinPlayers

	res := models.WsResponse{
		Method: models.NOTIFY.String(),
		Params: models.RoomResponse{
			IsWait:   len(room.Players) < min,
			RoomInfo: room,
		},
	}

	for _, p := range room.Players {
		conn, ex := WsConnPool[p.ConnId]
		if !ex {
			continue
		}

		logp := newLogParams(p.ConnId, conn.RemoteAddr())
		logp.Method = models.NOTIFY

		conn.sendJson("NTFY", res, logp)
	}
}

// <summary>: JSONデータを送信します
func (conn *WsConnection) sendJson(action string, res models.WsResponse, logp logParams) {
	if err := conn.WriteJSON(res); err == nil {
		logp.log(fmt.Sprintf("<%s> 送信完了：%+v", action, res))

	} else {
		logp.IsProcError = true
		e := fmt.Sprintf("<%s> メッセージの送信に失敗しました：%s", action, err)
		logp.log(e)
	}
}
