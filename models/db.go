package models

import "github.com/go-gorp/gorp"

// <summary>: 対応しているボードゲームの情報
var BgScore = map[string]BgPartialData{}

type MstrBoardgame struct {
	Id              string `db:"id, primarykey" json:"id"`
	UniqueName      string `db:"unique_name" json:"unique_name"`
	Title           string `db:"title" json:"title"`
	MinPlayers      int    `db:"min_players" json:"min_players"`
	MaxPlayers      int    `db:"max_players" json:"max_players"`
	PlayingTime     string `db:"playing_time" json:"playing_time"`
	MinAge          int    `db:"min_age" json:"min_age"`
	IsExpansion     bool   `db:"is_expansion" json:"is_expansion"`
	ExpansionBaseId string `db:"expansion_base_id" json:"expansion_base_id"`
	ProductUrl      string `db:"product_url" json:"product_url"`
	BodogeHoobbyNet bool   `db:"bodoge_hoobby_net" json:"bodoge_hoobby_net"`
	ScoreTool       bool   `db:"score_tool" json:"score_tool"`
}

type MstrUser struct {
	Id          string `db:"id, primarykey" json:"id"`
	UserName    string `db:"user_name" json:"user_name"`
	MailAddress string `db:"mail_address" json:"mail_address"`
	AuthKey     string `db:"auth_key" json:"auth_key"`
}

type MstrColor struct {
	GameId string `db:"game_id" json:"game_id"`
	Color  string `db:"color" json:"color"`
}

type TranOwn struct {
	UserId string `db:"user_id" json:"user_id"`
	GameId string `db:"game_id" json:"game_id"`
}

type BgScoreSupport struct {
	GameId     string `db:"game_id"`
	Title      string `db:"title"`
	MinPlayers int    `db:"min_players"`
	MaxPlayers int    `db:"max_players"`
	Color      string `db:"color"`
}

// MapStructsToTables 構造体と物理テーブルの紐付け
func MapStructsToTables(dbmap *gorp.DbMap) {
	dbmap.AddTableWithName(MstrBoardgame{}, "M_BOARDGAME").SetKeys(false, "Id")
	dbmap.AddTableWithName(MstrUser{}, "M_USER").SetKeys(false, "Id")
	dbmap.AddTableWithName(MstrColor{}, "M_COLOR")
	dbmap.AddTableWithName(TranOwn{}, "T_OWN")
}
