package ws

import (
	"fmt"

	"bgtools-api/models"
)

// <summary>: [Method] CREA に関する動作を定義します
func actionCreate(req models.WsRequest) {
	logp := newLogParams(req.ConnId, req.ClientIP)
	logp.Prefix = fmt.Sprintf("<%s>", models.CREATE.String())

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log("送信されたconnection_idが不正です")

		return
	}

	// 別室に既に入室していればエラー
	if pc.RoomId != "" {
		pc.sendError(models.ErrEnteredAnotherRoom, logp)
		return
	}

	// リクエストされた部屋情報がなければエラー
	if _, exist := RoomPool.Get(req.RoomId); exist {
		pc.sendError(models.ErrRoomExisted, logp)
		return
	}

	data, exist := models.BgScore[req.GameId]

	// リクエストされたボードゲーム情報がなければエラー
	if !exist {
		pc.sendError(models.ErrBoardgameNotFound, logp)
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

	pc.sendJson(response, logp)
}

// <summary>: [Method] JOIN に関する動作を定義します
func actionJoin(req models.WsRequest) {
	logp := newLogParams(req.ConnId, req.ClientIP)
	logp.Prefix = fmt.Sprintf("<%s>", models.JOIN.String())

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log("送信されたconnection_idが不正です")

		return
	}

	// 別室に既に入室していればエラー
	if pc.RoomId != "" {
		pc.sendError(models.ErrEnteredAnotherRoom, logp)
		return
	}

	room, exist := RoomPool.Get(req.RoomId)

	// リクエストされた部屋情報がなければエラー
	if !exist {
		pc.sendError(models.ErrRoomNotFound, logp)
		return
	}

	// リクエストされた部屋情報とゲームが不一致であればエラー
	if room.GameId != req.GameId {
		pc.sendError(models.ErrMismatchGame, logp)
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
		pc.sendError(models.ErrSameColorExistedInRoom, logp)
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
	pc.sendJson(response, logp)

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
		l.Prefix = fmt.Sprintf("<%s>", models.NOTIFY.String())
		response.Method = models.NOTIFY.String()

		inpc.sendJson(response, l)
	}
}

// <summary>: [Method] LEAV に関する動作を定義します
func actionLeave(req models.WsRequest) {
	logp := newLogParams(req.ConnId, req.ClientIP)
	logp.Prefix = fmt.Sprintf("<%s>", models.LEAVE.String())

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log("送信されたconnection_idが不正です")

		return
	}

	// リクエストされた部屋情報がなければエラー
	if _, exist := RoomPool.Get(req.RoomId); !exist {
		pc.sendError(models.ErrRoomNotFound, logp)
		return
	}

	notify := deletePlayerInfo(req.ConnId)

	logp.Method = models.OK
	response := models.WsResponse{
		Method: models.OK.String(),
		Params: "",
	}

	pc.sendJson(response, logp)

	// notifyが空文字でなければブロードキャスト
	if notify != "" {
		notifyOtherPlayers(notify)
	}
}

// <summary>: [Method] BRDC に関する動作を定義します
func actionBroadcast(req models.WsRequest) {
	logp := newLogParams(req.ConnId, req.ClientIP)
	logp.Prefix = fmt.Sprintf("<%s>", models.BROADCAST.String())

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log("送信されたconnection_idが不正です")

		return
	}

	room, exist := RoomPool.Get(req.RoomId)

	// リクエストされた部屋情報がなければエラー
	if !exist {
		pc.sendError(models.ErrRoomNotFound, logp)
		return
	}

	// リクエストされた部屋情報とゲームが不一致であればエラー
	if room.GameId != req.GameId {
		pc.sendError(models.ErrMismatchGame, logp)
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
		pc.sendError(models.ErrNotInRoom, logp)
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

	pc.sendJson(response, logp)

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
		l.Prefix = logp.Prefix
		response.Method = models.BROADCAST.String()

		inpc.sendJson(response, l)
	}
}

// <summary>: [Method] NONE に関する動作を定義します
func actionNone(req models.WsRequest) {
	logp := newLogParams(req.ConnId, req.ClientIP)
	logp.Prefix = fmt.Sprintf("<%s>", models.NONE.String())

	pc, ok := PlayerPool.Get(req.ConnId)
	if !ok {
		logp.IsProcError = true
		logp.log("送信されたconnection_idが不正です")

		return
	}

	pc.sendError(models.ErrInvalidMethod, logp)
}
