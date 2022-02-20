package ws

import (
	"fmt"

	"bgtools-api/models"
)

// <summary>: [Method] create_room に関する動作を定義します
func actionCreateRoom(req models.WsRequest) {
	conn, ok := WsConnPool[req.PlayerInfo.ConnId]
	if !ok {
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
			Method: models.FAILED.String(),
			Params: models.ErrRoomExisted,
		}
	}

	if err := conn.WriteJSON(response); err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
	}
}
