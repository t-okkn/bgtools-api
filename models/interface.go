package models

// <summary>: プレーヤーの情報
type PlayerInfoSet struct {
	ConnId      string `json:"connection_id"`
	PlayerColor string `json:"player_color"`
}

// <summary>: 部屋のゲーム内容と部屋にいるプレーヤー情報
type RoomInfoSet struct {
	GameId  string          `json:"game_id"`
	Players []PlayerInfoSet `json:"players"`
}

// <summary>: WebSocketでの受信用データの構造体
type WsRequest struct {
	Method      string   `json:"method"`
	ConnId      string   `json:"connection_id"`
	RoomId      string   `json:"room_id"`
	GameId      string   `json:"game_id"`
	PlayerColor string   `json:"player_color"`
	Points      []int    `json:"points"`
}

// <summary>: WebSocketからの返却用データの構造体
type WsResponse struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// <summary>: 接続時、Response内のParamsに使用される構造体
type ConnectResponse struct {
	ConnId string `json:"connection_id"`
}

// <summary>: 部屋の情報伝達時、Response内のParamsに使用される構造体
type RoomResponse struct {
	IsWait   bool        `json:"is_wait"`
	RoomId   string      `json:"room_id"`
	RoomInfo RoomInfoSet `json:"room"`
}

// <summary>: 得点のブロードキャスト時、Response内のParamsに使用される構造体
type PointResponse struct {
	Player PlayerInfoSet `json:"player"`
	Points []int         `json:"points"`
}

// <summary>: MethodがOKの時、特に伝達する情報がない場合に使用される構造体
type OKMessage struct {
	Message string `json:"message"`
}

// <summary>: エラーに関する情報を格納する構造体
type ErrorMessage struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// <summary>: 部屋の存在確認に使用される構造体
type CheckRoomResult struct {
	IsExist bool   `json:"is_exist"`
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
	ConnId       string          `json:"connection_id"`
	RoomId       string          `json:"room_id"`
	GameId       string          `json:"game_id"`
	GameData     BgPartialData   `json:"game_data"`
	PlayerColor  string          `json:"player_color"`
	OtherPlayers []PlayerInfoSet `json:"other_players"`
}

// <summary>: 部屋情報を一覧表示するための構造体
type RoomSummary struct {
	RoomId   string          `json:"room_id"`
	GameId   string          `json:"game_id"`
	GameData BgPartialData   `json:"game_data"`
	Players  []PlayerInfoSet `json:"players"`
}
