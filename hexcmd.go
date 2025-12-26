package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"tools/log"
	"tools/util"

	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/urfave/cli/v2"
)

// 0xd0 range - eof operations.
const (
	CALLTOKEN              vm.OpCode = 0xd0
	TOKENBALANCE           vm.OpCode = 0xd1
	CALLTOKENVALUE         vm.OpCode = 0xd2
	CALLTOKENID            vm.OpCode = 0xd3
	ISCONTRACT             vm.OpCode = 0xd4
	FREEZE                 vm.OpCode = 0xd5
	UNFREEZE               vm.OpCode = 0xd6
	FREEZEEXPIRETIME       vm.OpCode = 0xd7
	VOTEWITNESS            vm.OpCode = 0xd8
	WITHDRAWREWARD         vm.OpCode = 0xd9
	FREEZEBALANCEV2        vm.OpCode = 0xda
	UNFREEZEBALANCEV2      vm.OpCode = 0xdb
	CANCELALLUNFREEZEV2    vm.OpCode = 0xdc
	WITHDRAWEXPIREUNFREEZE vm.OpCode = 0xdd
	DELEGATERESOURCE       vm.OpCode = 0xde
	UNDELEGATERESOURCE     vm.OpCode = 0xdf
)

var tronOpCodeToString = [256]string{
	CALLTOKEN:              "CALLTOKEN",
	TOKENBALANCE:           "TOKENBALANCE",
	CALLTOKENVALUE:         "CALLTOKENVALUE",
	CALLTOKENID:            "CALLTOKENID",
	ISCONTRACT:             "ISCONTRACT",
	FREEZE:                 "FREEZE",
	UNFREEZE:               "UNFREEZE",
	FREEZEEXPIRETIME:       "FREEZEEXPIRETIME",
	VOTEWITNESS:            "VOTEWITNESS",
	WITHDRAWREWARD:         "WITHDRAWREWARD",
	FREEZEBALANCEV2:        "FREEZEBALANCEV2",
	UNFREEZEBALANCEV2:      "UNFREEZEBALANCEV2",
	CANCELALLUNFREEZEV2:    "CANCELALLUNFREEZEV2",
	WITHDRAWEXPIREUNFREEZE: "WITHDRAWEXPIREUNFREEZE",
	DELEGATERESOURCE:       "DELEGATERESOURCE",
	UNDELEGATERESOURCE:     "UNDELEGATERESOURCE",
}

