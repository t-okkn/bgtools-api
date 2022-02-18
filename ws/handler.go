package ws

import (
	"fmt"
	"net/http"

	"bgtools-api/models"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %s", err)
		return
	}

	res := models.WsResponse{
		Method:  models.NONE.String(),
		Message: "Hello WebSocket!!",
	}

	if err = conn.WriteJSON(res); err != nil {
		fmt.Printf("Failed to ssend message: %s", err)
	}
}
