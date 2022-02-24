package ws

import (
	"fmt"
	"time"

	"bgtools-api/models"
)

// <summary>: [Method] CRRM に関する動作を定義します
func actionCreateRoom(req models.WsRequest) {
	logp := newLogParams(req.PlayerInfo.ConnId, req.ClientIP)
	start := time.Now()

	conn, ok := WsConnPool[req.PlayerInfo.ConnId]
	if !ok {
		logp.IsProcError = true
		logp.log("<CRRM> 送信されたconnection_idが不正です")

		return
	}

	var response models.WsResponse
	_, exsit := RoomPool[req.PlayerInfo.RoomId]

	if !exsit {
		data := models.BGCollection[req.PlayerInfo.GameId]

		players := make(map[string]string, data.MaxPlayers)
		players[req.PlayerInfo.ConnId] = req.PlayerInfo.PlayerColor

		room := models.RoomInfoSet{
			GameId:  req.PlayerInfo.GameId,
			Players: players,
		}
		RoomPool[req.PlayerInfo.RoomId] = room

		response = models.WsResponse{
			Method: models.OK.String(),
			Params: models.RoomResponse{
				IsWait:   data.MinPlayers >= 2,
				RoomInfo: room,
			},
		}

	} else {
		response = models.WsResponse{
			Method: models.ERROR.String(),
			Params: models.ErrRoomExisted,
		}
	}

	d := getProcTime(start)

	if err := conn.WriteJSON(response); err == nil {
		logp.Method = models.ParseMethod(response.Method)
		logp.log(fmt.Sprintf("<CRRM> 処理時間：%v, 送信完了：%+v", d, response))

	} else {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<CRRM> メッセージの送信に失敗しました：%s", err))
	}
}

// <summary>: [Method] NONE に関する動作を定義します
func actionNone(req models.WsRequest) {
	logp := newLogParams(req.PlayerInfo.ConnId, req.ClientIP)
	start := time.Now()

	conn, ok := WsConnPool[req.PlayerInfo.ConnId]
	if !ok {
		logp.IsProcError = true
		logp.log("<NONE> 送信されたconnection_idが不正です")

		return
	}

	response := models.WsResponse{
		Method: models.ERROR.String(),
		Params: models.ErrInvalidMethod,
	}

	d := getProcTime(start)

	if err := conn.WriteJSON(response); err == nil {
		logp.Method = models.ParseMethod(response.Method)
		logp.log(fmt.Sprintf("<NONE> 処理時間：%v, 送信完了：%+v", d, response))

	} else {
		logp.IsProcError = true
		logp.log(fmt.Sprintf("<NONE> メッセージの送信に失敗しました：%s", err))
	}
}
