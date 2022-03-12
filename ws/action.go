package ws

import (
	"fmt"

	"bgtools-api/models"
)

// <summary>: [Method] CREA に関する動作を定義します
func actionCreate(req models.WsRequest) {
	action := models.CREATE.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	// 別室に既に入室していればエラー
	if pc.RoomId != "" {
		pc.sendError(action, models.ErrEnteredAnotherRoom, logp)
		return
	}

	// リクエストされた部屋情報がなければエラー
	if _, exist := RoomPool.Get(req.RoomId); exist {
		pc.sendError(action, models.ErrRoomExisted, logp)
		return
	}

	data, exist := models.BgScore[req.GameId]

	// リクエストされたボードゲーム情報がなければエラー
	if !exist {
		pc.sendError(action, models.ErrBoardgameNotFound, logp)
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
	PlayerPool.SetRoomId(req.ConnId, req.RoomId)

	logp.Method = models.OK
	response := models.WsResponse{
		Method: models.OK.String(),
		Params: models.RoomResponse{
			IsWait:   1 < data.MinPlayers,
			RoomInfo: room,
		},
	}

	pc.sendJson(action, response, logp)
}

// <summary>: [Method] JOIN に関する動作を定義します
func actionJoin(req models.WsRequest) {
	action := models.JOIN.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	// 別室に既に入室していればエラー
	if pc.RoomId != "" {
		pc.sendError(action, models.ErrEnteredAnotherRoom, logp)
		return
	}

	room, exist := RoomPool.Get(req.RoomId)

	// リクエストされた部屋情報がなければエラー
	if !exist {
		pc.sendError(action, models.ErrRoomNotFound, logp)
		return
	}

	// リクエストされた部屋情報とゲームが不一致であればエラー
	if room.GameId != req.GameId {
		pc.sendError(action, models.ErrMismatchGame, logp)
		return
	}

	col_ex := false

	for _, p := range room.Players {
		if p.PlayerColor == req.PlayerColor {
			col_ex = true
			break
		}
	}

	// 同じプレイヤー色を使おうとしていればエラー
	if col_ex {
		pc.sendError(action, models.ErrSameColorExistedInRoom, logp)
		return
	}

	player := models.PlayerInfoSet{
		ConnId:      req.ConnId,
		PlayerColor: req.PlayerColor,
	}
	room.Players = append(room.Players, player)

	min := models.BgScore[room.GameId].MinPlayers

	logp.Method = models.OK
	response := models.WsResponse{
		Method: models.OK.String(),
		Params: models.RoomResponse{
			IsWait:   len(room.Players) < min,
			RoomInfo: room,
		},
	}

	RoomPool.Set(req.RoomId, room)
	PlayerPool.SetRoomId(req.ConnId, req.RoomId)
	pc.sendJson(action, response, logp)

	for _, p := range room.Players {
		if p.ConnId == req.ConnId {
			continue
		}

		inpc, ex := PlayerPool.Get(p.ConnId)
		if !ex {
			continue
		}

		l := newLogParams(p.ConnId, inpc.C.RemoteAddr())
		l.Method = models.NOTIFY
		response.Method = models.NOTIFY.String()

		inpc.sendJson(action, response, l)
	}
}

// <summary>: [Method] LEAV に関する動作を定義します
func actionLeave(req models.WsRequest) {
	action := models.LEAVE.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	// リクエストされた部屋情報がなければエラー
	if _, exist := RoomPool.Get(req.RoomId); !exist {
		pc.sendError(action, models.ErrRoomNotFound, logp)
		return
	}

	notify := deletePlayerInfo(req.ConnId)

	logp.Method = models.OK
	response := models.WsResponse{
		Method: models.OK.String(),
		Params: "",
	}

	pc.sendJson(action, response, logp)

	// notifyが空文字でなければブロードキャスト
	if notify != "" {
		notifyOtherPlayers(notify)
	}
}

// <summary>: [Method] BRDC に関する動作を定義します
func actionBroadcast(req models.WsRequest) {
	action := models.BROADCAST.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	room, exist := RoomPool.Get(req.RoomId)

	// リクエストされた部屋情報がなければエラー
	if !exist {
		pc.sendError(action, models.ErrRoomNotFound, logp)
		return
	}

	// リクエストされた部屋情報とゲームが不一致であればエラー
	if room.GameId != req.GameId {
		pc.sendError(action, models.ErrMismatchGame, logp)
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
		pc.sendError(action, models.ErrNotInRoom, logp)
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
	response := models.WsResponse{
		Method: models.OK.String(),
		Params: point,
	}

	pc.sendJson(action, response, logp)

	for _, p := range room.Players {
		if p.ConnId == req.ConnId {
			continue
		}

		inpc, ex := PlayerPool.Get(p.ConnId)
		if !ex {
			continue
		}

		l := newLogParams(p.ConnId, inpc.C.RemoteAddr())
		l.Method = models.BROADCAST
		response.Method = models.BROADCAST.String()

		inpc.sendJson(action, response, l)
	}
}

// <summary>: [Method] NONE に関する動作を定義します
func actionNone(req models.WsRequest) {
	action := models.NONE.String()
	logp := newLogParams(req.ConnId, req.ClientIP)

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<%s> 送信されたconnection_idが不正です", action))

		return
	}

	pc.sendError(action, models.ErrInvalidMethod, logp)
}
