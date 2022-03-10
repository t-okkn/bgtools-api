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

// <summary>: 【エラー】部屋の中に入室していません
var ErrNotInRoom = ErrorMessage{
	Error: "E004",
	Message: "指定された部屋にまだ入室していません",
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

// <summary>: 【エラー】部屋のゲーム内容に誤りがある
var ErrMismatchGame = ErrorMessage{
	Error: "E202",
	Message: "指定された部屋のゲーム内容が違います",
}

// <summary>: 【エラー】部屋に同一の接続元からインしようとしている
var ErrConnectionDuplicated = ErrorMessage{
	Error: "E203",
	Message: "指定された部屋には既に同一接続元から入室しています",
}

// <summary>: 【エラー】部屋に同じ色のプレイヤーが参加しようとしている
var ErrColorDuplicated = ErrorMessage{
	Error: "E204",
	Message: "指定された部屋には既に同色プレイヤーが入室しています",
}