var (
	hexAddrCommand = cli.Command{
		Name:  "addr",
		Usage: "Convert addr between hex, TRON-addr and eth-addr",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("addr subcommand only needs one arg")
			}
			arg0 := c.Args().Get(0)
			addr := new(common.Address)
			if len(arg0) == 34 && arg0[0] == 'T' {
				// TRON address in Base58
				addrBytes, _, err := base58.CheckDecode(arg0)
				if err != nil {
					return err
				}
				addr.SetBytes(addrBytes)
			} else if len(arg0) == 42 && strings.HasPrefix(arg0, "41") {
				// TRON address in hexadecimal with 0x41 prefix
				if addrBytes, ok := utils.FromHex(arg0[2:]); ok {
					addr.SetBytes(addrBytes)
				} else {
					return errors.New("input is begin with 0x41, but can`t convert to hex")
				}
			} else if bigint, ok := math.ParseBig256(arg0); ok {
				// input is in dec or hex
				addr.SetBytes(bigint.Bytes())
			} else {
				return errors.New("input is not recognized, please append 0x prefix if in hex")
			}
			log.NewLog("addr", fmt.Sprintf("%s (%s)", base58.CheckEncode(addr.Bytes(), 0x41), addr.String()))
			return nil
		},
	}
	hexIntCommand = cli.Command{
		Name:  "int",
		Usage: "Convert num between dec and hex",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("int subcommand only needs num arg")
			}
			arg0 := c.Args().Get(0)
			var bigint *big.Int
			var ok bool
			if len(arg0) >= 2 && (arg0[:2] == "0x" || arg0[:2] == "0X") {
				bigint, ok = new(big.Int).SetString(arg0[2:], 16)
			} else {
				bigint, ok = new(big.Int).SetString(arg0, 10)
			}
			if ok {
				log.NewLog("in hex", bigint.Bytes())
				log.NewLog("in dec", fmt.Sprintf("%s (%s len:%d)", bigint.String(), formatBigInt(bigint), len(bigint.String())))
				return nil
			} else {
				return errors.New("only accept input in dec or hex")
			}
		},
	}
	hexMaxCommand = cli.Command{
		Name:  "max",
		Usage: "Get max value for the type like uint-x",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("max subcommand only needs size arg")
			}
			arg0 := c.Args().Get(0)
			size, err := strconv.Atoi(arg0)
			if err != nil {
				return err
			}
			if size <= 0 {
				return errors.New("input uint size should be greater than 0")
			}
			maxValue := new(big.Int)
			maxValue.Lsh(big.NewInt(1), uint(size))
			maxValue.Sub(maxValue, big.NewInt(1))
			log.NewLog("max hex", maxValue.Bytes())
			log.NewLog("max dec", maxValue)
			return nil
		},
	}
	hexStrCommand = cli.Command{
		Name:  "str",
		Usage: "Convert str between ascii and hex",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("hex command only needs single arg")
			}
			arg0 := c.Args().Get(0)

			// if len(arg0) > 1024 {
			// 	return errors.New("input string is too long, max length is 1024 characters")
			// }

			// check if input is in hex
			if argBytes, ok := utils.FromHex(arg0); ok {
				// special case, first byte is `backspace`
				if len(argBytes) > 0 && argBytes[0] == 0x08 {
					argBytes = argBytes[1:]
				}
				log.NewLog("in ascii", utils.ToReadableASCII(argBytes))
			} else {
				// otherwise treat input as str
				log.NewLog("in hex", []byte(arg0))
			}
			return nil
		},
	}
	hexCodeCommand = cli.Command{
		Name:  "code",
		Usage: "Convert hex to bytecode",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("code command only needs single arg")
			}
			arg0 := c.Args().Get(0)

			// check if input is in hex
			if argBytes, ok := utils.FromHex(arg0); ok {
				var sb strings.Builder
				for i := 0; i < len(argBytes); i++ {
					opCode := vm.OpCode(argBytes[i])
					if opCode >= CALLTOKEN && opCode <= UNDELEGATERESOURCE {
						sb.WriteString(fmt.Sprintf("[%d] 0x%02x %s\n", i, argBytes[i], tronOpCodeToString[opCode]))
					} else if opCode.IsPush() {
						dataLen := opCode - vm.PUSH0
						if i+1+int(dataLen) > len(argBytes) {
							break
						}
						data := hex.EncodeToString(argBytes[i+1 : i+1+int(dataLen)])
						sb.WriteString(fmt.Sprintf("[%d] 0x%02x %s 0x%s\n", i, argBytes[i], opCode.String(), data))
						i += int(dataLen) // skip the data bytes
					} else {
						if strings.Contains(opCode.String(), "not defined") {
							break
						}
						sb.WriteString(fmt.Sprintf("[%d] 0x%02x %s\n", i, argBytes[i], opCode.String()))
					}
				}
				log.NewLog("bytecode", "\n"+sb.String())
			} else {
				return errors.New("input is not in hex format")
			}
			return nil
		},
	}
	hexKeyCommand = cli.Command{
		Name:  "key",
		Usage: "Calculate the address corresponding to the private key",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("key command only needs single arg")
			}
			arg0 := c.Args().Get(0)

			// check if input is in hex
			if argBytes, ok := utils.FromHex(arg0); ok {
				if len(argBytes) != 32 {
					return errors.New("input should be 32 bytes hex string")
				}
				privateKey, err := crypto.ToECDSA(argBytes)
				if err != nil {
					return fmt.Errorf("invalid private key: %v", err)
				}
				addr := crypto.PubkeyToAddress(privateKey.PublicKey)
				log.NewLog("key addr", fmt.Sprintf("%s (%s)", base58.CheckEncode(addr.Bytes(), 0x41), addr.String()))
			} else {
				return errors.New("input is not in hex format")
			}
			return nil
		},
	}
)
