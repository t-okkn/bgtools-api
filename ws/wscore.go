package ws

import (
	"fmt"
	"net/http"

	"bgtools-api/models"

	"github.com/google/uuid"
)

// <summary>: Websocket接続時に行われる動作
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %s\n", err)
		return
	}

	wsconn := &WsConnection{conn}
	obj, _ := uuid.NewRandom()
	uuidStr := obj.String()

	WsConnPool[uuidStr] = wsconn

	res := models.WsResponse{
		Method: models.CONNCTED.String(),
		Params: models.ConnectedResponse{
			ConnId: uuidStr,
		},
	}

	if err := conn.WriteJSON(res); err != nil {
		fmt.Printf("Failed to send message: %s\n", err)
	}

	go readRequests(wsconn)
}

// <summary>: Websocketでの電文のやり取りを行います
func ListenAndServe() {
	for {
		// メッセージが入るまで、ここでブロック
		e := <-chWsReq
		var action func(models.WsRequest)

		switch models.ParseMethod(e.Method) {
		case models.CREATE_ROOM:
			action = actionCreateRoom
		
		default:
			fmt.Println("Method is NONE")
			continue
		}

		action(e)
	}
}

// <summary>: 受信した内容を読み取ります
func readRequests(conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %v\n", r)
		}
	}()

	var req models.WsRequest

	for {
		err := conn.ReadJSON(&req)
		if err == nil {
			chWsReq <- req
		}
	}
}
