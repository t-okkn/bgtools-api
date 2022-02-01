package main

import (
	"github.com/gin-gonic/gin"
)

// summary => 待ち受けるサーバのルーターを定義します
// return::*gin.Engine =>
// remark => httpHandlerを受け取る関数にそのまま渡せる
/////////////////////////////////////////
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", sample)

	 return router
}

// summary => サンプル
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func sample(c *gin.Context) {
	//途中で処理を中断
	//c.Abort()
}

