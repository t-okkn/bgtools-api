package models

// WebSocketsからの返却用データの構造体
type WsResponse struct {
	Method  string `json:"method"`
	Message string `json:"message"`
}