package models

import "strings"

type Method string

const (
	NONE Method = "none"
	BC Method = "broadcast"
)

func (m Method) String() string {
	return string(m)
}

func ParseMethod(s string) (m Method) {
	switch strings.ToLower(s) {
	case "b", "bc", "broadcast":
		m = BC

	default:
		m = NONE
	}

	return
}