package main

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
	"math/big"
	"strconv"
	"strings"
)

var (
	addrCommand = cli.Command{
		Name:  "addr",
		Usage: "convert hex between hex, tron-addr, eth-addr",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("1 input only")
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
				input := dropHexPrefix(input)
				addr.SetBytes(hexutils.HexToBytes(input))
			}

			fmt.Printf("%s result:\n  hex:%s\n  tron: %s\n  eth: %s\n",
				input, addr.Hex(), base58.CheckEncode(addr.Bytes(), 0x41), addr.String())
			return nil
		},
	}

	intCommand = cli.Command{
		Name:  "int",
		Usage: "convert hex between hex, int",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return errors.New("empty input")
			}
			input := c.Args().Get(0)
			var base = 16
			if c.Args().Len() > 1 {
				base, _ = strconv.Atoi(c.Args().Get(1))
			}
			input = dropHexPrefix(input)
			value := new(big.Int)
			value.SetString(input, base)
			bigValue, _ := uint256.FromBig(value)

			fmt.Printf("%s (base %d) result:\n  uint: %s\n  hex: %s\n", input, base, bigValue, bigValue.Hex())
			return nil
		},
	}

	maxCommand = cli.Command{
		Name:  "max",
		Usage: "get max uint-x, example: max 256",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("input uint size [8~256]")
			}
			size, _ := strconv.Atoi(c.Args().Get(0))
			if size < 8 || size > 256 || size%8 != 0 {
				return errors.New("input uint size [8~256]")
			}
			bigValue, _ := uint256.FromHex("0x100")
			one, _ := uint256.FromHex("0x1")
			bigValue.Lsh(bigValue, uint(size-8)).Sub(bigValue, one)

			fmt.Printf("max uint%d result:\n  uint: %s\n  hex: %s\n", size, bigValue, bigValue.Hex())
			return nil
		},
	}

	strCommand = cli.Command{
		Name:  "str",
		Usage: "convert hex between str",
		Action: func(c *cli.Context) error {
			return errors.New("NOT SUPPORT YET")
		},
	}
)
