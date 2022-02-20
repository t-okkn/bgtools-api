package models

// <summary>: 【エラー】部屋が存在しません
var ErrRoomNotFound = FailedResponse{
	Error: "E001",
	Message: "指定された部屋は存在しません",
}

// <summary>: 【エラー】対象の接続が存在しません
var ErrTargetConnectionNotFound = FailedResponse{
	Error: "E002",
	Message: "指定された接続元は存在しません",
}

// <summary>: 【エラー】対象のボードゲーム情報が存在しません
var ErrBoardgameNotFound = FailedResponse{
	Error: "E003",
	Message: "指定されたボードゲームは存在しません",
}

// <summary>: 【エラー】部屋が既に存在します
var ErrRoomExisted = FailedResponse{
	Error: "E101",
	Message: "指定された部屋は既に存在しています",
}

