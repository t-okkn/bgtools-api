package models

import "strings"

type Method string

const (
	NONE        Method = "NONE"
	BROADCAST   Method = "BROADCAST"
	CREATE_ROOM Method = "CREATEROOM"
	JOIN_ROOM   Method = "JOINROOM"

	CONNCTED Method = "CONNCTED"
	ACCEPTED Method = "ACCEPTED"
	ERROR    Method = "ERROR"
)

// <summary>: Methodを文字列として表現します
func (m Method) String() string {
	return string(m)
}

// <summary>: 文字列をMethodとして表現します
func ParseMethod(s string) (m Method) {
	switch strings.ToLower(s) {
	case "b", "bc", "broadcast":
		m = BROADCAST

	case "cr", "createroom":
		m = CREATE_ROOM

	case "jr", "joinroom":
		m = JOIN_ROOM

	case "connected":
		m = CONNCTED

	case "accepted":
		m = ACCEPTED

	case "error":
		m = ERROR

	default:
		m = NONE
	}

	return
}
