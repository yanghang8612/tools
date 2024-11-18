package utils

import "github.com/btcsuite/btcd/btcutil/base58"

func ToAddress(s string) ([]byte, bool) {
	if len(s) == 34 && s[0] == 'T' {
		addrBytes, _, err := base58.CheckDecode(s)
		if err == nil {
			return addrBytes, true
		}
	} else {
		if addrBytes, ok := FromHex(s); ok && len(addrBytes) == 20 {
			return addrBytes, true
		}
	}
	return nil, false
}
