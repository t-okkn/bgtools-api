package models

import "net"

// <summary>: 接続元プレーヤーの固有情報
type PlayerInfoSet struct {
	RoomId      string `json:"room_id"`
	GameId      string `json:"game_id"`
	PlayerColor string `json:"player_color"`
	ConnId      string `json:"connection_id"`
}

// <summary>: 部屋のゲーム内容と部屋にいるプレーヤー情報
type RoomInfoSet struct {
	GameId  string            `json:"game_id"`
	Players map[string]string `json:"players"`
}

// <summary>: WebSocketsでの受信用データの構造体
type WsRequest struct {
	Method     string        `json:"method"`
	PlayerInfo PlayerInfoSet `json:"player"`
	Points     []int         `json:"points"`
	ClientIP   net.Addr      `json:"-"`
}

// <summary>: WebSocketsからの返却用データの構造体
type WsResponse struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// <summary>: 接続時、Response内のParamsに使用される構造体
type ConnectedResponse struct {
	ConnId string `json:"connection_id"`
}

// <summary>: 部屋へ接続をしに来た時、Response内のParamsに使用される構造体
type RoomResponse struct {
	IsWait   bool        `json:"is_wait"`
	RoomInfo RoomInfoSet `json:"room"`
}

// <summary>: エラーに関する情報を格納する構造体
type ErrorMessage struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// <summary>: 部屋の存在確認に使用される構造体
type CheckRoomResult struct {
	IsExsit bool   `json:"is_exsit"`
	GameId  string `json:"game_id"`
}

// <summary>: スコアツール対応のボードゲームデータを格納します
type BgPartialData struct {
	Title      string   `json:"title"`
	MinPlayers int      `json:"min_players"`
	MaxPlayers int      `json:"max_players"`
	Colors     []string `json:"colors"`
}

// <summary>: 接続情報を一覧表示するための構造体
type ConnectionSummary struct {
	ConnId       string            `json:"connection_id"`
	RoomId       string            `json:"room_id"`
	GameId       string            `json:"game_id"`
	GameData     BgPartialData     `json:"gama_data"`
	PlayerColor  string            `json:"player_color"`
	OtherPlayers map[string]string `json:"other_players"`
}

// <summary>: 部屋情報を一覧表示するための構造体
type RoomSummary struct {
	RoomId   string            `json:"room_id"`
	GameId   string            `json:"game_id"`
	GameData BgPartialData     `json:"gama_data"`
	Players  map[string]string `json:"players"`
}
