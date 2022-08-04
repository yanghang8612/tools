package main

import "strings"

func dropHexPrefix(str string) string {
	if strings.HasPrefix(str, "x") || strings.HasPrefix(str, "X") {
		str = str[1:]
	}
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}
	return str
}
