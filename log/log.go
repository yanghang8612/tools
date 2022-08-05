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
		bytesContent, ok := log.content.([]byte)
		if ok {
			fmt.Printf(titleFormat+"0x%x\n", title, bytesContent)
			continue
		}
		stringContent, ok := log.content.(string)
		if ok {
			fmt.Printf(titleFormat+"%s\n", title, stringContent)
			continue
		}
		bigIntContent, ok := log.content.(*big.Int)
		if ok {
			fmt.Printf(titleFormat+"%d\n", title, bigIntContent)
			continue
		}
	}
}

func NewLog(title string, content interface{}) {
	pendingLogs = append(pendingLogs, Log{title: title, content: content})
}
