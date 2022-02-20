package models

// 接続元プレーヤーの固有情報
type PlayerInfoSet struct {
	RoomId      string `json:"room_id"`
	GameId      string `json:"game_id"`
	PlayerColor string `json:"player_color"`
	ConnId      string `json:"connection_id"`
}

// 部屋のゲーム内容と部屋にいるプレーヤー情報
type RoomInfoSet struct {
	GameId  string
	Players map[string]string
}

// WebSocketsでの受信用データの構造体
type WsRequest struct {
	Method     string        `json:"method"`
	PlayerInfo PlayerInfoSet `json:"player_info"`
	Points     []int         `json:"points"`
}

// WebSocketsからの返却用データの構造体
type WsResponse struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// 接続時、Response内のParamsに使用される構造体
type ConnectedResponse struct {
	ConnId string `json:"connection_id"`
}

// 部屋へ接続をしに来た時、Response内のParamsに使用される構造体
type RoomResponse struct {
	IsWait   bool        `json:"is_wait"`
	RoomInfo RoomInfoSet `json:"room_info"`
}

type FailedResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type CheckRoomResult struct {
	IsExsit bool   `json:"is_exsit"`
	GameId  string `json:"game_id"`
}

type ConnectionSummary struct {
	ConnId       string            `json:"connection_id"`
	RoomId       string            `json:"room_id"`
	GameId       string            `json:"game_id"`
	GameData     BoardgameData     `json:"gama_data"`
	PlayerColor  string            `json:"player_color"`
	OtherPlayers map[string]string `json:"other_players"`
}

type RoomSummary struct {
	RoomId   string            `json:"room_id"`
	GameId   string            `json:"game_id"`
	GameData BoardgameData     `json:"gama_data"`
	Players  map[string]string `json:"players"`
}
