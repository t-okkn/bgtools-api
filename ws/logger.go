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

// <summary>: Websocket用ログのパラメータを示す構造体
type logParams struct {
	ClientIP    string
	ConnId      string
	Method      models.Method
	IsProcError bool
	Message     string
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
func newLogParams() logParams {
	return logParams{
		ClientIP:    "",
		ConnId:      "",
		Method:      models.NONE,
		IsProcError: false,
		Message:     "",
	}
}

// <summary>: logParamsの情報からログを書き込みます
func (p logParams) log() {
	resetColor := reset
	if !enableColorMode {
		resetColor = ""
	}

	tag := "WS"
	var str string

	if p.IsProcError {
		tag = "WS-Error"

		str = fmt.Sprintf("[%s] %v | %15s | %s | 【処理エラー】%s",
			tag,
			time.Now().Format("2006/01/02 - 15:04:05"),
			p.ClientIP,
			p.ConnId,
			p.Message,
		)
	} else {
		str = fmt.Sprintf("[%s] %v | %15s | %s |%s %-10s %s %s",
			tag,
			time.Now().Format("2006/01/02 - 15:04:05"),
			p.ClientIP,
			p.ConnId,
			p.methodColor(), p.Method.String(), resetColor,
			p.Message,
		)
	}

	fmt.Fprintf(outputDestination, "%s\n", str)
}

// <summary>: Method用の色を出力します
func (p logParams) methodColor() string {
	if !enableColorMode {
		return ""
	}

	switch p.Method {
	case models.CREATE_ROOM:
		return blue

	case models.JOIN_ROOM:
		return cyan

	case models.BROADCAST:
		return yellow

	case models.CONNCTED:
		return bgreen

	case models.ACCEPTED:
		return green

	case models.ERROR:
		return red

	default:
		return white
	}
}

// <summary>: 処理にかかった時間を文字列で取得します
func getProcTime(start time.Time) string {
	return fmt.Sprintf("%v", time.Since(start))
}
