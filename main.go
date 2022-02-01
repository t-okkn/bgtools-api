package main

import (
	"flag"
	"fmt"
	// "nginx/unit"
)

const LISTEN_PORT string = ":8506"

var (
	Version string
	Revision string
)

// summary => main関数（サーバを開始します）
/////////////////////////////////////////
func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		fmt.Println(Version, Revision)
		return
	}

	// unit.ListenAndServe(LISTEN_PORT, SetupRouter())
	SetupRouter().Run(LISTEN_PORT)
}
