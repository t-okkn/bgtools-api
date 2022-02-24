package models

import "strings"

type Method string

const (
	NONE        Method = "NONE"
	BROADCAST   Method = "BRDC"
	CREATE_ROOM Method = "CRRM"
	JOIN_ROOM   Method = "JNRM"

	CONNCTED Method = "CONN"
	OK       Method = "OK"
	ERROR    Method = "ERR"
)

// <summary>: Methodを文字列として表現します
func (m Method) String() string {
	return string(m)
}

// <summary>: 文字列をMethodとして表現します
func ParseMethod(s string) (m Method) {
	switch strings.ToLower(s) {
	case "brdc":
		m = BROADCAST

	case "crrm":
		m = CREATE_ROOM

	case "jnrm":
		m = JOIN_ROOM

	case "conn":
		m = CONNCTED

	case "ok":
		m = OK

	case "err":
		m = ERROR

	default:
		m = NONE
	}

	return
}
