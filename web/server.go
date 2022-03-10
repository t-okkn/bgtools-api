package web

import (
	"database/sql"
	"fmt"
	"net/http"

	"bgtools-api/db"
	"bgtools-api/models"
	"bgtools-api/ws"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
)

// <summary>: 待ち受けるサーバのルーターを定義します
// <remark>: httpHandlerを受け取る関数にそのまま渡せる
func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")

	//v1.GET("/boardgames", getBoardgames)
	//v1.GET("/boardgames/:gameId", getBoardgames)

	score := v1.Group("score")

	score.GET("/entry", wsEntry)
	score.GET("/rooms/:roomId", checkRoom)
	score.GET("/boardgames", getScoreSupported)
	score.GET("/boardgames/:gameId", getScoreSupported)

	stat := score.Group("statistics")

	stat.GET("/rooms", getRooms)
	stat.GET("/rooms/:roomId", getRooms)
	stat.GET("/connections", getConnections)
	stat.GET("/connections/:connId", getConnections)

	//admin := v1.Group("admin")

	//admin.POST("/boardgames", setBoardgames)
	//admin.PUT("/boardgames/:gameId", updateBoardgames)

	r, err := initDB()
	if err != nil {
		fmt.Printf("initDB: %v\n", err)
	}

	if err := db.LoadBgDataForScore(r); err != nil {
		fmt.Printf("LoadBgData: %v\n", err)
	}

	db.BgRepo = r

	return router
}

// <summary>: WebSocket系の処理が実行されます
func wsEntry(c *gin.Context) {
	ws.EntryPoint(c.Writer, c.Request)
}

// <summary>: 部屋情報が存在しているか確認します
func checkRoom(c *gin.Context) {
	roomid := c.Param("roomId")

	info, exsit := ws.RoomPool.Get(roomid)

	rv := models.CheckRoomResult{
		IsExsit: exsit,
		GameId:  "",
	}

	if exsit {
		rv.GameId = info.GameId
	}

	c.JSON(http.StatusOK, rv)
}

// <summary>: ボードゲーム情報を取得します
func getScoreSupported(c *gin.Context) {
	gameid := c.Param("gameId")

	if gameid == "" {
		c.JSON(http.StatusOK, models.BgScore)

	} else {
		data, ok := models.BgScore[gameid]

		if ok {
			res := make(map[string]models.BgPartialData, 1)
			res[gameid] = data

			c.JSON(http.StatusOK, res)

		} else {
			c.JSON(http.StatusBadRequest, models.ErrBoardgameNotFound)
		}
	}
}

// <summary>: 部屋情報を取得します
func getRooms(c *gin.Context) {
	roomid := c.Param("roomId")
	summary := make([]models.RoomSummary, 0, ws.RoomPool.Count())

	pack := func(id string, room models.RoomInfoSet) {
		gameid := room.GameId

		rs := models.RoomSummary{
			RoomId:   id,
			GameId:   gameid,
			GameData: models.BgScore[gameid],
			Players:  room.Players,
		}

		summary = append(summary, rs)
	}

	if roomid == "" {
		ws.RoomPool.Range(pack)

	} else {
		room, exsit := ws.RoomPool.Get(roomid)

		if exsit {
			pack(roomid, room)

		} else {
			c.JSON(http.StatusBadRequest, models.ErrRoomNotFound)
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}

// <summary>: 接続情報を取得します
func getConnections(c *gin.Context) {
	connid := c.Param("connId")
	summary := make([]models.ConnectionSummary, 0, ws.ConnPool.Count())

	empty := func(id string) models.ConnectionSummary {
		return models.ConnectionSummary{
			ConnId:       id,
			RoomId:       "",
			GameId:       "",
			PlayerColor:  "",
			OtherPlayers: []models.PlayerInfoSet{},
			GameData: models.BgPartialData{
				Title:      "",
				MinPlayers: 0,
				MaxPlayers: 0,
				Colors:     []string{},
			},
		}
	}

	if connid == "" {
		conn_keys := ws.ConnPool.GetKeys()

		inner := func(roomid string, room models.RoomInfoSet) {
			for _, player := range room.Players {
				_, ok := conn_keys[player.ConnId]
				if !ok {
					continue
				}

				gameid := room.GameId
				other := make([]models.PlayerInfoSet, 0, len(room.Players))

				for _, in_player := range room.Players {
					if player.ConnId != in_player.ConnId {
						other = append(other, in_player)
					}
				}

				cs := models.ConnectionSummary{
					ConnId:       player.ConnId,
					RoomId:       roomid,
					GameId:       gameid,
					GameData:     models.BgScore[gameid],
					PlayerColor:  player.PlayerColor,
					OtherPlayers: other,
				}

				delete(conn_keys, player.ConnId)
				summary = append(summary, cs)
			}
		}

		ws.RoomPool.Range(inner)

		if len(conn_keys) != 0 {
			for key := range conn_keys {
				cs := empty(key)
				summary = append(summary, cs)
			}
		}

		c.JSON(http.StatusOK, summary)

	} else {
		if _, exsit := ws.ConnPool.Get(connid); !exsit {
			c.JSON(http.StatusBadRequest, models.ErrConnectionNotFound)
			return
		}

		found := false
		var cs models.ConnectionSummary

		inner := func(roomid string, room models.RoomInfoSet) {
			other := make([]models.PlayerInfoSet, 0, len(room.Players))

			for _, player := range room.Players {
				if connid == player.ConnId {
					found = true
					gameid := room.GameId

					cs = models.ConnectionSummary{
						ConnId:      player.ConnId,
						RoomId:      roomid,
						GameId:      gameid,
						GameData:    models.BgScore[gameid],
						PlayerColor: player.PlayerColor,
					}

				} else {
					other = append(other, player)
				}
			}

			if found {
				cs.OtherPlayers = other
				return
			}
		}

		ws.RoomPool.Range(inner)

		if !found {
			cs = empty(connid)
		}

		res := []models.ConnectionSummary{cs}
		c.JSON(http.StatusOK, res)
	}
}

// <summary>: DBとの接続についての初期処理
func initDB() (*db.BgRepository, error) {
	driver, dsn, err := db.GetDataSourceName()
	if err != nil {
		return nil, err
	}

	var dbmap *gorp.DbMap

	switch driver {
	case "mysql":
		op, err := sql.Open(driver, dsn)
		if err != nil {
			return nil, err
		}

		dial := gorp.MySQLDialect{
			Engine:   "InnoDB",
			Encoding: "utf8mb4",
		}

		dbmap = &gorp.DbMap{
			Db:              op,
			Dialect:         dial,
			ExpandSliceArgs: true,
		}

		models.MapStructsToTables(dbmap)
	}

	return db.NewRepository(dbmap), nil
}
