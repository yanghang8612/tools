package main

import (
	utils "tools/util"

	"errors"
	"tools/log"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/urfave/cli/v2"
)

var (
	addrNumCommand = cli.Command{
		Name:  "num",
		Usage: "Pad the given num to address",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("num subcommand needs num arg")
			}
			addr := common.HexToAddress(c.Args().Get(0))
			log.NewLog("eth addr", addr.String())
			log.NewLog("tron addr", base58.CheckEncode(addr.Bytes(), 0x41))
			return nil
		},
	}
	addrDecodeCommand = cli.Command{
		Name:  "decode",
		Usage: "Decode base58 encoded address to eth address",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("decode subcommand needs base58 encoded address arg")
			}
			base58Addr := c.Args().Get(0)
			if len(base58Addr) != 34 {
				return errors.New("base58 encoded address length must be 34")
			}
			if base58Addr[0] != 'T' {
				return errors.New("base58 encoded address must begin with 'T'")
			}
			addr, _, err := base58.CheckDecode(base58Addr)
			if err == nil {
				log.NewLog("eth addr", common.BytesToAddress(addr).String())
				return nil
			} else {
				return err
			}
		},
	}
	addrEncodeCommand = cli.Command{
		Name:  "encode",
		Usage: "Encode eth address to base58 encoded address",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return errors.New("decode subcommand needs eth address arg")
			}
			addr := utils.DropHexPrefix(c.Args().Get(0))
			if len(addr) != 40 {
				return errors.New("eth address length must be 20")
			}
			ethAddr := hexutils.HexToBytes(addr)
			log.NewLog("tron addr", base58.CheckEncode(ethAddr, 0x41))
			return nil
		},
	}
)
