package main

import (
	"flag"
	"fmt"

	"bgtools-api/models"
	"bgtools-api/web"
	"bgtools-api/ws"
)

const LISTEN_PORT string = ":8506"

var (
	Version string
	Revision string
)

// <summary>: main関数（サーバを開始します）
func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		fmt.Println(Version, Revision)
		return
	}

	if err := models.LoadBoardgameData(); err != nil {
		panic("Fail to load boardgame data")
	}

	go ws.ListenAndServe()

	web.SetupRouter().Run(LISTEN_PORT)
}
