package log

import (
	"fmt"
	"math/big"
	"strconv"
)

type Log struct {
	title   string
	content interface{}
}

var pendingLogs []Log

func FlushLogsToConsole() {
	maxTitleLength := 0
	for _, log := range pendingLogs {
		if len(log.title)+2 > maxTitleLength {
			maxTitleLength = len(log.title) + 2
		}
	}
	titleFormat := "%" + strconv.Itoa(maxTitleLength) + "s - "
	for _, log := range pendingLogs {
		title := "[" + log.title + "]"
		switch log.content.(type) {
		case []byte, [32]byte:
			fmt.Printf(titleFormat+"0x%x\n", title, log.content)
		case string:
			fmt.Printf(titleFormat+"%s\n", title, log.content)
		case int, uint, int64, uint64, *big.Int, big.Int:
			fmt.Printf(titleFormat+"%d\n", title, log.content)
		default:
			fmt.Printf(titleFormat+"%v\n", title, log.content)
		}
	}
}

func NewLog(title string, content interface{}) {
	pendingLogs = append(pendingLogs, Log{title: title, content: content})
}
