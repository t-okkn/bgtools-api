package ws

import (
	"fmt"
	"net/http"

	"bgtools-api/models"

	"github.com/gorilla/websocket"
)

// <summary>: プレイヤーの接続情報をまとめた構造体
type PlayerConn struct {
	C      *websocket.Conn
	RoomId string
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

	// <summary>: プレイヤーの接続情報プール
	PlayerPool = NewPlayerMap()

	// <summary>: 部屋情報プール
	RoomPool = NewRoomMap()
)

// <summary>: WebSocket接続時に行われる動作
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	connid, err := getConnId(r.RemoteAddr)
	logp := newLogParams(connid)

	if err != nil {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("不正な接続元からのアクセスです: %s", err))

		return
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("WebSocketのUpgradeに失敗しました: %s", err))

		return
	}

	pconn := PlayerConn{
		C:      conn,
		RoomId: "",
	}
	PlayerPool.Set(connid, pconn)

	logp.ConnId = connid
	logp.Method = models.CONNECT
	logp.Prefix = fmt.Sprintf("<%s>", models.CONNECT.String())

	res := models.WsResponse{
		Method: models.CONNECT.String(),
		Params: models.ConnectResponse{
			ConnId: connid,
		},
	}

	pconn.sendJson(res, logp)

	go readRequests(connid, pconn)
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
func readRequests(id string, pc PlayerConn) {
	defer func() {
		if r := recover(); r != nil {
			elogp := newLogParams(id)
			elogp.IsProcError = true

			deleteConnection(id)
			elogp.log(fmt.Sprintf("予期せぬエラーが発生しました: %s", r))
		}
	}()

	var req models.WsRequest

	for {
		logp := newLogParams(id)

		if err := pc.C.ReadJSON(&req); err == nil {
			logp.Method = models.ParseMethod(req.Method)
			logp.log(fmt.Sprintf("メッセージ受信: %+v", req))

			if !isCorrectConnId(req.ConnId, pc.C.RemoteAddr().String()) {
				pc.sendError(models.ErrIllegalConnId, logp)
				continue
			}

			chWsReq <- req

		} else {
			// TODO: 他のCloseCodeのときはどうするか検討
			// そもそもどういう状況でどんなCloseCodeになるか要調査
			if websocket.IsCloseError(err, websocket.CloseNoStatusReceived) {
				logp.Method = models.DISCONNECT
				logp.Prefix = "[close-1005]"
				logp.log("NoStatusReceived: 接続が切断されました")

				if n := deleteConnection(id); n != "" {
					notifyOtherPlayers(n)
				}

				break

			} else {
				logp.IsProcError = true
				logp.log(fmt.Sprintf("メッセージの受信に失敗しました: %s", err))
			}
		}
	}
}

// <summary>: 接続情報を削除します
func deleteConnection(id string) (notify string) {
	notify = deletePlayerInfo(id)
	PlayerPool.Delete(id)

	return
}

// <summary>: プレイヤー情報を部屋情報プールから削除します
func deletePlayerInfo(id string) (notify string) {
	notify = ""

	pc, ok := PlayerPool.Get(id)
	if !ok {
		return
	}

	room, ok := RoomPool.Get(pc.RoomId)
	if !ok {
		return
	}

	if len(room.Players) <= 1 {
		RoomPool.Delete(pc.RoomId)

	} else {
		for i, player := range room.Players {
			if player.ConnId == id {
				room.Players[i] = room.Players[len(room.Players)-1]
				break
			}
		}

		room.Players = room.Players[:len(room.Players)-1]
		RoomPool.Set(pc.RoomId, room)
		notify = pc.RoomId
	}

	PlayerPool.SetRoomId(id, "")
	return
}

// <summary>: 部屋にいる他のプレイヤーに通知します
func notifyOtherPlayers(roomid string) {
	if roomid == "" {
		return
	}

	room, exist := RoomPool.Get(roomid)
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
		pc, ex := PlayerPool.Get(p.ConnId)
		if !ex {
			continue
		}

		logp := newLogParams(p.ConnId)
		logp.Method = models.NOTIFY
		logp.Prefix = fmt.Sprintf("<%s>", models.NOTIFY.String())

		pc.sendJson(res, logp)
	}
}

// <summary>: エラー内容を送信します
func (pc PlayerConn) sendError(err models.ErrorMessage, logp logParams) {
	logp.Method = models.ERROR
	res := models.WsResponse{
		Method: models.ERROR.String(),
		Params: err,
	}

	pc.sendJson(res, logp)
}

// <summary>: JSONデータを送信します
func (pc PlayerConn) sendJson(res models.WsResponse, logp logParams) {
	if err := pc.C.WriteJSON(res); err == nil {
		logp.log(fmt.Sprintf("送信完了: %+v", res))

	} else {
		logp.IsProcError = true
		e := fmt.Sprintf("メッセージの送信に失敗しました: %s", err)
		logp.log(e)
	}
}
