package ws

import (
	"fmt"
	"time"

	"bgtools-api/models"
)

// <summary>: [Method] CREATEROOM に関する動作を定義します
func actionCreateRoom(req models.WsRequest) {
	logp := newLogParams()
	start := time.Now()

	logp.ConnId = req.PlayerInfo.ConnId
	logp.ClientIP = req.ClientIP

	conn, ok := WsConnPool[req.PlayerInfo.ConnId]
	if !ok {
		logp.IsProcError = true
		logp.Message = "<CREATEROOM> 送信されたconnection_idが不正です"
		logp.log()

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
			Method: models.ACCEPTED.String(),
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
		logp.Message =
			fmt.Sprintf("<CREATEROOM> 処理時間：%v, 送信完了：%+v", d, response)
		logp.log()

	} else {
		logp.IsProcError = true
		logp.Message =
			fmt.Sprintf("<CREATEROOM> メッセージの送信に失敗しました：%s", err)
		logp.log()
	}
}

// <summary>: [Method] NONE に関する動作を定義します
func actionNone(req models.WsRequest) {
	logp := newLogParams()
	start := time.Now()

	logp.ConnId = req.PlayerInfo.ConnId
	logp.ClientIP = req.ClientIP

	conn, ok := WsConnPool[req.PlayerInfo.ConnId]
	if !ok {
		logp.IsProcError = true
		logp.Message = "<NONE> 送信されたconnection_idが不正です"
		logp.log()

		return
	}

	response := models.WsResponse{
		Method: models.ERROR.String(),
		Params: models.ErrInvalidMethod,
	}

	d := getProcTime(start)

	if err := conn.WriteJSON(response); err == nil {
		logp.Method = models.ParseMethod(response.Method)
		logp.Message =
			fmt.Sprintf("<NONE> 処理時間：%v, 送信完了：%+v", d, response)
		logp.log()

	} else {
		logp.IsProcError = true
		logp.Message =
			fmt.Sprintf("<NONE> メッセージの送信に失敗しました：%s", err)
		logp.log()
	}
}
