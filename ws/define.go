package ws

import (
	"net/http"

	"bgtools-api/models"

	"github.com/gorilla/websocket"
)

type WsConnection struct {
	*websocket.Conn
}

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	chWsReq = make(chan models.WsRequest)

	WsConnPool = map[string]*WsConnection{}

	RoomPool = map[string]models.RoomInfoSet{}
)
