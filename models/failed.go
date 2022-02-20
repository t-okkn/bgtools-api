package models

var ErrRoomNotFound = FailedResponse{
	Error: "E001",
	Message: "指定された部屋は存在しません",
}

var ErrTargetConnectionNotFound = FailedResponse{
	Error: "E002",
	Message: "指定された接続元は存在しません",
}

var ErrBoardgameNotFound = FailedResponse{
	Error: "E003",
	Message: "指定されたボードゲームは存在しません",
}

var ErrRoomExisted = FailedResponse{
	Error: "E101",
	Message: "指定された部屋は既に存在しています",
}

