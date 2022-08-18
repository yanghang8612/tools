package main

import (
	"errors"
	"math/big"
	"strconv"
	"strings"
	"tools/log"
	"tools/util"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
)

var (
	hexAddrCommand = cli.Command{
		Name:  "addr",
		Usage: "Convert addr between hex, TRON-addr and eth-addr",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("addr needs single arg")
			}
			input := c.Args().Get(0)
			addr := new(common.Address)

			if input[0] == 'T' {
				bytes, _, err := base58.CheckDecode(input)
				if err != nil {
					return err
				}
				addr.SetBytes(bytes)
			} else {
				if len(input) == 42 && strings.HasPrefix(input, "41") {
					input = input[2:]
				}
				addr.SetBytes(utils.HexToBytes(input))
			}
			log.NewLog("eth addr", addr.String())
			log.NewLog("tron addr", base58.CheckEncode(addr.Bytes(), 0x41))
			return nil
		},
	}
	hexIntCommand = cli.Command{
		Name:  "int",
		Usage: "Convert num between dec and hex",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return errors.New("int subcommand needs at least num arg")
			}
			input := c.Args().Get(0)
			var base = 16
			if c.Args().Len() > 1 {
				base, _ = strconv.Atoi(c.Args().Get(1))
			}
			value := new(big.Int)
			value.SetString(utils.DropHexPrefix(input), base)
			bigValue, _ := uint256.FromBig(value)
			log.NewLog("in hex", bigValue.Hex())
			log.NewLog("in dec", bigValue.ToBig())
			return nil
		},
	}
	hexMaxCommand = cli.Command{
		Name:  "max",
		Usage: "Get max value for the type like uint-x",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("max subcommand needs size arg")
			}
			size, _ := strconv.Atoi(c.Args().Get(0))
			if size < 8 || size > 256 {
				return errors.New("input uint size should be in [8~256]")
			}
			if size%8 != 0 {
				return errors.New("input uint size should be mod by 8")
			}
			bigValue, _ := uint256.FromHex("0x100")
			one, _ := uint256.FromHex("0x1")
			bigValue.Lsh(bigValue, uint(size-8)).Sub(bigValue, one)
			log.NewLog("max hex", bigValue.Hex())
			log.NewLog("max dec", bigValue.ToBig())
			return nil
		},
	}
	hexStrCommand = cli.Command{
		Name:  "str",
		Usage: "convert hex between str",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("hex command needs one arg")
			}
			arg := c.Args().Get(0)
			// input is num in hex
			if utils.ContainHexPrefix(arg) {
				argBytes := hexutils.HexToBytes(utils.DropHexPrefix(arg))
				if len(argBytes) <= 32 {
					log.NewLog("in decimal", new(big.Int).SetBytes(argBytes))
				}
				// special case, first byte is `backspace`
				if len(argBytes) > 0 && argBytes[0] == 0x08 {
					argBytes = argBytes[1:]
				}
				log.NewLog("in ascii", utils.ToReadableASCII(argBytes))
			} else {
				// otherwise input is str
				log.NewLog("in hex", []byte(arg))
			}
			return nil
		},
	}
)
