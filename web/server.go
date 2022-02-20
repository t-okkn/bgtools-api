package web

import (
	"net/http"

	"bgtools-api/models"
	"bgtools-api/ws"

	"github.com/gin-gonic/gin"
)

// <summary>: 待ち受けるサーバのルーターを定義します
// <remark>: httpHandlerを受け取る関数にそのまま渡せる
func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")

	v1.GET("/join", wsEntry)
	v1.GET("/check/rooms/:roomId", checkRoom)
	v1.GET("/boardgames", getBoardgames)
	v1.GET("/boardgames/:gameId", getBoardgames)

	stat := v1.Group("statistics")

	stat.GET("/rooms", getRooms)
	stat.GET("/rooms/:roomId", getRooms)
	stat.GET("/connections", getConnections)
	stat.GET("/connections/:connId", getConnections)

	return router
}

// <summary>: WebSocket系の処理が実行されます
func wsEntry(c *gin.Context) {
	ws.EntryPoint(c.Writer, c.Request)
}

// <summary>: 部屋情報が存在しているか確認します
func checkRoom(c *gin.Context) {
	roomid := c.Param("roomId")

	info, exsit := ws.RoomPool[roomid]

	rv := models.CheckRoomResult{
		IsExsit: exsit,
		GameId:  "",
	}

	if exsit {
		rv.GameId = info.GameId
	}

	c.JSON(http.StatusOK, rv)
}

func getBoardgames(c *gin.Context) {
	gameid := c.Param("gameId")

	if gameid == "" {
		c.JSON(http.StatusOK, models.BGCollection)

	} else {
		data, ok := models.BGCollection[gameid]

		if ok {
			res := make(map[string]models.BoardgameData, 1)
			res[gameid] = data

			c.JSON(http.StatusOK, models.BGCollection)

		} else {
			c.JSON(http.StatusBadRequest, models.ErrBoardgameNotFound)
		}
	}
}

func getRooms(c *gin.Context) {
	roomid := c.Param("roomId")
	summary := make([]models.RoomSummary, 0, len(ws.RoomPool))

	pack := func(id string, room models.RoomInfoSet) {
		gameid := room.GameId

		rs := models.RoomSummary{
			RoomId:   id,
			GameId:   gameid,
			GameData: models.BGCollection[gameid],
			Players:  room.Players,
		}

		summary = append(summary, rs)
	}

	if roomid == "" {
		for id, room := range ws.RoomPool {
			pack(id, room)
		}

	} else {
		room, exsit := ws.RoomPool[roomid]

		if exsit {
			pack(roomid, room)

		} else {
			c.JSON(http.StatusBadRequest, models.ErrRoomNotFound)
			return
		}

		c.JSON(http.StatusOK, summary)
	}
}

func getConnections(c *gin.Context) {
	connid := c.Param("connId")
	summary := make([]models.ConnectionSummary, 0, len(ws.WsConnPool))

	if connid == "" {
		conns := make(map[string]struct{}, len(ws.WsConnPool))

		for key := range ws.WsConnPool {
			conns[key] = struct{}{}
		}

		for roomid, room := range ws.RoomPool {
			for cid, color := range room.Players {
				_, ok := conns[cid]
				if !ok {
					continue
				}

				gameid := room.GameId
				other := make(map[string]string, len(room.Players))

				for in_cid, in_color := range room.Players {
					if cid != in_cid {
						other[in_cid] = in_color
					}
				}

				cs := models.ConnectionSummary{
					ConnId:       connid,
					RoomId:       roomid,
					GameId:       gameid,
					GameData:     models.BGCollection[gameid],
					PlayerColor:  color,
					OtherPlayers: other,
				}

				delete(conns, cid)
				summary = append(summary, cs)
			}
		}

		if len(conns) == 0 {
			c.JSON(http.StatusOK, summary)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

	} else {
		found := false
		var cs models.ConnectionSummary

		for roomid, room := range ws.RoomPool {
			other := make(map[string]string, len(room.Players))

			for cid, color := range room.Players {
				if connid == cid {
					found = true
					gameid := room.GameId

					cs = models.ConnectionSummary{
						ConnId:      connid,
						RoomId:      roomid,
						GameId:      gameid,
						GameData:    models.BGCollection[gameid],
						PlayerColor: color,
					}

				} else {
					other[cid] = color
				}
			}

			if found {
				cs.OtherPlayers = other
				break
			}
		}

		if found {
			summary = append(summary, cs)
			c.JSON(http.StatusOK, summary)

		} else {
			c.JSON(http.StatusBadRequest, models.ErrTargetConnectionNotFound)
		}
	}
}
