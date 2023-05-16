package utils

import (
	"bytes"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func Has0xPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

func FromHex(s string) ([]byte, bool) {
	if Has0xPrefix(s) {
		return common.FromHex(s), true
	}
	return nil, false
}

func FromDec(s string) (*big.Int, bool) {
	return new(big.Int).SetString(s, 10)
}

func ToReadableASCII(s []byte) string {
	return strings.ReplaceAll(string(bytes.ToValidUTF8(s, nil)), "\n", "↵")
}
