package ws

import (
	"fmt"
	"io"
	"os"
	"time"

	"bgtools-api/models"
)

const (
	blue    string = "\033[97;44m"
	cyan    string = "\033[97;46m"
	red     string = "\033[97;41m"
	magenta string = "\033[97;45m"
	green   string = "\033[97;42m"
	bgreen  string = "\033[30;42m"
	yellow  string = "\033[30;103m"
	white   string = "\033[30;47m"
	reset   string = "\033[0m"
)

var (
	enableColorMode   bool      = true
	outputDestination io.Writer = os.Stdout
)

// <summary>: WebSocket用ログのパラメータを示す構造体
type logParams struct {
	ClientIP    string
	ConnId      string
	Prefix      string
	Method      models.Method
	IsProcError bool
}

// <summary>: ログの色つけを有効化します
func EnableColorLog() {
	enableColorMode = true
}

// <summary>: ログの色つけを無効化します
func DisableColorLog() {
	enableColorMode = false
}

// <summary>: ログの出力先を変更します
// <remark>: defaultは標準出力
func ChangeOutputDestination(dest io.Writer) {
	outputDestination = dest
}

// <summary>: 新規logParams構造体を生成します
func newLogParams(connid string) logParams {
	ip, _ := getIpPort(connid)

	return logParams{
		ClientIP:    ip,
		ConnId:      connid,
		Method:      models.NONE,
		IsProcError: false,
	}
}

// <summary>: logParamsの情報からログを書き込みます
func (p logParams) log(message string) {
	resetColor := reset
	if !enableColorMode {
		resetColor = ""
	}

	tag := "WS"
	var str string

	if p.IsProcError {
		tag = "WS-Error"
		redColor := red

		if !enableColorMode {
			redColor = ""
		}

		str = fmt.Sprintf("[%s] %v | %15s | %s | %s[ERROR]%s %s",
			tag,
			time.Now().Format("2006/01/02 - 15:04:05"),
			p.ClientIP,
			p.ConnId,
			redColor, resetColor, message,
		)

	} else {
		prefix := p.Prefix + " "
		if p.Prefix == "" {
			prefix = ""
		}

		str = fmt.Sprintf("[%s] %v | %15s | %s |%s %-10s %s| %s%s",
			tag,
			time.Now().Format("2006/01/02 - 15:04:05"),
			p.ClientIP,
			p.ConnId,
			p.methodColor(), p.Method.String(), resetColor,
			prefix, message,
		)
	}

	fmt.Fprintf(outputDestination, "%s\n", str)
}

// <summary>: Method用の色を出力します
func (p logParams) methodColor() string {
	if !enableColorMode {
		return ""
	}

	//TODO: まだ作っていない専用色を作る
	switch p.Method {
	case models.CREATE:
		return blue

	case models.JOIN:
		return cyan

	case models.LEAVE:
		return magenta

	case models.BROADCAST:
		return yellow

	case models.CONNECT:
		return bgreen

	case models.DISCONNECT:
		return magenta

	case models.EJECT:
		return magenta

	case models.NOTIFY:
		return yellow

	case models.OK:
		return green

	case models.ERROR:
		return red

	default:
		return white
	}
}
