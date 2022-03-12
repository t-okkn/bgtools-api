package models

import "strings"

type Method string

const (
	BROADCAST Method = "BROADCAST"
	CREATE    Method = "CREATE"
	JOIN      Method = "JOIN"
	LEAVE     Method = "LEAVE"

	NONE       Method = "NONE"
	CONNECT    Method = "CONNECT"
	DISCONNECT Method = "DISCONNECT"
	EJECT      Method = "EJECT"
	NOTIFY     Method = "NOTIFY"
	OK         Method = "OK"
	ERROR      Method = "ERROR"
)

// <summary>: Methodを文字列として表現します
func (m Method) String() string {
	return string(m)
}

// <summary>: 文字列をMethodとして表現します
func ParseMethod(s string) (m Method) {
	switch strings.ToUpper(s) {
	case "BROADCAST":
		m = BROADCAST

	case "CREATE":
		m = CREATE

	case "JOIN":
		m = JOIN

	case "LEAVE":
		m = LEAVE

	case "CONNECT":
		m = CONNECT

	case "DISCONNECT":
		m = DISCONNECT

	case "EJECT":
		m = EJECT

	case "NOTIFY":
		m = NOTIFY

	case "OK":
		m = OK

	case "ERROR":
		m = ERROR

	default:
		m = NONE
	}

	return
}
