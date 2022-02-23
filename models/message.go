package models

// <summary>: 【エラー】部屋が存在しません
var ErrRoomNotFound = ErrorMessage{
	Error: "E001",
	Message: "指定された部屋は存在しません",
}

// <summary>: 【エラー】対象の接続が存在しません
var ErrTargetConnectionNotFound = ErrorMessage{
	Error: "E002",
	Message: "指定された接続元は存在しません",
}

// <summary>: 【エラー】対象のボードゲーム情報が存在しません
var ErrBoardgameNotFound = ErrorMessage{
	Error: "E003",
	Message: "指定されたボードゲームは存在しません",
}

// <summary>: 【エラー】無効なメソッド
var ErrInvalidMethod = ErrorMessage{
	Error: "E101",
	Message: "無効なメソッドが指定されました",
}

// <summary>: 【エラー】部屋が既に存在します
var ErrRoomExisted = ErrorMessage{
	Error: "E201",
	Message: "指定された部屋は既に存在しています",
}
