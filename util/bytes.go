package utils

import (
	"bytes"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func Has0xPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

func FromHex(s string) ([]byte, bool) {
	// try to decode hex
	if data, err := hexutil.Decode(s); err == nil {
		return data, true
	}

	// otherwise, it may be string or invalid hex
	return nil, false
}

func FromDec(s string) (*big.Int, bool) {
	return new(big.Int).SetString(s, 10)
}

func ToReadableASCII(s []byte) string {
	return strings.ReplaceAll(string(bytes.ToValidUTF8(s, nil)), "\n", "↵")
}
