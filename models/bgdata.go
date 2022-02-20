package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type BoardgameData struct {
	Title      string `json:"title"`
	MinPlayers int    `json:"min_players"`
	MaxPlayers int    `json:"max_players"`
}

var BGCollection = map[string]BoardgameData{}

func RepackBoardgameData() error {
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
