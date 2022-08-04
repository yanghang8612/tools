package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
)

var (
	vmPadCommand = cli.Command{
		Name:  "pad",
		Usage: "Pad num or bytes to 32bytes",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("pad subcommand needs num or bytes arg")
			}
			arg := c.Args().Get(0)
			// input is bytes
			if containHexPrefix(arg) {
				arg = dropHexPrefix(arg)
				argBytes := hexutils.HexToBytes(arg)
				if len(argBytes) > 32 {
					return errors.New("input type is likely bytes, but its length is greater than 32, there is no need to pack it")
				}
				lackBytes := make([]byte, 32-len(argBytes))
				bigEndianRes := append(lackBytes, argBytes...)
				fmt.Printf("[padded hex in big endian] 0x%x\n", bigEndianRes)
				littleEndianRes := append(argBytes, lackBytes...)
				fmt.Printf("[padded hex in little endian] 0x%x\n", littleEndianRes)
			} else {
				// otherwise input must be num
				if num, ok := new(big.Int).SetString(c.Args().Get(0), 10); ok {
					fmt.Printf("[origin hex] 0x%x\n", num.Bytes())
					res := make([]byte, 32-len(num.Bytes()))
					res = append(res, num.Bytes()...)
					fmt.Printf("[padded hex] 0x%x\n", res)
				} else {
					return errors.New("input type is likely num, but cannot covert it")
				}
			}
			return nil
		},
	}
	vmSplitCommand = cli.Command{
		Name:  "split",
		Usage: "Spilt data to each 32bytes",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("split subcommand needs data arg")
			}
			data := hexutils.HexToBytes(dropHexPrefix(c.Args().Get(0)))
			if len(data)%32 == 4 {
				rspData := doGet(fmt.Sprintf("https://www.4byte.directory/api/v1/signatures/"+
					"?hex_signature=%x", data[:4]))
				var rsp Rsp4Bytes
				err := json.Unmarshal(rspData, &rsp)
				if err != nil {
					return err
				}

				if rsp.Count != 0 {
					fmt.Printf("[selector]: %x - %s\n", data[:4], rsp.Results[rsp.Count-1].Signature)
				} else {
					fmt.Printf("[selector]: %x\n", data[:4])
				}
				data = data[4:]
			}
			if len(data)%32 != 0 {
				return errors.New("data must be 32*N")
			}
			fmt.Println("[each data word]:")
			format := "0x%02x: %x\n"
			if len(data) > 8*32 {
				format = "0x%03x: %x\n"
			}
			for i := 0; i < len(data)/32; i++ {
				fmt.Printf(format, i*32, data[i*32:i*32+32])
			}
			return nil
		},
	}
	vmUnpackCommand = cli.Command{
		Name:  "unpack",
		Usage: "Unpack data with given types",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				return errors.New("unpack subcommand needs data and type args")
			}
			data, _ := hex.DecodeString(dropHexPrefix(c.Args().Get(1)))
			args := abi.Arguments{}
			for _, arg := range strings.Split(c.Args().Get(0), ",") {
				solType, _ := abi.NewType(arg, "", nil)
				args = append(args, abi.Argument{Type: solType})
			}
			if res, err := args.UnpackValues(data); err == nil {
				fmt.Printf("[unpack result]:\n")
				for i, r := range res {
					printSol(r, &args[i].Type, "arg", i, 1)
				}
				return nil
			} else {
				return err
			}
		},
	}
	vm4bytesCommand = cli.Command{
		Name:  "4bytes",
		Usage: "Get 4bytes selector for given func or event",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("4bytes subcommand needs func or event signature arg")
			}
			signature := c.Args().Get(0)
			selector := crypto.Keccak256([]byte(signature))[:4]
			fmt.Printf("[origin hex] 0x%x\n", selector)
			res := make([]byte, 32-len(selector))
			selector = append(selector, res...)
			fmt.Printf("[padded hex] 0x%x\n", selector)
			return nil
		},
	}
)

func printSol(param interface{}, paramTy *abi.Type, name string, index, offset int) {
	printSeparator(offset, "  ", "", "- ")
	switch paramTy.T {
	case abi.ArrayTy:
		fmt.Printf("[%s-%02d]: %s\n", name, index, paramTy.String())
		paramArray := reflect.ValueOf(param)
		for i := 0; i < paramArray.Len(); i++ {
			printSol(paramArray.Index(i).Interface(), paramTy.Elem, "array", i, offset+1)
		}
	case abi.SliceTy:
		fmt.Printf("[%s-%02d]: %s\n", name, index, paramTy.String())
		paramSlice := reflect.ValueOf(param)
		for i := 0; i < paramSlice.Len(); i++ {
			printSol(paramSlice.Index(i).Interface(), paramTy.Elem, "slice", i, offset+1)
		}
	case abi.BytesTy, abi.FixedBytesTy:
		fmt.Printf("[%s-%02d]: %s, %#x\n", name, index, paramTy.String(), param)
	case abi.AddressTy:
		fmt.Printf("[%s-%02d]: %s, %v - %s\n", name, index, paramTy.String(), param, base58.CheckEncode(param.(common.Address).Bytes(), 0x41))
	case abi.IntTy, abi.UintTy:
		fmt.Printf("[%s-%02d]: %s, %v", name, index, paramTy.String(), param)
		intWithDot := formatBigInt(param.(*big.Int))
		if strings.ContainsAny(intWithDot, ",") {
			fmt.Printf(" - %s", intWithDot)
		}
		if len(param.(*big.Int).String()) >= 6 {
			fmt.Printf(" (%d)", len(param.(*big.Int).String()))
		}
		fmt.Println()
	default:
		fmt.Printf("[%s-%02d]: %s, %v\n", name, index, paramTy.String(), param)
		//fmt.Printf("[Parameter-%d]: %T, %#x\n", index, param, param)
	}
}

func printSeparator(repeat int, symbol, prefix, suffix string) {
	fmt.Print(prefix)
	for i := 0; i < repeat; i++ {
		fmt.Print(symbol)
	}
	fmt.Print(suffix)
}

func formatBigInt(n *big.Int) string {
	var (
		text  = n.String()
		buf   = make([]byte, len(text)+len(text)/3)
		comma = 0
		i     = len(buf) - 1
	)
	for j := len(text) - 1; j >= 0; j, i = j-1, i-1 {
		c := text[j]

		switch {
		case c == '-':
			buf[i] = c
		case comma == 3:
			buf[i] = ','
			i--
			comma = 0
			fallthrough
		default:
			buf[i] = c
			comma++
		}
	}
	return string(buf[i+1:])
}
