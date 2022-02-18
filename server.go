package main

import (
	"bgtools-api/ws"

	"github.com/gin-gonic/gin"
)

// <summary>: 待ち受けるサーバのルーターを定義します
// <remark>: httpHandlerを受け取る関数にそのまま渡せる
func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")

	v1.GET("/ws", wsServe)

	return router
}

// <summary>: WebSocket系の処理が実行されます
func wsServe(c *gin.Context) {
	ws.Handler(c.Writer, c.Request)
}
