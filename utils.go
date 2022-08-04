package main

import (
	"strings"
)

func containHexPrefix(str string) bool {
	return strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X")
      || strings.HasPrefix(str, "x") || strings.HasPrefix(str, "X")
}

func dropHexPrefix(str string) string {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	} else if strings.HasPrefix(str, "x") || strings.HasPrefix(str, "X") {
		str = str[1:]
	}
	return str
}
