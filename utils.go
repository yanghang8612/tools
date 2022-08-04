package main

import "strings"

func containHexPrefix(str string) bool {
	return strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X")
}

func dropHexPrefix(str string) string {
	if containHexPrefix(str) {
		str = str[2:]
	}
	return str
}
