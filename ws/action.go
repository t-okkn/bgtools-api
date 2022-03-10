package ws

import (
	"fmt"

	"bgtools-api/models"
)

// <summary>: [Method] CREA に関する動作を定義します
func actionCreate(req models.WsRequest) {
	action := models.CREATE.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	conn, ok := ConnPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	var response models.WsResponse

	// リクエストされた部屋情報がなければエラー
	if _, exist := RoomPool.Get(req.RoomId); exist {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrRoomExisted,
		}

		conn.sendJson(action, response, logp)
		return
	}

	data, exist := models.BgScore[req.GameId]

	// リクエストされたボードゲーム情報がなければエラー
	if !exist {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrBoardgameNotFound,
		}

		conn.sendJson(action, response, logp)
		return
	}

	players := make([]models.PlayerInfoSet, 1, data.MaxPlayers)
	players[0] = models.PlayerInfoSet{
		ConnId:      req.ConnId,
		PlayerColor: req.PlayerColor,
	}

	room := models.RoomInfoSet{
		GameId:  req.GameId,
		Players: players,
	}
	RoomPool.Set(req.RoomId, room)

	logp.Method = models.OK
	response = models.WsResponse{
		Method: models.OK.String(),
		Params: models.RoomResponse{
			IsWait:   1 < data.MinPlayers,
			RoomInfo: room,
		},
	}

	conn.sendJson(action, response, logp)
}

// <summary>: [Method] JOIN に関する動作を定義します
func actionJoin(req models.WsRequest) {
	action := models.JOIN.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	conn, ok := ConnPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	var response models.WsResponse
	room, exist := RoomPool.Get(req.RoomId)

	// リクエストされた部屋情報がなければエラー
	if !exist {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrRoomNotFound,
		}

		conn.sendJson(action, response, logp)
		return
	}

	// リクエストされた部屋情報とゲームが不一致であればエラー
	if room.GameId != req.GameId {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrMismatchGame,
		}

		conn.sendJson(action, response, logp)
		return
	}

	conn_dup := false
	col_dup := false

	for _, p := range room.Players {
		if p.ConnId == req.ConnId {
			conn_dup = true
		}

		if p.PlayerColor == req.PlayerColor {
			col_dup = true
		}
	}

	// 同じ部屋に入ろうとしていればエラー
	if conn_dup {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrConnectionDuplicated,
		}

		conn.sendJson(action, response, logp)
		return
	}

	// 同じプレイヤー色を使おうとしていればエラー
	if col_dup {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrColorDuplicated,
		}

		conn.sendJson(action, response, logp)
		return
	}

	player := models.PlayerInfoSet{
		ConnId:      req.ConnId,
		PlayerColor: req.PlayerColor,
	}
	room.Players = append(room.Players, player)

	min := models.BgScore[room.GameId].MinPlayers

	logp.Method = models.OK
	response = models.WsResponse{
		Method: models.OK.String(),
		Params: models.RoomResponse{
			IsWait:   len(room.Players) < min,
			RoomInfo: room,
		},
	}

	RoomPool.Set(req.RoomId, room)
	conn.sendJson(action, response, logp)

	for _, p := range room.Players {
		if p.ConnId == req.ConnId {
			continue
		}

		in_conn, ex := ConnPool.Get(p.ConnId)
		if !ex {
			continue
		}

		l := newLogParams(p.ConnId, in_conn.RemoteAddr())
		l.Method = models.NOTIFY
		response.Method = models.NOTIFY.String()

		in_conn.sendJson(action, response, l)
	}
}

// <summary>: [Method] LEAV に関する動作を定義します
func actionLeave(req models.WsRequest) {
	action := models.LEAVE.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	conn, ok := ConnPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	var response models.WsResponse
	notify := ""

	if _, exist := RoomPool.Get(req.RoomId); !exist {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrRoomNotFound,
		}

		conn.sendJson(action, response, logp)
		return
	}

	notify = deletePlayerInfo(req.ConnId)

	logp.Method = models.OK
	response = models.WsResponse{
		Method: models.OK.String(),
		Params: "",
	}

	conn.sendJson(action, response, logp)

	// notifyが空文字のままならブロードキャストはしない
	if notify == "" {
		return
	}

	room, _ := RoomPool.Get(req.RoomId)

	for _, p := range room.Players {
		in_conn, ex := ConnPool.Get(p.ConnId)
		if !ex {
			continue
		}

		l := newLogParams(p.ConnId, in_conn.RemoteAddr())
		l.Method = models.NOTIFY

		min := models.BgScore[room.GameId].MinPlayers
		res := models.WsResponse{
			Method: models.NOTIFY.String(),
			Params: models.RoomResponse{
				IsWait:   len(room.Players) < min,
				RoomInfo: room,
			},
		}

		in_conn.sendJson(action, res, l)
	}
}

// <summary>: [Method] BRDC に関する動作を定義します
func actionBroadcast(req models.WsRequest) {
	action := models.BROADCAST.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	conn, ok := ConnPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	var response models.WsResponse
	room, exist := RoomPool.Get(req.RoomId)

	// リクエストされた部屋情報がなければエラー
	if !exist {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrRoomNotFound,
		}

		conn.sendJson(action, response, logp)
		return
	}

	// リクエストされた部屋情報とゲームが不一致であればエラー
	if room.GameId != req.GameId {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrMismatchGame,
		}

		conn.sendJson(action, response, logp)
		return
	}

	ex_conn := false

	for _, p := range room.Players {
		if p.ConnId == req.ConnId {
			ex_conn = true
			break
		}
	}

	// 部屋にプレイヤーが入室していなければエラー
	if !ex_conn {
		logp.Method = models.ERROR
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrNotInRoom,
		}

		conn.sendJson(action, response, logp)
		return
	}

	point := models.PointResponse{
		Points: req.Points,
		Player: models.PlayerInfoSet{
			ConnId:      req.ConnId,
			PlayerColor: req.PlayerColor,
		},
	}

	logp.Method = models.OK
	response = models.WsResponse{
		Method: models.OK.String(),
		Params: point,
	}

	conn.sendJson(action, response, logp)

	for _, p := range room.Players {
		if p.ConnId == req.ConnId {
			continue
		}

		in_conn, ex := ConnPool.Get(p.ConnId)
		if !ex {
			continue
		}

		l := newLogParams(p.ConnId, in_conn.RemoteAddr())
		l.Method = models.BROADCAST
		response.Method = models.BROADCAST.String()

		in_conn.sendJson(action, response, l)
	}
}

// <summary>: [Method] NONE に関する動作を定義します
func actionNone(req models.WsRequest) {
	action := models.NONE.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	conn, ok := ConnPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	logp.Method = models.ERROR
	response := models.WsResponse{
		Method: models.ERROR.String(),
		Params: models.ErrInvalidMethod,
	}

	conn.sendJson(action, response, logp)
}
