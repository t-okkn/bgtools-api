package models

import "strings"

type Method string

const (
	BROADCAST  Method = "BRDC"
	CREATE     Method = "CRET"
	JOIN       Method = "JOIN"
	LEAVE      Method = "LEAV"

	NONE    Method = "NONE"
	CONNECT Method = "CONN"
	EJECT   Method = "EJCT"
	NOTIFY  Method = "NTFY"
	OK      Method = "OK"
	ERROR   Method = "ERR"
)

// <summary>: Methodを文字列として表現します
func (m Method) String() string {
	return string(m)
}

// <summary>: 文字列をMethodとして表現します
func ParseMethod(s string) (m Method) {
	switch strings.ToUpper(s) {
	case "BRDC":
		m = BROADCAST

	case "CRET":
		m = CREATE

	case "JOIN":
		m = JOIN

	case "LEAV":
		m = LEAVE

	case "CONN":
		m = CONNECT

	case "EJCT":
		m = EJECT

	case "NTFY":
		m = NOTIFY

	case "OK":
		m = OK

	case "ERR":
		m = ERROR

	default:
		m = NONE
	}

	return
}
