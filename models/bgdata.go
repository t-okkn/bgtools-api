package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// <summary>: ボードゲームのデータを格納します
type BoardgameData struct {
	Title      string `json:"title"`
	MinPlayers int    `json:"min_players"`
	MaxPlayers int    `json:"max_players"`
}

// <summary>: 対応しているボードゲームの情報
var BGCollection = map[string]BoardgameData{}

// <summary>: ボードゲームのデータを読み込みます
func LoadBoardgameData() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	file := filepath.Join(filepath.Dir(exe), "bgdata.json")

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, &BGCollection)
}
