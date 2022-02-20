package models

import "strings"

type Method string

const (
	NONE        Method = "none"
	BROADCAST   Method = "broadcast"
	CREATE_ROOM Method = "create_room"
	JOIN_ROOM   Method = "join_room"

	CONNCTED Method = "connected"
	ACCEPTED Method = "accepted"
	FAILED   Method = "failed"
)

func (m Method) String() string {
	return string(m)
}

func ParseMethod(s string) (m Method) {
	switch strings.ToLower(s) {
	case "b", "bc", "broadcast":
		m = BROADCAST

	case "cr", "create_room":
		m = CREATE_ROOM

	case "jr", "join_room":
		m = JOIN_ROOM

	default:
		m = NONE
	}

	return
}
